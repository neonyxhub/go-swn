package swn

import (
	"bufio"
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	auth_pb "go.neonyx.io/go-swn/internal/swn/pb"
	"go.neonyx.io/go-swn/pkg/crypto"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNotAuthoried       = errors.New("not authorized")
	ErrAuthEmptyData      = errors.New("empty auth data")
	ErrAuthEmptyChallenge = errors.New("empty auth challenge")
)

// Perform swn authentification with destination string
func (s *SWN) Auth(destination string) (bool, error) {
	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		return false, err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		s.Log.Sugar().Infoln(err)
		return false, err
	}
	peerId, _ := info.ID.Marshal()

	deviceId, err := s.GetDeviceForPeerId(&auth_pb.PeerId{Data: peerId})

	if err != nil {
		s.Log.Sugar().Infoln(err)
		return false, err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	s.Peer.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	stream, err := s.Peer.Host.NewStream(context.Background(), info.ID, "/swnauth/1.0.0")
	if err != nil {
		return false, fmt.Errorf(`error creating stream: %v`, err)
	}
	s.Log.Sugar().Infoln("Started auth process")

	// Create a buffered stream so that read and writes are non-blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	c := make(chan bool)

	var wg sync.WaitGroup

	wg.Add(1)
	go s.performOutcomingAuth(rw, deviceId, wg.Done, c)
	wg.Wait()

	auth := <-c
	return auth, nil
}

// Function to perform incoming auth
func (s *SWN) PerformIncomingAuth(rw *bufio.ReadWriter) error {
	senderDevice, _ := rw.ReadString('\n')
	if senderDevice == "" {
		return ErrAuthEmptyData
	}

	var incomeDeviceId auth_pb.DeviceId
	err := proto.Unmarshal([]byte(senderDevice), &incomeDeviceId)
	if err != nil {
		return err
	}
	s.Log.Sugar().Infof("Got device id: %v\n", string(incomeDeviceId.DeviceId))

	authRow, err := s.GetAuthRowFromDeviceId(&incomeDeviceId)
	if err != nil {
		return err
	}
	getterPriv, _ := x509.ParsePKCS1PrivateKey(authRow.MyDevicePrivateKey)

	senderPub, _ := x509.ParsePKCS1PublicKey(authRow.OtherDevicePubKey)

	nonce, _ := crypto.GetNonce()

	challenge, err := crypto.GenerateChallenge(senderPub, nonce)
	if err != nil {
		return err
	}

	s.Log.Sugar().Infof("Generated challenge: %v\n", string(challenge))
	rw.Write(append(challenge, byte('\n')))
	rw.Flush()

	signedChallenge, _ := rw.ReadString('\n')
	if signedChallenge == "" {
		return ErrAuthEmptyChallenge
	}
	s.Log.Sugar().Infof("Got hashed data: %v\n", signedChallenge)

	auth, err := crypto.CheckResponse([]byte(signedChallenge), nonce, getterPriv)
	if err != nil {
		return err
	}

	if auth {
		rw.WriteString(fmt.Sprintf("%s\n", "nice"))
		rw.Flush()
		return nil
	} else {
		rw.WriteString(fmt.Sprintf("%s\n", "notnice"))
		rw.Flush()
		return ErrNotAuthoried
	}

}

func (s *SWN) GetAuthRowFromDeviceId(deviceId *auth_pb.DeviceId) (*auth_pb.AuthInfo, error) {
	key := s.Ds.NewKey("auth_storage/incoming/" + string(deviceId.DeviceId))
	authRaw, err := s.Ds.Get(key, nil)

	if err != nil {
		return nil, err
	}
	var authRow auth_pb.AuthInfo
	err = proto.Unmarshal(authRaw, &authRow)

	return &authRow, err
}

// Do the authentification challenge-response
func (s *SWN) performOutcomingAuth(rw *bufio.ReadWriter, deviceId *auth_pb.DeviceId, done func(), c chan bool) {
	defer done()
	rw.WriteString(string(deviceId.DeviceId) + "\n")
	rw.Flush()

	challenge, _ := rw.ReadString('\n')
	s.Log.Sugar().Infof("got challenge: %v", challenge)

	authRow, err := s.GetAuthInfo(deviceId)

	if err != nil {
		return
	}
	senderPriv, _ := x509.ParsePKCS1PrivateKey(authRow.MyDevicePrivateKey)
	getterPub, _ := x509.ParsePKCS1PublicKey(authRow.OtherDevicePubKey)

	ans, err := crypto.DecryptWithPrivateKey([]byte(challenge), senderPriv)

	if err != nil {
		return
	}
	resp, err := crypto.EncryptWithPublicKey(ans, getterPub)

	if err != nil {
		return
	}

	rw.WriteString(string(resp) + "\n")
	rw.Flush()

	str, _ := rw.ReadString('\n')
	s.Log.Sugar().Infof("got response: %v", str)

	c <- true
}

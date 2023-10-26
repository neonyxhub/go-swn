package swn

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/x509"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"google.golang.org/protobuf/proto"

	"github.com/go-errors/errors"
	"go.neonyx.io/go-swn/internal/swn/pb"
	"go.neonyx.io/go-swn/pkg/crypto"
)

const (
	AUTH_OK      = "OK"
	AUTH_NOK     = "NOK"
	AUTH_TIMEOUT = 10 * time.Second
)

var (
	ErrNotAuthorized = errors.Errorf("not authorized")
)

func (s *SWN) IsAuthenticated(conn network.Conn) bool {
	if conn.IsClosed() {
		s.Log.Sugar().Infof("connection for %v is closed, not authorized\n", conn.RemotePeer())
		delete(s.AuthDeviceMap, conn.RemotePeer().String())
		return false
	}

	_, ok := s.AuthDeviceMap[conn.RemotePeer().String()]
	return ok
}

// Perform incoming auth from AuthHandler
func (s *SWN) AuthIn(stream network.Stream) error {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// 0. sends ACK/NACK if authenticated
	if s.IsAuthenticated(stream.Conn()) {
		if err := WriteB64(rw, []byte(AUTH_OK)); err != nil {
			return err
		}
		return nil
	} else {
		if err := WriteB64(rw, []byte(AUTH_NOK)); err != nil {
			return err
		}
	}

	// 1. receives DeviceAuthRequest from sender
	s.Log.Info("waiting for DeviceAuthRequest")
	reqRaw, err := ReadB64(rw)
	if err != nil {
		return err
	}
	req := &pb.DeviceAuthRequest{}
	if err = proto.Unmarshal(reqRaw, req); err != nil {
		return errors.Errorf("failed to Unmarshal DeviceAuthRequest: %v", err)
	}

	// will be stored if challenge-response is ACK
	senderDeviceId := req.Data

	// challenge
	senderDevPubKey, err := x509.ParsePKCS1PublicKey(req.SenderDevPubKey)
	if err != nil {
		return errors.Errorf("failed to parse SenderDevPubKey: %v", err)
	}
	nonce, err := crypto.GetNonce()
	if err != nil {
		return err
	}
	challenge, err := crypto.GenerateChallenge(senderDevPubKey, nonce)
	if err != nil {
		return errors.Errorf("failed to GenerateChallenge: %v", err)
	}

	// 2. send to sender a challenge, encrypted with sender's pubkey
	if err = WriteB64(rw, challenge); err != nil {
		return err
	}

	// 3. receive's response on challenge
	senderHashedNonce, err := ReadB64(rw)
	if err != nil {
		return err
	}

	// 4. send ACK/NACK
	localNonceHash := sha256.Sum256(nonce)
	if bytes.Equal(localNonceHash[:], senderHashedNonce) {
		if err = WriteB64(rw, []byte(AUTH_OK)); err != nil {
			return err
		}

		// register current connection auth
		s.AuthDeviceMap[stream.Conn().RemotePeer().String()] = senderDeviceId
		s.Log.Sugar().Infof("authenticated peer with deviceId=%v", senderDeviceId)

		return nil
	} else {
		if err = WriteB64(rw, []byte(AUTH_NOK)); err != nil {
			return err
		}

		return ErrNotAuthorized
	}
}

// Perform outgoing swn authentification with given multiaddress destination string
func (s *SWN) AuthOut(destination string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AUTH_TIMEOUT)
	defer cancel()

	destInfo, err := s.Peer.AddRemotePeer(destination, AUTH_TIMEOUT)
	if err != nil {
		return false, err
	}

	stream, err := s.Peer.Host.NewStream(ctx, destInfo.ID, HID_AUTH)
	if err != nil {
		return false, err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// 0. receive ACK/NACK if authenticated
	s.Log.Info("reading ACK/NACK if authenticated")
	resp, err := ReadB64(rw)
	if err != nil {
		return false, err
	}

	if len(resp) == 2 || len(resp) == 3 {
		if bytes.Equal(resp, []byte(AUTH_OK)) {
			s.Log.Info("already authenticated!")
			return true, nil
		}
	} else {
		return false, errors.Errorf("irrelevant auth response: %v", resp)
	}

	reqRaw, err := proto.Marshal(&pb.DeviceAuthRequest{
		Data:            s.Device.Id,
		SenderDevPubKey: s.Device.GetPubKeyRaw(),
	})
	if err != nil {
		return false, err
	}

	s.Log.Info("sending local deviceId to remote swn auth")
	if err = WriteB64(rw, reqRaw); err != nil {
		return false, err
	}

	// 2. receive challenge with encrypted nonce from outgoing swn
	s.Log.Info("reading a challenge from remote swn")
	challenge, err := ReadB64(rw)
	if err != nil {
		return false, err
	}

	nonce, err := crypto.DecryptWithPrivateKey([]byte(challenge), s.Device.PrivKey)
	if err != nil {
		return false, err
	}

	hashedNonce := sha256.Sum256(nonce)

	// 3. response to outgoing swn with hashed nonce
	s.Log.Info("responding to remote swn's challenge")
	if err = WriteB64(rw, hashedNonce[:]); err != nil {
		return false, err
	}

	// 4. receive ACK from destination that nonce is valid
	ack, err := ReadB64(rw)
	if err != nil {
		return false, err
	}

	if string(ack) == AUTH_OK {
		s.Log.Info("received OK on AuthOut")
		return true, nil
	} else {
		s.Log.Info("received NOK on AuthOut")
		return false, ErrNotAuthorized
	}
}

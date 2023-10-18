package swn

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"google.golang.org/protobuf/proto"

	"github.com/go-errors/errors"
	"go.neonyx.io/go-swn/internal/swn/pb"
	"go.neonyx.io/go-swn/pkg/crypto"
)

const (
	AUTH_ACK     = "ACK"
	AUTH_NACK    = "NACK"
	AUTH_TIMEOUT = 10 * time.Second
)

var (
	ErrNotAuthorized = errors.Errorf("not authorized")

	// connId: deviceId
	AuthDeviceMap = make(map[string][]byte)
)

func (s *SWN) IsAuthorized(conn network.Conn) bool {
	if conn.IsClosed() {
		s.Log.Sugar().Infof("stream %v is closed, not authorized\n", conn.ID())
		delete(AuthDeviceMap, conn.ID())
		return false
	}

	_, ok := AuthDeviceMap[conn.ID()]
	return ok
}

func writeB64(rw *bufio.ReadWriter, req []byte) error {
	encoded := base64.StdEncoding.EncodeToString(req)

	if _, err := rw.WriteString(encoded + "\n"); err != nil {
		return err
	}

	if err := rw.Flush(); err != nil {
		return err
	}

	return nil
}

func readB64(rw *bufio.ReadWriter) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AUTH_TIMEOUT)
	defer cancel()

	resultCh := make(chan []byte)
	errorCh := make(chan error, 1)

	go func() {
		// Read the encoded message from the stream.
		encodedReq, err := rw.ReadString('\n')
		if err != nil {
			errorCh <- err
			return
		}
		if len(encodedReq) == 0 {
			errorCh <- errors.Errorf("empty buffer")
			return
		}

		req, err := base64.StdEncoding.DecodeString(encodedReq)
		if err != nil {
			errorCh <- err
			return
		}

		resultCh <- req
	}()

	select {
	case res := <-resultCh:
		return res, nil
	case err := <-errorCh:
		return nil, err
	case <-ctx.Done():
		return nil, errors.Errorf("auth timeout: %v sec passed", AUTH_TIMEOUT)
	}
}

// Perform incoming auth from AuthHandler
func (s *SWN) AuthIn(stream network.Stream) error {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	if s.IsAuthorized(stream.Conn()) {
		if err := writeB64(rw, []byte(AUTH_ACK)); err != nil {
			return err
		}
		return nil
	}

	// 0. sends local device public key to sender
	if err := writeB64(rw, s.Device.GetPubKeyRaw()); err != nil {
		return err
	}

	// 1. receives DeviceAuthRequest from sender
	reqRaw, err := readB64(rw)
	if err != nil {
		return err
	}
	req := &pb.DeviceAuthRequest{}
	if err = proto.Unmarshal(reqRaw, req); err != nil {
		return errors.Errorf("failed to Unmarshal DeviceAuthRequest: %v", err)
	}

	data, err := crypto.DecryptWithPrivateKey(req.Data, s.Device.PrivKey)
	if err != nil || len(data) == 0 {
		return errors.Errorf("failed to DecryptWithPrivateKey: %v", err)
	}

	// will be stored if challenge-response is ACK
	senderDeviceId := data

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
	if err = writeB64(rw, challenge); err != nil {
		return err
	}

	// 3. receive's response on challenge
	senderHashedNonce, err := readB64(rw)
	if err != nil {
		return err
	}

	// 4. send ACK/NACK
	localNonceHash := sha256.Sum256(nonce)
	if bytes.Equal(localNonceHash[:], senderHashedNonce) {
		if err = writeB64(rw, []byte(AUTH_ACK)); err != nil {
			return err
		}

		AuthDeviceMap[stream.Conn().ID()] = senderDeviceId
		s.Log.Sugar().Infof("authenticated peer with deviceId=%v", senderDeviceId)

		return nil
	} else {
		if err = writeB64(rw, []byte(AUTH_NACK)); err != nil {
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

	// 0. receive from destination its device public key or ACK if authenticated
	s.Log.Info("reading a destination device public key")
	resp, err := readB64(rw)
	if err != nil {
		return false, err
	}

	if len(resp) == 3 && bytes.Equal(resp, []byte(AUTH_ACK)) {
		s.Log.Info("already authenticated!")
		return true, nil
	}

	destPubKey, err := x509.ParsePKCS1PublicKey(resp)
	if err != nil {
		return false, err
	}

	// 1. send current device Id, encrypting with destination pubkey
	encDevId, err := crypto.EncryptWithPublicKey(s.Device.Id, destPubKey)
	if err != nil {
		return false, err
	}
	reqRaw, err := proto.Marshal(&pb.DeviceAuthRequest{
		Data:            encDevId,
		SenderDevPubKey: s.Device.GetPubKeyRaw(),
	})
	if err != nil {
		return false, err
	}

	s.Log.Info("sending local deviceId to remote swn auth")
	if err = writeB64(rw, reqRaw); err != nil {
		return false, err
	}

	// 2. receive challenge with encrypted nonce from outgoing swn
	s.Log.Info("reading a challenge from remote swn")
	challenge, err := readB64(rw)
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
	if err = writeB64(rw, hashedNonce[:]); err != nil {
		return false, err
	}

	// 4. receive ACK from destination that nonce is valid
	ack, err := readB64(rw)
	if err != nil {
		return false, err
	}

	if string(ack) == AUTH_ACK {
		s.Log.Info("received ACK on AuthOut")
		return true, nil
	} else {
		s.Log.Info("received NACK on AuthOut")
		return false, ErrNotAuthorized
	}
}

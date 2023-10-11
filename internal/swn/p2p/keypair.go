package p2p

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"os"

	libp2p_crypto "github.com/libp2p/go-libp2p/core/crypto"
)

var (
	ErrInvalidPrivKeyEd25519 = errors.New("invalid ed25519 private key")
)

type KeyPair25519 struct {
	PrivKey     libp2p_crypto.PrivKey
	PubKey      libp2p_crypto.PubKey
	PrivKeyPath string
	IsGenerated bool
}

func (kp *KeyPair25519) ReadFromFile() error {
	prvKeyBytes, err := os.ReadFile(kp.PrivKeyPath)
	if err != nil {
		return err
	}

	if len(prvKeyBytes) != ed25519.PrivateKeySize {
		return ErrInvalidPrivKeyEd25519
	}

	//prvKey := ed25519.PrivateKey(prvKeyBytes)
	//pubKey := prvKey.Public()

	kp.PrivKey = nil
	kp.PubKey = nil
	kp.IsGenerated = false

	return nil
}

func (kp *KeyPair25519) Gen() error {
	// TODO: store priv & pub key at config.P2p.PrivKeyPath
	prvKey, pubKey, err := libp2p_crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return err
	}

	kp.PrivKey = prvKey
	kp.PubKey = pubKey
	kp.IsGenerated = true

	return nil
}

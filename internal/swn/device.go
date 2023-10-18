package swn

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

type Device struct {
	// MSB 12-bytes of swn "hardware" pubkey hash
	Id      []byte
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
}

func (d *Device) GetPrivKeyRaw() []byte {
	return x509.MarshalPKCS1PrivateKey(d.PrivKey)
}

func (d *Device) GetPubKeyRaw() []byte {
	return x509.MarshalPKCS1PublicKey(d.PubKey)
}

func (d *Device) ParsePrivKeyRaw(der []byte) error {
	prvKey, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return err
	}
	d.PrivKey = prvKey

	return nil
}

func (d *Device) ParsePubKeyRaw(der []byte) error {
	pubKey, err := x509.ParsePKCS1PublicKey(der)
	if err != nil {
		return err
	}
	d.PubKey = pubKey

	return nil
}

func (d *Device) GenKeyPair() error {
	prvKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	d.PrivKey = prvKey
	d.PubKey = &prvKey.PublicKey

	return nil
}

func (d *Device) GenDeviceId() error {
	pubKeyRaw := d.GetPubKeyRaw()
	hash := sha256.New()
	if _, err := hash.Write(pubKeyRaw); err != nil {
		return err
	}

	d.Id = hash.Sum(nil)[:12]

	return nil
}

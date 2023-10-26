package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"reflect"
)

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// Encrypts nonce with responder public key
func GenerateChallenge(pubKey *rsa.PublicKey, nonce []byte) ([]byte, error) {
	decryptedChallenge := nonce

	challenge, err := EncryptWithPublicKey(decryptedChallenge, pubKey)
	return challenge, err
}

// Check if response was reencoded in right way
func CheckResponse(response []byte, nonce []byte, priv *rsa.PrivateKey) (bool, error) {
	ans, err := DecryptWithPrivateKey(response, priv)

	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(nonce, ans), nil
}

// Generated random 32 bytes as nonce
func GetNonce() ([]byte, error) {
	nonce := make([]byte, 32)

	_, err := rand.Read(nonce)

	if err != nil {
		return nil, err
	}
	return nonce, nil
}

func hashNonce(nonce []byte) {

}

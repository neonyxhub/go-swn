package crypto_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"

	"go.neonyx.io/go-swn/pkg/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pub := &priv.PublicKey
	msg := []byte("secret message")

	cipher, err := crypto.EncryptWithPublicKey(msg, pub)
	require.NoError(t, err)

	plaintext, err := crypto.DecryptWithPrivateKey(cipher, priv)
	require.NoError(t, err)
	require.Equal(t, "secret message", string(plaintext))
}

package swn_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.neonyx.io/go-swn/internal/swn"
)

func TestGenKeyPair(t *testing.T) {
	d := &swn.Device{}
	err := d.GenKeyPair()
	require.NoError(t, err)

	require.NoError(t, err)
	require.NotEmpty(t, d.GetPrivKeyRaw())

	require.NoError(t, err)
	require.NotEmpty(t, d.GetPubKeyRaw())
}

func TestGenDeviceId(t *testing.T) {
	d := &swn.Device{}
	err := d.GenKeyPair()
	require.NoError(t, err)

	err = d.GenDeviceId()
	require.NoError(t, err)
	require.Equal(t, 12, len(d.Id))
}

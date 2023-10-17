package swn_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAuthInfo(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)

	deviceAuth, err := swn.GetDeviceAuth()
	require.NoError(t, err)
	require.NotEmpty(t, deviceAuth.PrivKey)
}

func TestSaveAuthInfo(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)

	err = swn.SaveDeviceAuth()
	require.NoError(t, err)

	deviceAuth, err := swn.GetDeviceAuth()
	require.NoError(t, err)
	require.NotEmpty(t, deviceAuth.PrivKey)
}

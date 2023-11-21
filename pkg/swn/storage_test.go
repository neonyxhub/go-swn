package swn_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/syndtr/goleveldb/leveldb"
)

func TestGetAuthInfo(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)

	err = swn.GetDeviceAuth()
	require.NoError(t, err)
	require.NotEmpty(t, swn.Device.PrivKey)
}

func TestSaveAuthInfo(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)

	err = swn.SaveDeviceAuth()
	require.NoError(t, err)

	err = swn.GetDeviceAuth()
	require.NoError(t, err)
	require.NotEmpty(t, swn.Device.PrivKey)
}

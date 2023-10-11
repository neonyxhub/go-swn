package swn_test

import (
	mrand "math/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/require"
)

type Pair struct {
	Priv crypto.PrivKey
	Pub  crypto.PubKey
}

func genPrivPubKeys(t *testing.T) (client *Pair) {
	r := mrand.New(mrand.NewSource(42))

	clientPriv, clientPub, err := crypto.GenerateEd25519Key(r)
	require.NoError(t, err)

	client = &Pair{Priv: clientPriv, Pub: clientPub}

	return
}

func TestGetAuthInfo(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)

	clientPair := genPrivPubKeys(t)
	deviceId, err := swn.SaveAuthInfo(swn.Peer.KeyPair.PubKey, clientPair.Priv)
	require.NoError(t, err)

	authInfo, err := swn.GetAuthInfo(deviceId)
	require.NoError(t, err)
	require.NotEmpty(t, authInfo)
}

func TestSaveAuthInfo(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)

	clientPair := genPrivPubKeys(t)

	deviceId, err := swn.SaveAuthInfo(swn.Peer.KeyPair.PubKey, clientPair.Priv)
	require.NoErrorf(t, err, `Error saving AuthInfo: %v`, err)

	authInfo, err := swn.GetAuthInfo(deviceId)
	require.NoErrorf(t, err, `Error getting AuthInfo: %v`, err)

	otherPubkey, err := crypto.UnmarshalEd25519PublicKey(authInfo.OtherDevicePubKey)
	require.NoErrorf(t, err, `Error unmarshalling otherPubkey: %v`, err)

	myDevicePrivateKey, err := crypto.UnmarshalEd25519PrivateKey(authInfo.MyDevicePrivateKey)
	require.NoErrorf(t, err, `Error unmarshalling myDevicePrivateKey: %v`, err)

	if !swn.Peer.KeyPair.PubKey.Equals(otherPubkey) {
		getData, _ := otherPubkey.Raw()
		wantData, _ := swn.Peer.KeyPair.PubKey.Raw()
		t.Fatalf("DevicePub are not equal: got: %v, want: %v", string(getData), string(wantData))
	}

	if !myDevicePrivateKey.Equals(clientPair.Priv) {
		getData, _ := myDevicePrivateKey.Raw()
		wantData, _ := clientPair.Priv.Raw()
		t.Fatalf("ClientPriv are not equal: got: %v, want: %v", string(getData), string(wantData))
	}
}

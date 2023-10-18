package swn_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.neonyx.io/go-swn/internal/swn"
)

func createGetterSender(t *testing.T) (*swn.SWN, *swn.SWN) {
	getter, err := newSWN(1)
	require.NoError(t, err)

	sender, err := newSWN(2)
	require.NoError(t, err)

	err = sender.Peer.EstablishConn(context.Background(), getter.Peer.Getp2pMA())
	require.NoError(t, err)

	return getter, sender
}

func TestIsAuthorized(t *testing.T) {
	getter, sender := createGetterSender(t)
	defer closeSWN(t, getter)
	defer closeSWN(t, sender)

	conns := getter.Peer.Host.Network().Conns()
	require.Equal(t, len(conns), 1)

	require.False(t, sender.IsAuthorized(conns[0]))
}

func TestAuthOut(t *testing.T) {
	getter, sender := createGetterSender(t)
	defer closeSWN(t, getter)
	defer closeSWN(t, sender)

	ack, err := sender.AuthOut(getter.Peer.Getp2pMA().String())
	require.NoError(t, err)
	require.True(t, ack)

	conns := getter.Peer.Host.Network().Conns()
	require.Equal(t, len(conns), 1)
	require.True(t, sender.IsAuthorized(conns[0]))

	// already authenticated
	ack, err = sender.AuthOut(getter.Peer.Getp2pMA().String())
	require.NoError(t, err)
	require.True(t, ack)

	// wrong destination string
	nack, err := sender.AuthOut("/abc/")
	require.Error(t, err)
	require.False(t, nack)
}

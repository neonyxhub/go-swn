package swn_test

import (
	"bufio"
	"bytes"
	"context"
	"sync"
	"testing"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"go.neonyx.io/go-swn/pkg/swn"
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

	require.False(t, sender.IsAuthenticated(conns[0]))
}

func TestAuthOut(t *testing.T) {
	getter, sender := createGetterSender(t)
	defer closeSWN(t, getter)
	defer closeSWN(t, sender)

	ack, err := sender.AuthOut(getter.Peer.Getp2pMA().String())
	require.NoError(t, err)
	require.True(t, ack)

	conns := sender.Peer.Host.Network().Conns()
	require.Equal(t, len(conns), 1)

	// already authenticated
	ack, err = sender.AuthOut(getter.Peer.Getp2pMA().String())
	require.NoError(t, err)
	require.True(t, ack)

	// wrong destination string
	nack, err := sender.AuthOut("/abc/")
	require.Error(t, err)
	require.False(t, nack)
}

func TestAuthIn(t *testing.T) {
	getter, sender := createGetterSender(t)
	defer closeSWN(t, getter)
	defer closeSWN(t, sender)

	var wg sync.WaitGroup
	wg.Add(1)

	stream, err := sender.Peer.Host.NewStream(context.Background(), getter.ID(), swn.HID_AUTH)
	require.NoError(t, err)

	go func(wg *sync.WaitGroup, stream network.Stream) {
		defer wg.Done()
		err := getter.AuthIn(stream)
		require.Error(t, err)
	}(&wg, stream)

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	sender.Log.Info("reading NOK because not authenticated")
	resp, err := swn.ReadB64(rw)
	// TODO: improve this test upon stream reset
	if err != nil {
		require.Error(t, err)
	}
	require.True(t, bytes.Equal([]byte(swn.AUTH_NOK), resp))

	err = swn.WriteB64(rw, []byte{})
	// TODO: improve this test upon stream reset
	if err != nil {
		require.Error(t, err)
	}

	wg.Wait()
}

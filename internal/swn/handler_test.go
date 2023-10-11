package swn_test

import (
	"bufio"
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
	"go.neonyx.io/go-swn/internal/swn"
	"google.golang.org/protobuf/proto"
)

func TestEventHandler(t *testing.T) {
	getter, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, getter)

	sender, err := newSWN(2)
	require.NoError(t, err)
	defer closeSWN(t, sender)

	err = sender.Peer.Host.Connect(context.Background(), peer.AddrInfo{
		ID:    getter.ID(),
		Addrs: getter.Peer.Host.Addrs(),
	})
	require.NoError(t, err)

	s, err := sender.Peer.Host.NewStream(context.Background(), getter.ID(), swn.HID_EVENTBUS)
	require.NoError(t, err)

	evt, _, err := mockEvent(1)
	require.NoError(t, err)

	raw, err := swn.PackEvent(evt)
	require.NoError(t, err)

	s.Write(raw)

	writeLen, err := s.Write(raw)
	require.NoError(t, err)

	require.Equal(t, len(raw), writeLen)

	evt2 := <-getter.GrpcServer.Bus.EventToLocal

	require.True(t, proto.Equal(evt, evt2))
}

func TestEventHandler2(t *testing.T) {
	getter, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, getter)

	sender, err := newSWN(2)
	require.NoError(t, err)
	defer closeSWN(t, sender)

	sender.Peer.EstablishConn(getter.Peer.Getp2pMA())
	conns := sender.Peer.Host.Network().ConnsToPeer(getter.ID())
	require.Equal(t, len(conns), 1)
	stream, err := conns[0].NewStream(context.Background())
	require.NoError(t, err)

	w := bufio.NewWriter(stream)
	w.Write([]byte{0x0a, 0x00})

	go sender.EventHandler(stream)

	select {
	case <-getter.GrpcServer.Bus.EventToLocal:
		t.Fatal("should not receive improper packed event")
	case <-time.After(10 * time.Millisecond):
		require.True(t, true, "should timeout as EventHandler can't process improper event")
	}
}

func TestAuthHandler(t *testing.T) {
	getter, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, getter)

	sender, err := newSWN(2)
	require.NoError(t, err)
	defer closeSWN(t, sender)

	sender.Peer.EstablishConn(getter.Peer.Getp2pMA())
	conns := sender.Peer.Host.Network().ConnsToPeer(getter.ID())
	require.Equal(t, len(conns), 1)
	_, err = conns[0].NewStream(context.Background())
	require.NoError(t, err)
	// TODO:

}

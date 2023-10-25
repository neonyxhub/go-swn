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

	// manually authorize
	getter.AuthDeviceMap[sender.ID().String()] = sender.Device.Id

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

	evt2 := <-getter.GrpcServer.Bus.EventUpstream

	require.True(t, proto.Equal(evt, evt2))
}

func TestEventHandler2(t *testing.T) {
	getter, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, getter)

	sender, err := newSWN(2)
	require.NoError(t, err)
	defer closeSWN(t, sender)

	// manually authorize
	getter.AuthDeviceMap[sender.ID().String()] = sender.Device.Id

	sender.Peer.EstablishConn(context.Background(), getter.Peer.Getp2pMA())
	conns := sender.Peer.Host.Network().ConnsToPeer(getter.ID())
	require.Equal(t, len(conns), 1)
	stream, err := conns[0].NewStream(context.Background())
	require.NoError(t, err)

	w := bufio.NewWriter(stream)
	w.Write([]byte{0x0a, 0x00})

	go sender.EventHandler(stream)

	select {
	case <-getter.GrpcServer.Bus.EventUpstream:
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

	sender.Peer.EstablishConn(context.Background(), getter.Peer.Getp2pMA())
	conns := sender.Peer.Host.Network().ConnsToPeer(getter.ID())
	require.Equal(t, len(conns), 1)

	stream, err := conns[0].NewStream(context.Background())
	require.NoError(t, err)

	go sender.AuthHandler(stream)

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	nack, err := swn.ReadB64(rw)
	require.Error(t, err)
	require.Empty(t, nack)
}

package swn_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/require"

	"go.neonyx.io/go-swn/pkg/swn/config"

	neo_swn "go.neonyx.io/go-swn/pkg/swn"
)

func newSWN(id int, busTimer ...time.Duration) (*neo_swn.SWN, error) {
	var cfg config.Config

	cfg.DataStore.Path = fmt.Sprintf("mock/db/%d", id)
	cfg.GrpcServer.Addr = fmt.Sprintf(":%d", 8090+id)
	cfg.EventBus = config.EVENTBUS_GRPC
	cfg.EventBusTimer = 1 * time.Second
	cfg.P2p.ConnLimit = []int{100, 400}
	cfg.P2p.Multiaddr = "/ip4/0.0.0.0/tcp/0"
	cfg.Log.Dev = true

	if len(busTimer) > 0 {
		cfg.EventBusTimer = busTimer[0]
	}

	opts := []libp2p.Option{}

	if err := os.RemoveAll("mock"); err != nil {
		return nil, err
	}

	swn, err := neo_swn.New(&cfg, opts...)
	if err != nil {
		panic(err)
	}

	if err = swn.Run(); err != nil {
		panic(err)
	}

	return swn, nil
}

func closeSWN(t *testing.T, swn *neo_swn.SWN) {
	err := swn.Stop()
	require.NoError(t, err)
	err = os.RemoveAll("mock")
	require.NoError(t, err)

	require.Empty(t, swn.Peer.Host.Network().ListenAddresses())
}

// tests New(), Run(), Stop() methods of swn
func TestNewRunStop(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)
	require.NotEmpty(t, swn.Peer.Host.Network().ListenAddresses())
}

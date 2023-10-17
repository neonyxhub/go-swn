package swn_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/require"
	neo_swn "go.neonyx.io/go-swn/internal/swn"
	"go.neonyx.io/go-swn/internal/swn/config"
)

func newSWN(id int) (*neo_swn.SWN, error) {
	var cfg config.Config

	cfg.DataStore.Path = fmt.Sprintf("mock/db/%d", id)
	cfg.GrpcServer.Addr = fmt.Sprintf(":%d", 8090+id)
	cfg.P2p.ConnLimit = []int{100, 400}
	cfg.P2p.Multiaddr = "/ip4/0.0.0.0/tcp/0"
	cfg.Log.Dev = true

	opts := []libp2p.Option{}

	swn, err := neo_swn.New(&cfg, opts...)
	if err != nil {
		return nil, err
	}

	swn.Run()

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
	require.NotEmpty(t, swn.GrpcServer.Listener.Addr())
	require.NotEmpty(t, swn.Peer.Host.Network().ListenAddresses())
}

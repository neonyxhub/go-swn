package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.neonyx.io/go-swn/internal/swn/config"
)

func TestParseConfigYaml(t *testing.T) {
	const testConfigYaml string = `
grpc_server:
  addr: :8080
datastore:
  path: mockDatastore
p2p:
  multiaddr: "/ip4/0.0.0.0/tcp/0"
  conn_limit: [100, 400]
log:
  dev: true
`
	data := []byte(testConfigYaml)
	cfg, err := config.ParseConfig(&data)
	require.NoError(t, err)
	require.Equal(t, cfg.GrpcServer.Addr, ":8080")
	require.Equal(t, cfg.DataStore.Path, "mockDatastore")
	require.Equal(t, cfg.P2p.Multiaddr, "/ip4/0.0.0.0/tcp/0")
	require.Equal(t, cfg.Log.Dev, true)
}

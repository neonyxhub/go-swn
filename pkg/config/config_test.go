package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.neonyx.io/go-swn/pkg/config"
)

func TestParseConfigYaml(t *testing.T) {
	const testConfigYaml string = `
grpc_server:
  addr: :50051
eventbus: grpc
eventbus_timer: 1s
datastore:
  path: mockDatastore
p2p:
  multiaddr:
    - "/ip4/0.0.0.0/tcp/65000"
  conn_limit: [100, 400]
log:
  dev: true
`
	data := []byte(testConfigYaml)
	cfg, err := config.ParseConfig(&data)
	require.NoError(t, err)
	require.Equal(t, cfg.GrpcServer.Addr, ":50051")
	require.Equal(t, cfg.DataStore.Path, "mockDatastore")
	require.Equal(t, len(cfg.P2p.Multiaddr), 1)
	require.Equal(t, cfg.Log.Dev, true)
	require.Equal(t, cfg.EventBus, "grpc")
	require.Equal(t, cfg.EventBusTimer, 1*time.Second)
}

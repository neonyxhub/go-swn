package p2p_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.neonyx.io/go-swn/internal/swn/config"
	"go.neonyx.io/go-swn/internal/swn/p2p"
	"go.neonyx.io/go-swn/pkg/logger"
)

func newPeer(t *testing.T, cfg *config.Config) *p2p.Peer {
	logCfg := &logger.LoggerCfg{
		Dev:      cfg.Log.Dev,
		OutPaths: cfg.Log.OutPaths,
		ErrPaths: cfg.Log.ErrPaths,
	}
	log, err := logger.New(logCfg)
	peer, err := p2p.New(cfg, log)
	require.NoError(t, err)
	return peer
}

func stopPeer(t *testing.T, peer *p2p.Peer) {
	peer.Stop()
	require.Empty(t, peer.Host.Network().Peers())
}

func TestNew(t *testing.T) {
	var cfg config.Config
	peer := newPeer(t, &cfg)
	defer stopPeer(t, peer)
}

func TestStop(t *testing.T) {
	var cfg config.Config
	peer := newPeer(t, &cfg)
	stopPeer(t, peer)
}

func TestEstablishConn(t *testing.T) {
	var getterCfg config.Config
	getter := newPeer(t, &getterCfg)
	defer stopPeer(t, getter)

	var senderCfg config.Config
	sender := newPeer(t, &senderCfg)
	defer stopPeer(t, sender)

	getterMultiAddr := getter.Getp2pMA()
	err := sender.EstablishConn(context.Background(), getterMultiAddr)
	require.NoError(t, err)

	senderConns := sender.Host.Network().Conns()
	require.Equal(t, len(senderConns), 1)

	connected := false
	for _, conn := range sender.Host.Network().Conns() {
		if conn.RemotePeer() == getter.Host.ID() {
			connected = true
		}
	}
	require.True(t, connected)
}

func TestGetp2pMA(t *testing.T) {
	var cfg config.Config
	peer := newPeer(t, &cfg)
	defer stopPeer(t, peer)

	maddr := peer.Getp2pMA()
	require.NotEmpty(t, maddr)
	maddrStr := maddr.String()
	parts := strings.Split(maddrStr, "/")
	peerId := parts[len(parts)-1]
	tcpPort := parts[4]

	require.Equal(t, peer.Host.ID().String(), peerId)
	require.Equal(t, maddrStr, fmt.Sprintf("/ip4/%v/tcp/%s/p2p/%s", peer.GetIpv4(), tcpPort, peerId))
}

func TestGetIpv4(t *testing.T) {
	var cfg config.Config
	peer := newPeer(t, &cfg)
	defer stopPeer(t, peer)

	ipv4 := peer.GetIpv4()
	require.NotEqual(t, ipv4, "127.0.0.1")
}

func TestGetTransportPort(t *testing.T) {
	var cfg config.Config
	peer := newPeer(t, &cfg)
	defer stopPeer(t, peer)

	port, err := peer.GetTransportPort("tcp")
	require.NoError(t, err)
	require.NotEmpty(t, port)

	port, err = peer.GetTransportPort("abc")
	require.Error(t, err)
	require.Empty(t, port)
}

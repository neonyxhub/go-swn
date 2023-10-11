package p2p_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.neonyx.io/go-swn/internal/swn/config"
	"go.neonyx.io/go-swn/internal/swn/p2p"
)

func newPeer(t *testing.T, cfg *config.Config) *p2p.Peer {
	peer, err := p2p.New(cfg)
	require.NoError(t, err)
	return peer
}

func stopPeer(t *testing.T, peer *p2p.Peer) {
	peer.Stop()
	require.Empty(t, peer.Host.Network().Peers())
}

func TestNew(t *testing.T) {
	var cfg config.Config
	cfg.P2p.PrivKeyPath = ""
	peer := newPeer(t, &cfg)
	defer stopPeer(t, peer)
	require.True(t, peer.KeyPair.IsGenerated, "should be generated")
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
	err := sender.EstablishConn(getterMultiAddr)
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
	require.Equal(t, maddrStr, fmt.Sprintf("/ip4/127.0.0.1/tcp/%s/p2p/%s", tcpPort, peerId))
}

func TestGetIpv4(t *testing.T) {
	var cfg config.Config
	peer := newPeer(t, &cfg)
	defer stopPeer(t, peer)

	ipv4 := peer.GetIpv4()
	require.NotEqual(t, ipv4, "127.0.0.1")
}

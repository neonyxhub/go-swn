package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	mstream "github.com/multiformats/go-multistream"

	"github.com/multiformats/go-multiaddr"

	"go.neonyx.io/go-swn/internal/swn/config"
	"go.neonyx.io/go-swn/pkg/bus/pb"
	"go.neonyx.io/go-swn/pkg/logger"
)

var (
	ErrNegotioateProtocol = func(args ...interface{}) error {
		return fmt.Errorf("failed to negotiate protocol: %w", args...)
	}

	Log logger.Logger
)

// Eventbus to send and retrieve events, coming from/to libp2p handlers
type Bus struct {
	Sender chan *pb.Event
}

type Peer struct {
	Host host.Host
	Bus  *Bus
}

func New(cfg *config.Config, opts ...libp2p.Option) (*Peer, error) {
	bus := &Bus{
		Sender: make(chan *pb.Event),
	}
	peer := &Peer{
		Bus: bus,
	}

	// prepare multiaddr
	if cfg.P2p.Multiaddr != "" {
		maddr, err := multiaddr.NewMultiaddr(cfg.P2p.Multiaddr)
		if err != nil {
			return nil, err
		}
		opts = append(opts, libp2p.ListenAddrs(maddr))
	}

	// keypair for sign & verify
	peerPrivKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	opts = append(opts, libp2p.Identity(peerPrivKey))

	// transport data between peers is encrypted with TLS
	opts = append(opts, libp2p.Security(libp2ptls.ID, libp2ptls.New))

	if len(cfg.P2p.ConnLimit) >= 2 {
		connMgr, err := connmgr.NewConnManager(
			cfg.P2p.ConnLimit[0], // Lowwater
			cfg.P2p.ConnLimit[1], // HighWater,
			connmgr.WithGracePeriod(1*time.Minute),
		)
		if err != nil {
			return nil, err
		}
		opts = append(opts, libp2p.ConnectionManager(connMgr))
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}
	peer.Host = host

	return peer, nil
}

func (p *Peer) Stop() error {
	return p.Host.Close()
}

func (p *Peer) Pretty(e *pb.Event) string {
	if e == nil {
		return "empty event"
	}
	return fmt.Sprintf(
		"Dest: %s, Lexicon: %s, Data: %s",
		string(e.Dest.GetAddr()),
		string(e.Lexicon.GetUri()),
		string(e.GetData()),
	)
}

// Add to current peer's PeerStore a remote peer by its destination info and ttl
func (p *Peer) AddRemotePeer(destination string, ttl time.Duration) (*peer.AddrInfo, error) {
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		return nil, err
	}

	destInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, err
	}

	for _, existingPeerId := range p.Host.Peerstore().PeersWithAddrs() {
		if existingPeerId == destInfo.ID {
			return destInfo, nil
		}
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	p.Host.Peerstore().AddAddrs(destInfo.ID, destInfo.Addrs, ttl)

	return destInfo, nil
}

func (p *Peer) GetActiveConns(destPeerId peer.ID) []network.Conn {
	active := []network.Conn{}
	for _, conn := range p.Host.Network().ConnsToPeer(destPeerId) {
		if !conn.IsClosed() {
			active = append(active, conn)
		}
	}

	return active
}

// Open connection with peer and adds its info to peerstore
func (p *Peer) EstablishConn(ctx context.Context, maddr multiaddr.Multiaddr) error {
	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	p.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	return p.Host.Connect(ctx, *info)
}

func (p *Peer) StreamOverConn(ctx context.Context, conn network.Conn, protos ...protocol.ID) (network.Stream, error) {
	// creates a new connection
	s, err := conn.NewStream(ctx)
	if err != nil {
		return nil, err
	}

	selected, err := mstream.SelectOneOf(protos, s)
	if err != nil {
		s.Reset()
		return nil, ErrNegotioateProtocol(err)
	}

	err = s.SetProtocol(protocol.ID(selected))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return s, nil
}

// Returns MultiAddr with non-localhost ipv4 and with /p2p/<peerId> prefix
func (p *Peer) Getp2pMA() multiaddr.Multiaddr {
	peerMa, _ := multiaddr.NewMultiaddr("/p2p/" + p.Host.ID().String())

	for _, ma := range p.Host.Addrs() {
		parts := strings.Split(ma.String(), "/")
		if len(parts) < 4 {
			Log.Sugar().Errorf("invalid multiaddr: %v", ma)
			continue
		}
		if parts[1] == "ip4" && parts[2] != "127.0.0.1" {
			return ma.Encapsulate(peerMa)
		}
	}

	return nil
}

// Returns non-localhost IPv4 address of current peer
func (p *Peer) GetIpv4() string {
	ma := p.Getp2pMA()
	parts := strings.Split(ma.String(), "/")

	return parts[2]
}

// Returns current peer's transport port of given protocol
func (p *Peer) GetTransportPort(protocol string) (string, error) {
	var port string
	var proto int

	switch protocol {
	case "tcp":
		proto = multiaddr.P_TCP
	case "udp":
		proto = multiaddr.P_UDP
	case "quic":
		proto = multiaddr.P_QUIC
	default:
		return "", fmt.Errorf("unknown transport protocol: %s", protocol)
	}

	for _, la := range p.Host.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(proto); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		return "", fmt.Errorf("port not found: %s", port)
	}

	return port, nil
}

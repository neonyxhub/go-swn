package p2p

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	mstream "github.com/multiformats/go-multistream"

	"go.neonyx.io/go-swn/internal/swn/config"
	"go.neonyx.io/go-swn/pkg/bus/pb"
	"go.neonyx.io/go-swn/pkg/logger"
)

var (
	ErrNegotioateProtocol = func(args ...interface{}) error {
		return fmt.Errorf("failed to negotiate protocol: %w", args...)
	}
)

// Eventbus to send and retrieve events, coming from/to libp2p handlers
type Bus struct {
	Sender chan *pb.Event
}

type Peer struct {
	Host    host.Host
	Bus     *Bus
	KeyPair *KeyPair25519
	Log     logger.Logger
}

func New(cfg *config.Config, opts ...libp2p.Option) (*Peer, error) {
	bus := &Bus{
		Sender: make(chan *pb.Event, 100),
	}
	keyPair := &KeyPair25519{PrivKeyPath: cfg.P2p.PrivKeyPath}
	peer := &Peer{
		Bus:     bus,
		KeyPair: keyPair,
	}

	// TODO: make sure that there is no duplicate options
	// prepare multiaddr
	if cfg.P2p.Multiaddr != "" {
		maddr, err := multiaddr.NewMultiaddr(cfg.P2p.Multiaddr)
		if err != nil {
			return nil, err
		}
		opts = append(opts, libp2p.ListenAddrs(maddr))
	}

	// prepare private key
	err := keyPair.ReadFromFile()
	if err != nil {
		if os.IsNotExist(err) {
			if err := keyPair.Gen(); err != nil {
				return nil, err
			}
			// TODO: appending should be either keyPair is read from file
			// or generated
			opts = append(opts, libp2p.Identity(peer.KeyPair.PrivKey))
		} else {
			return nil, err
		}
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

// Open connection with peer and adds its info to peerstore
func (p *Peer) EstablishConn(maddr multiaddr.Multiaddr) error {
	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	p.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	return p.Host.Connect(context.Background(), *info)
}

func (p *Peer) StreamOverConn(ctx context.Context, conn network.Conn, protos ...protocol.ID) (network.Stream, error) {
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

func (p *Peer) Getp2pMA() multiaddr.Multiaddr {
	peerMa, _ := multiaddr.NewMultiaddr("/p2p/" + p.Host.ID().String())
	return p.Host.Addrs()[0].Encapsulate(peerMa)
}

func (p *Peer) GetIpv4() string {
	var ipv4 string

	for _, ma := range p.Host.Addrs() {
		parts := strings.Split(ma.String(), "/")
		if parts[1] == "ip4" && parts[2] != "127.0.0.1" {
			return parts[2]
		}
	}

	return ipv4
}

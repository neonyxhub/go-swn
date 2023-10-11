package swn

import (
	"context"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	leveldb_opt "github.com/syndtr/goleveldb/leveldb/opt"

	neo_ds "go.neonyx.io/go-swn/internal/ds"
	"go.neonyx.io/go-swn/internal/ds/drivers"
	"go.neonyx.io/go-swn/internal/swn/config"
	"go.neonyx.io/go-swn/internal/swn/grpc_server"
	"go.neonyx.io/go-swn/internal/swn/p2p"
	"go.neonyx.io/go-swn/pkg/logger"
)

type Handler struct {
	Id   string
	Func func(network.Stream)
}

// Main structure of module with necessary pointers to components
type SWN struct {
	Cfg        *config.Config
	Ds         drivers.DataStore
	DsCfg      *drivers.DataStoreCfg
	GrpcServer *grpc_server.GrpcServer
	Peer       *p2p.Peer
	Handlers   []*Handler
	Log        logger.Logger
	Ctx        context.Context
	CtxCancel  context.CancelFunc
}

// New creates an instance of SWN with libp2p peer, datastore, gRPC server, P2PBus
func New(cfg *config.Config, opts ...libp2p.Option) (*SWN, error) {
	ctx, cancel := context.WithCancel(context.Background())

	swn := SWN{
		Cfg:       cfg,
		Ctx:       ctx,
		CtxCancel: cancel,
	}

	logCfg := &logger.LoggerCfg{
		Dev:      cfg.Log.Dev,
		OutPaths: cfg.Log.OutPaths,
		ErrPaths: cfg.Log.ErrPaths,
	}
	log, err := logger.New(logCfg)
	if err != nil {
		return nil, err
	}
	swn.Log = log

	// new libp2p peer
	peer, err := p2p.New(cfg, opts...)
	if err != nil {
		return nil, err
	}
	peer.Log = log
	swn.Peer = peer

	// new DataStore driver
	swn.DsCfg = &drivers.DataStoreCfg{
		Path:    cfg.DataStore.Path,
		Options: leveldb_opt.Options{},
	}
	driver, err := neo_ds.New(swn.DsCfg)
	if err != nil {
		return nil, err
	}
	swn.Ds = driver

	swn.ApplyDefaultHandlers()

	swn.GrpcServer = grpc_server.New()
	swn.GrpcServer.Log = log
	swn.GrpcServer.Bus.PeerId = []byte(swn.ID())

	return &swn, nil
}

// Serves gRPC server, set p2p network stream handlers and starts event listening
func (s *SWN) Run() error {
	s.Log.Sugar().Infof("starting gRPC server on %s", s.Cfg.GrpcServer.Addr)
	if err := s.GrpcServer.Serve(s.Cfg.GrpcServer.Addr); err != nil {
		s.Ds.Close()
		return err
	}

	s.Log.Sugar().Infof("starting %d handlers", len(s.Handlers))
	for _, h := range s.Handlers {
		s.Peer.Host.SetStreamHandler(protocol.ID(h.Id), func(stream network.Stream) {
			h.Func(stream)
		})
	}

	for _, p := range s.Peer.Host.Mux().Protocols() {
		s.Log.Sugar().Infof("protocol: %v", p)
	}

	if err := s.StartEventListening(); err != nil {
		s.Ds.Close()
		return err
	}

	if s.Cfg.Debug {
		err := s.DebugSavePeerInfo()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SWN) Stop() error {
	if s.Cfg.Debug {
		err := s.DebugDeletePeerInfo()
		if err != nil {
			return err
		}
	}

	s.Log.Sugar().Infof("stopping gRPC server on %s", s.Cfg.GrpcServer.Addr)
	if err := s.GrpcServer.Stop(); err != nil {
		return err
	}

	for _, handler := range s.Handlers {
		s.Peer.Host.RemoveStreamHandler(protocol.ID(handler.Id))
	}

	s.StopEventListening()

	if err := s.Peer.Stop(); err != nil {
		return err
	}

	if err := s.Ds.Close(); err != nil {
		return err
	}

	return nil
}

func (s *SWN) GetPeerTransportPort(p string) (string, error) {
	var port string
	var proto int

	switch p {
	case "tcp":
		proto = multiaddr.P_TCP
	case "udp":
		proto = multiaddr.P_UDP
	case "quic":
		proto = multiaddr.P_QUIC
	default:
		return "", fmt.Errorf("unknown transport protocol: %s", p)
	}

	for _, la := range s.Peer.Host.Network().ListenAddresses() {
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

func (s *SWN) ID() peer.ID {
	return s.Peer.Host.ID()
}

func (s *SWN) GetPeerMAddrs() string {
	maddrs := []string{}
	for _, addr := range s.Peer.Host.Addrs() {
		maddrs = append(maddrs, addr.String())
	}
	return strings.Join(maddrs, ",")
}

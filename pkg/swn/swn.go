package swn

import (
	"context"

	"github.com/go-errors/errors"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"go.neonyx.io/go-swn/pkg/bus"
	"go.neonyx.io/go-swn/pkg/config"
	"go.neonyx.io/go-swn/pkg/ds"
	"go.neonyx.io/go-swn/pkg/ds/drivers"
	"go.neonyx.io/go-swn/pkg/logger"
	"go.neonyx.io/go-swn/pkg/swn/p2p"

	"go.neonyx.io/go-swn/internal/grpcserver"
	"go.neonyx.io/go-swn/internal/natsclient"
)

type Handler struct {
	Id   string
	Func network.StreamHandler
}

// Main structure of module with necessary pointers to components
type SWN struct {
	// swn configuration for gRPC, p2p, logger etc.
	Cfg *config.Config

	Log logger.Logger

	// local datastore interface and configuration
	Ds    drivers.DataStore
	DsCfg *drivers.DataStoreCfg

	// implementation for sending and receiving Event in p2p network
	EventBus bus.EventBus

	// internal main structure for event IO channels
	EventIO *bus.EventIO

	// peer structure with p2p logic
	Peer *p2p.Peer

	// local swn "hardware" information like deviceId
	Device *Device

	// p2p remote peer Id: remote peer's DeviceId
	AuthDeviceMap map[string][]byte

	// slice of p2p stream handlers
	Handlers []Handler

	// parent context of swn state with cancel function
	Ctx       context.Context
	CtxCancel context.CancelFunc
}

// New creates an instance of SWN with libp2p peer, datastore, gRPC server, P2PBus
func New(cfg *config.Config, opts ...libp2p.Option) (*SWN, error) {
	logCfg := &logger.LoggerCfg{
		Name:     cfg.Log.Name,
		Dev:      cfg.Log.Dev,
		OutPaths: cfg.Log.OutPaths,
		ErrPaths: cfg.Log.ErrPaths,
	}
	log, err := logger.New(logCfg)
	if err != nil {
		return nil, err
	}

	log.Info("creating a new SWN")
	ctx, cancel := context.WithCancel(context.Background())

	swn := SWN{
		Cfg:           cfg,
		Ctx:           ctx,
		CtxCancel:     cancel,
		AuthDeviceMap: make(map[string][]byte),
		Log:           log,
	}

	// new libp2p peer
	peer, err := p2p.New(cfg, log, opts...)
	if err != nil {
		return nil, err
	}
	swn.Peer = peer

	// new DataStore driver
	swn.DsCfg = &drivers.DataStoreCfg{
		Path:    cfg.DataStore.Path,
		Options: opt.Options{},
	}
	driver, err := ds.New(swn.DsCfg)
	if err != nil {
		return nil, err
	}
	swn.Ds = driver

	// init device
	swn.Device = &Device{}
	if err = swn.CheckDeviceId(); err != nil {
		return nil, err
	}

	// main EventBus channels management
	timeout := cfg.EventBusTimer
	swn.EventIO = bus.New(timeout, timeout)

	// how to handle downstreaming and upstreaming Event in p2p
	switch cfg.EventBus {
	case config.EVENTBUS_EVENTIO:
		log.Info("set EvenIO eventbus")
		swn.EventBus = swn.EventIO

	case config.EVENTBUS_NATS:
		log.Info("set NATS eventbus")
		natsClient, err := natsclient.New(cfg.Nats.Url, swn.EventIO, log)
		if err != nil {
			return nil, err
		}
		swn.EventBus = natsClient

	case config.EVENTBUS_GRPC:
		log.Info("set gRPC eventbus")
		// init internal gRPC management
		grpcServer := grpcserver.New(cfg, swn.EventIO, log)
		grpcServer.PeerId = []byte(swn.ID())
		//swn.GrpcServer = grpcServer
		swn.EventBus = grpcServer

	default:
		return nil, errors.Errorf("unknown eventbus: %v", cfg.EventBus)
	}

	swn.ApplyDefaultHandlers()

	return &swn, nil
}

// Starts eventbus, set p2p network stream handlers and starts event listening
func (s *SWN) Run() error {
	s.Log.Sugar().Infof("starting eventbus as %s", s.Cfg.EventBus)
	if err := s.EventBus.Run(); err != nil {
		if err := s.Ds.Close(); err != nil {
			return err
		}
		return err
	}

	s.Log.Sugar().Infof("starting %d handlers", len(s.Handlers))
	for _, h := range s.Handlers {
		s.Peer.Host.SetStreamHandler(protocol.ID(h.Id), h.Func)
	}

	for _, p := range s.Peer.Host.Mux().Protocols() {
		s.Log.Sugar().Infof("protocol: %v", p)
	}

	if err := s.StartEventListening(); err != nil {
		if err := s.Ds.Close(); err != nil {
			return err
		}
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

	s.Log.Sugar().Infof("stopping eventbus %s", s.Cfg.EventBus)
	if err := s.EventBus.Stop(); err != nil {
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

	if err := s.EventBus.Stop(); err != nil {
		return err
	}

	return nil
}

func (s *SWN) ID() peer.ID {
	return s.Peer.Host.ID()
}

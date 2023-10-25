package grpc_server

import (
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"

	"go.neonyx.io/go-swn/internal/swn/config"
	"go.neonyx.io/go-swn/pkg/bus/pb"
	logger "go.neonyx.io/go-swn/pkg/logger"
)

const (
	GRPC_TIMEOUT = 10 * time.Second
)

var (
	ErrGrpcTimeout = errors.New("gRPC server timeout")

	Log logger.Logger
)

type GrpcServer struct {
	Server   *grpc.Server
	Listener net.Listener
	Bus      *SWNBusServer
}

// New creates a new gRPC server and register a new SWNBusServer service
func New(cfg *config.Config) *GrpcServer {
	s := grpc.NewServer()

	sBus := &SWNBusServer{
		EventDownstream: make(chan *pb.Event),
		EventUpstream:   make(chan *pb.Event),
	}

	pb.RegisterSWNBusServer(s, sBus)

	return &GrpcServer{Server: s, Bus: sBus}
}

// Serve serves gRPC connection on already registered services.
// This function is blocking until server is stopped
func (s *GrpcServer) Serve(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.Listener = listener

	go func() {
		if err := s.Server.Serve(listener); err != nil {
			Log.Sugar().Errorf("gRPC Serve() failed: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops gRPC server letting active connections to complete
func (s *GrpcServer) Stop() error {
	done := make(chan bool, 1)

	go func() {
		s.Server.GracefulStop()
		close(done)
	}()

	select {
	case <-time.After(GRPC_TIMEOUT):
		return ErrGrpcTimeout
	case <-done:
		return nil
	}
}

// Get gRPC server port
func (s *GrpcServer) GetPort() int {
	return s.Listener.Addr().(*net.TCPAddr).Port
}

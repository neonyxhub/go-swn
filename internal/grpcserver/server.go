package grpcserver

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.neonyx.io/go-swn/pkg/bus"
	"go.neonyx.io/go-swn/pkg/config"
	"go.neonyx.io/go-swn/pkg/logger"

	"go.neonyx.io/go-swn/pkg/bus/pb"
)

const (
	GRPC_TIMEOUT = 10 * time.Second
)

var (
	ErrGrpcTimeout     = errors.New("gRPC server timeout")
	ErrNoLocalListener = errors.New("no local listener is presented")
)

type swnBusServer struct {
	pb.UnimplementedSWNBusServer
}

type GrpcServer struct {
	*grpc.Server
	net.Listener

	bus        *swnBusServer
	eventIOPtr *bus.EventIO
	Log        logger.Logger
	Cfg        *config.Config
	PeerId     []byte
}

// New creates a new gRPC server and register a new GrpcServer service
func New(cfg *config.Config, eventIO *bus.EventIO, log logger.Logger) *GrpcServer {
	s := grpc.NewServer()

	sBus := &swnBusServer{}
	pb.RegisterSWNBusServer(s, sBus)

	return &GrpcServer{
		Server:     s,
		eventIOPtr: eventIO,
		Log:        log,
		Cfg:        cfg,
		bus:        sBus,
	}
}

// Serve serves gRPC connection on already registered services.
// This function is blocking until server is stopped
func (s *GrpcServer) Run() error {
	listener, err := net.Listen("tcp", s.Cfg.GrpcServer.Addr)
	if err != nil {
		return err
	}

	s.Listener = listener

	go func() {
		if err := s.Serve(listener); err != nil {
			s.Log.Sugar().Errorf("gRPC Serve() failed: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops gRPC server letting active connections to complete
func (s *GrpcServer) Stop() error {
	done := make(chan bool, 1)

	go func() {
		s.GracefulStop()
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

// Send Event from SWN to gRPC LocalFunnelEvents call
// Since this eventbus impl. is gRPC server, then we should receive gRPC call
// first to actually send upstream event
func (s *GrpcServer) SendUpstream(event *pb.Event) error {
	return s.eventIOPtr.SendUpstream(event)
}

// TODO: add gRPC status and error handling

// Gets events from local sender stream and passes then to EventIO Downstream channel
func (s *GrpcServer) LocalDistributeEvents(stream pb.SWNBus_LocalDistributeEventsServer) error {
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamEventsResponse{})
		}
		if err != nil {
			return err
		}
		if err := s.eventIOPtr.RecvDownstream(stream.Context(), event); err != nil {
			return err
		}
	}
}

// Gets events from EventUpstream channel and passes them to local listener
func (s *GrpcServer) LocalFunnelEvents(in *pb.ListenEventsRequest, stream pb.SWNBus_LocalFunnelEventsServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event := <-s.eventIOPtr.Upstream:
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}

func (s *GrpcServer) GetPeerInfo(ctx context.Context, in *emptypb.Empty) (*pb.Peer, error) {
	defer ctx.Done()
	return &pb.Peer{Id: s.PeerId}, nil
}

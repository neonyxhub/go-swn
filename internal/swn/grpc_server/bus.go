package grpc_server

import (
	"context"

	"io"
	"sync"

	"github.com/go-errors/errors"
	"go.neonyx.io/go-swn/pkg/bus/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrNoLocalListener = errors.New("no local listener is presented")
)

type SWNBusServer struct {
	pb.UnimplementedSWNBusServer
	sync.Mutex

	EventDownstream  chan *pb.Event
	EventUpstream    chan *pb.Event
	EventUpstreamBuf []*pb.Event

	PeerId []byte
}

func (sBus *SWNBusServer) FlushUpstreamBuffer() {
	sBus.Lock()
	bufferCopy := make([]*pb.Event, len(sBus.EventUpstreamBuf))
	copy(bufferCopy, sBus.EventUpstreamBuf)
	sBus.EventUpstreamBuf = nil
	sBus.Unlock()

	for _, event := range bufferCopy {
		sBus.EventUpstream <- event
	}
}

// Gets events from local sender stream and passes then to EventDownstream channel
func (sBus *SWNBusServer) LocalDistributeEvents(stream pb.SWNBus_LocalDistributeEventsServer) error {
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamEventsResponse{})
		}
		if err != nil {
			return err
		}
		sBus.EventDownstream <- event
	}
}

// Gets events from EventUpstream channel and passes them to local listener
func (sBus *SWNBusServer) LocalFunnelEvents(in *pb.ListenEventsRequest, srv pb.SWNBus_LocalFunnelEventsServer) error {
	// upon next listener, flush upstream buffer
	go sBus.FlushUpstreamBuffer()

	for {
		select {
		case <-srv.Context().Done():
			return nil
		case event := <-sBus.EventUpstream:
			if err := srv.Send(event); err != nil {
				return err
			}
		}
	}
}

func (sBus *SWNBusServer) GetPeerInfo(ctx context.Context, in *emptypb.Empty) (*pb.Peer, error) {
	defer ctx.Done()
	return &pb.Peer{Id: sBus.PeerId}, nil
}

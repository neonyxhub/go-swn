package grpc_server

import (
	"context"
	"io"
	"log"

	"go.neonyx.io/go-swn/pkg/bus/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SWNBusServer struct {
	pb.UnimplementedSWNBusServer
	EventFromLocal chan *pb.Event
	EventToLocal   chan *pb.Event
	PeerId         []byte

	HasListener bool
}

// Gets events from local sender stream and passes then to EventFromLocal channel
func (sBus *SWNBusServer) LocalDistributeEvents(stream pb.SWNBus_LocalDistributeEventsServer) error {
	log.Printf("start listening to stream events")
	for {
		event, err := stream.Recv()
		sBus.EventFromLocal <- event
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamEventsResponse{})
		}
		if err != nil {
			return err
		}
		log.Printf("dest: %v, method: %v, data: %v\n",
			string(event.Dest.GetAddr()), event.Lexicon.GetUri(), string(event.GetData()))
	}
}

// Gets events from EventToLocal channel and passes them to local listener
func (sBus *SWNBusServer) LocalFunnelEvents(in *pb.ListenEventsRequest, srv pb.SWNBus_LocalFunnelEventsServer) error {
	sBus.HasListener = true

	defer func() { sBus.HasListener = false }()

	for {
		select {
		case <-srv.Context().Done():
			return nil
		case event := <-sBus.EventToLocal:
			if err := srv.Send(event); err != nil {
				log.Printf("send error %v", err)
			}
		}
	}
}

func (sBus *SWNBusServer) GetPeerId(ctx context.Context, in *emptypb.Empty) (*pb.Peer, error) {
	defer ctx.Done()
	return &pb.Peer{Id: sBus.PeerId}, nil
}

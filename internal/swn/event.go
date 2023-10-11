package swn

import (
	"context"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"

	"go.neonyx.io/go-swn/pkg/bus/pb"
)

var (
	ErrNoLocalListener       = errors.New("no local listener is presented")
	ErrNoExistingConnection  = errors.New("no existing connection is presented")
	ErrIncompleteStreamWrite = errors.New("incomplete stream write")
	ErrEmptyEvent            = errors.New("empty event")
)

// Passes event to listeners, connected over grpc method
func (s *SWN) EventToLocalListener(event *pb.Event) error {
	s.Log.Sugar().Infoln("passing event to local listener")

	if !s.GrpcServer.Bus.HasListener {
		s.Log.Sugar().Infoln("local listener is not detected")
		s.GrpcServer.Bus.EventToLocal <- event
		// TODO: store event to storage
		return ErrNoLocalListener
	}

	s.GrpcServer.Bus.EventToLocal <- event

	s.Log.Sugar().Infoln("wrote event to channel")

	return nil
}

// Listens for incoming requests and passes them to the network
func (s *SWN) StartEventListening() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	s.Log.Sugar().Infof("start event listening on peer: %v", s.GetPeerMAddrs())
	go func(s *SWN) {
		for {
			select {
			case <-s.Ctx.Done():
				return
			case evt := <-s.GrpcServer.Bus.EventFromLocal:
				s.Log.Sugar().Infof("got event to pass: %v", s.Peer.Pretty(evt))
				err := s.PassEventToNetwork(evt)
				if err != nil {
					s.Log.Sugar().Errorf("error passing event to network: %v", err)
				}
			}
		}
	}(s)

	return
}

func (s *SWN) StopEventListening() {
	s.CtxCancel()
}

// Gets event, resolves destination and passes it with libp2p connection
func (s *SWN) PassEventToNetwork(evt *pb.Event) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if evt == nil {
		s.Log.Sugar().Errorln(ErrEmptyEvent)
		return ErrEmptyEvent
	}

	ma, err := multiaddr.NewMultiaddrBytes(evt.Dest.GetAddr())
	if err != nil {
		return err
	}

	info, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		return err
	}

	err = s.Peer.EstablishConn(ma)
	if err != nil {
		return err
	}

	conns := s.Peer.Host.Network().ConnsToPeer(info.ID)
	if len(conns) != 0 {
		// s, err := s.Host.NewStream(ctx, info.ID, HID_EVENTBUS)
		// if err != nil {
		// 	return err
		// }
		// return StreamPassEvent(ctx, evt, s)
		return s.ConnPassEvent(ctx, evt, conns[0])
	}

	return ErrNoExistingConnection
}

// Packs event and passes it over an existing stream
func StreamPassEvent(ctx context.Context, evt *pb.Event, s network.Stream) error {
	rawEvt, err := PackEvent(evt)
	if err != nil {
		return err
	}

	n, err := s.Write(rawEvt)
	if n != len(rawEvt) {
		return ErrIncompleteStreamWrite
	}

	return err
}

// Packs event, opens a stream over an existing connection and writes packed event
func (s *SWN) ConnPassEvent(ctx context.Context, evt *pb.Event, conn network.Conn) error {
	rawEvt, err := PackEvent(evt)
	if err != nil {
		return err
	}

	for _, s := range conn.GetStreams() {
		if s.Protocol() == HID_EVENTBUS {
			n, err := s.Write(rawEvt)
			if n != len(rawEvt) {
				return ErrIncompleteStreamWrite
			}

			return err
		}
	}

	// creates a new stream and pass event to HID_EVENTBUS
	stream, err := s.Peer.StreamOverConn(ctx, conn, HID_EVENTBUS)
	if err != nil {
		return err
	}

	n, err := stream.Write(rawEvt)
	if n != len(rawEvt) {
		return ErrIncompleteStreamWrite
	}

	return err
}

func PackEvent(evt *pb.Event) ([]byte, error) {
	rawEvt, err := proto.Marshal(evt)
	if err != nil {
		return nil, err
	}
	l := len(rawEvt)

	return append([]byte(fmt.Sprintf("%v\n", l)), rawEvt...), nil
}

package swn

import (
	"bufio"

	"github.com/libp2p/go-libp2p/core/network"
)

const (
	HID_AUTH     = "/swn/auth/1.0.0"
	HID_EVENTBUS = "/swn/eventbus/1.0.0"
)

func (s *SWN) ApplyDefaultHandlers() {
	s.Handlers = []Handler{
		{
			Id:   HID_AUTH,
			Func: s.AuthHandler,
		},
		{
			Id:   HID_EVENTBUS,
			Func: s.EventHandler,
		},
	}
}

func (s *SWN) RegisterNewHandler(h ...Handler) {
	s.Handlers = append(s.Handlers, h...)
}

// Handle incoming auth from another SWN via HID_AUTH protocol
func (s *SWN) AuthHandler(stream network.Stream) {
	s.Log.Sugar().Infof("got auth stream: connId=%v", stream.Conn().ID())

	if err := s.AuthIn(stream); err != nil {
		s.Log.Sugar().Errorln(err)
		s.Log.Sugar().Warnf("closing stream for failed auth stream: %s", stream.Conn().ID())
		stream.Conn().Close()
	}
}

// Handle protobuf Events via HID_EVENTBUS protocol
func (s *SWN) EventHandler(stream network.Stream) {
	s.Log.Sugar().Infof("got event stream: connId=%v", stream.Conn().ID())

	if !s.IsAuthorized(stream.Conn()) {
		s.Log.Sugar().Warnf("closing stream for unauthorized connection: %s", stream.Conn().ID())
		stream.Conn().Close()
		return
	}

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	evt, err := UnpackEvent(rw)
	if err != nil {
		s.Log.Sugar().Errorf("failed to UnpackEvent: %v", err)
		return
	}

	s.Log.Sugar().Infof("got event: %s", s.Peer.Pretty(evt))

	s.EventToLocalListener(evt)
}

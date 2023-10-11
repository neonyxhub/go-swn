package swn

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p/core/network"
	"go.neonyx.io/go-swn/pkg/bus/pb"
	"google.golang.org/protobuf/proto"
)

const (
	HID_AUTH     = "/swn/auth/1.0.0"
	HID_EVENTBUS = "/swn/eventbus/1.0.0"
)

func (s *SWN) ApplyDefaultHandlers() {
	s.Handlers = []*Handler{
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

func (s *SWN) RegisterNewHandler(h ...*Handler) {
	s.Handlers = append(s.Handlers, h...)
}

// Authorize to SWN
func (s *SWN) AuthHandler(stream network.Stream) {
	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	err := s.PerformIncomingAuth(rw)
	if err != nil {
		s.Log.Sugar().Errorln(err)
		s.Log.Sugar().Warnf("closing stream for failed auth stream: %s", stream.Conn().ID())
		stream.Conn().Close()
	}
}

// Regular libp2p handler, to handle Events coming into HID_EVENTBUS protocol
func (s *SWN) EventHandler(stream network.Stream) {
	s.Log.Info("got event stream")
	r := bufio.NewReader(stream)
	for {
		// TODO: implement another way of reading packed event bytes instead of \n
		lenEvt, err := r.ReadBytes('\n')
		if err != nil {
			if err == network.ErrReset {
				// TODO: handle ErrReset properly
				s.Log.Sugar().Errorln(err)
				break
			}

			s.Log.Sugar().Errorf("error reading event length: %v", err)
			continue
		}

		lenEvtInt, err := strconv.Atoi(strings.Trim(string(lenEvt), "\n"))
		if err != nil {
			s.Log.Sugar().Errorf("event length is not number: %v", err)
			continue
		}

		rawEvt := make([]byte, lenEvtInt)
		_, err = r.Read(rawEvt)
		if err != nil {
			s.Log.Sugar().Errorf("error reading event: %v", err)
			continue
		}

		evt := &pb.Event{}
		err = proto.Unmarshal(rawEvt, evt)
		if err != nil {
			s.Log.Sugar().Errorf("error unmarshalling event: %v", err)
			continue
		}

		s.Log.Sugar().Infof("got event: %s", s.Peer.Pretty(evt))

		s.EventToLocalListener(evt)
	}
}

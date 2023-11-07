package bus

import (
	"context"
	"time"

	"github.com/go-errors/errors"
	"go.neonyx.io/go-swn/pkg/bus/pb"
)

var (
	ErrRecvTimeout = errors.Errorf("EventIO.Recv() timeout")
	ErrSendTimeout = errors.Errorf("EventIO.Send() timeout")
)

type EventBus interface {
	ProduceUpstream(event *pb.Event) error
}

// Generic struct to send and recv Event from channels
type EventIO struct {
	dnTimeout time.Duration
	upTimeout time.Duration

	// chan<- send only, <-chan recv only
	Downstream     chan<- *pb.Event
	Upstream       chan *pb.Event
	UpstreamBufCnt int
}

func New(up chan *pb.Event, dn chan<- *pb.Event, dnTimeout, upTimeout time.Duration) *EventIO {
	return &EventIO{
		dnTimeout:  dnTimeout,
		upTimeout:  upTimeout,
		Downstream: dn,
		Upstream:   up,
	}
}

// Recv Event from client to SWN
func (e *EventIO) Recv(ctx context.Context, event *pb.Event) error {
	select {
	case <-ctx.Done():
		return nil
	case e.Downstream <- event:
		return nil
	case <-time.After(e.dnTimeout):
		return ErrRecvTimeout
	}
}

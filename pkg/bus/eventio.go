package bus

import (
	"context"
	"time"

	"github.com/go-errors/errors"
	"go.neonyx.io/go-swn/pkg/bus/pb"
)

var (
	ErrRecvTimeout = errors.Errorf("EventIO recv timeout")
	ErrSendTimeout = errors.Errorf("EventIO send timeout")
)

// Generic struct to send and recv Event from channels
type EventIO struct {
	dnTimeout time.Duration
	upTimeout time.Duration

	Downstream     chan *pb.Event
	Upstream       chan *pb.Event
	UpstreamBufCnt int
}

func New(dnTimeout, upTimeout time.Duration) *EventIO {
	return &EventIO{
		dnTimeout: dnTimeout,
		upTimeout: upTimeout,
		// TODO: add buf size
		Downstream: make(chan *pb.Event),
		Upstream:   make(chan *pb.Event),
	}
}

// Recv Event from client to SWN
func (e *EventIO) RecvDownstream(ctx context.Context, event *pb.Event) error {
	select {
	case <-ctx.Done():
		return nil
	case e.Downstream <- event:
		return nil
	case <-time.After(e.dnTimeout):
		return ErrRecvTimeout
	}
}

// Send Event SWN from client
func (e *EventIO) SendUpstream(event *pb.Event) error {
	select {
	case e.Upstream <- event:
		return nil
	// TODO: use sync.Pool as *time.Timer to optimize timer GC
	case <-time.After(e.upTimeout):
		e.UpstreamBufCnt++
		return ErrSendTimeout
	}
}

// Nothing to run
func (e *EventIO) Run() error {
	return nil
}

// Nothing to stop
func (e *EventIO) Stop() error {
	return nil
}

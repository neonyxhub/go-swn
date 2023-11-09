package bus

import (
	"go.neonyx.io/go-swn/pkg/bus/pb"
)

// Interface to handle downstreaming and upstreaming Event in p2p
// Implementations should be able handle sending upstream Event.
// Downstream event should be handled in own receiving way, but writing to
// internal EventIO Downstream channel.
// EventIO implementation of SendUpstream methods can be
// called if no pre-processing is required
type EventBus interface {
	SendUpstream(event *pb.Event) error
	Stop() error
}

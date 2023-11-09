package natsclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"go.neonyx.io/go-swn/pkg/bus"
	"go.neonyx.io/go-swn/pkg/bus/pb"
	"go.neonyx.io/go-swn/pkg/logger"
)

const (
	TOPIC_MODULE_RESP = "module.resp"
)

var (
	ErrInvalidEventLexicon = errors.Errorf("invalid Event lexicon")
)

type NatsClient struct {
	nc         *nats.Conn
	eventIOPtr *bus.EventIO

	Subs []*nats.Subscription

	Log logger.Logger
}

func New(url string, eventIO *bus.EventIO, logger logger.Logger) (*NatsClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	// TODO: implement protobuf encoder?
	//nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	natsClient := &NatsClient{
		nc:         nc,
		eventIOPtr: eventIO,
		Log:        logger,
	}

	sub, err := nc.Subscribe(TOPIC_MODULE_RESP, natsClient.ModuleRespHandler)
	if err != nil {
		return nil, err
	}

	natsClient.Subs = append(natsClient.Subs, sub)

	return natsClient, nil
}

// Subscribe async for modules response subject and pass over p2p network
func (n *NatsClient) ModuleRespHandler(m *nats.Msg) {
	n.Log.Sugar().Infof("received module.resp: %v", m)

	event := &pb.Event{}
	if err := proto.Unmarshal(m.Data, event); err != nil {
		n.Log.Sugar().Errorf("failed to unmarshal module.resp: %v", err)
		return
	}

	n.Log.Sugar().Infof("received module.resp: %v", m)

	if err := n.eventIOPtr.RecvDownstream(context.Background(), event); err != nil {
		n.Log.Sugar().Errorf("failed to send to SWN upon module.resp: %v", m)
	}
}

// Publish Event to module
func (n *NatsClient) SendUpstream(event *pb.Event) error {
	var moduleSubj string

	uri := event.Lexicon.GetUri()
	parts := strings.Split(uri, "/")
	if len(parts) != 2 {
		return ErrInvalidEventLexicon
	}

	// wildcard
	moduleSubj = fmt.Sprintf("%s.>", parts[1])

	eventRaw, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	n.Log.Sugar().Infof("publishing Event to %s", moduleSubj)

	return n.nc.Publish(moduleSubj, eventRaw)
}

func (n *NatsClient) Stop() error {
	n.Log.Info("unsubscribing from NATS")

	for _, sub := range n.Subs {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

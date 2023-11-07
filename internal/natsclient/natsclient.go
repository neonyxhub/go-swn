package natsclient

import (
	"github.com/nats-io/nats.go"
	"go.neonyx.io/go-swn/pkg/bus/pb"
)

type NatsClient struct {
	conn *nats.Conn
}

func New(url string) (*NatsClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsClient{conn: conn}, nil
}

func (n *NatsClient) ProduceUpstream(event *pb.Event) error {
	return nil
}

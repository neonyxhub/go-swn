package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/multiformats/go-multiaddr"
	"go.neonyx.io/go-swn/pkg/bus/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type debugPeerInfo struct {
	GrpcServerPort int    `yaml:"grpc_server_port"`
	PeerIpv4       string `yaml:"peer_ipv4"`
	PeerId         string `yaml:"peer_id"`
	TransportPort  string `yaml:"transport_port"`
}

type debugPeers struct {
	Peers []debugPeerInfo
}

var (
	swn1, swn2 debugPeerInfo
)

func fetchPeers() {
	p := &debugPeers{}
	const debugYml = "/testdata/debug.yml"

	for retry := 0; retry < 3; retry++ {
		time.Sleep(1 * time.Second)

		if _, err := os.Stat(debugYml); os.IsNotExist(err) {
			log.Println("debug.yml is missing")
			continue
		}

		data, err := os.ReadFile(debugYml)
		if err != nil {
			log.Fatalf("failed to read debug.yml: %v", err)
		}

		err = yaml.Unmarshal(data, p)
		if err != nil {
			log.Fatalf("failed to unmarshal debug.yml: %v", err)
		}

		if len(p.Peers) != 2 {
			log.Println("should be 2 peers (swn1, swn2) in debug.yml")
			continue
		}

		break
	}

	swn1 = p.Peers[0]
	swn2 = p.Peers[1]
}

func consumer(swnclient2 pb.SWNBusClient, done chan bool) {
	responseStream, err := swnclient2.LocalFunnelEvents(context.Background(), &pb.ListenEventsRequest{})
	if err != nil {
		log.Fatalf("failed to call LocalFunnelEvents: %v", err)
	}

	log.Println("waiting on consumer")
	event, err := responseStream.Recv()
	if err != nil {
		log.Fatalf("failed to receive from stream: %v", err)
	}

	log.Printf("received event: Dest: %v", event.Dest.GetAddr())

	done <- true
}

func main() {
	fetchPeers()

	log.Printf("swn1: %s", swn1.PeerId)
	log.Printf("swn2: %s", swn2.PeerId)

	// connect to swn1 gRPC
	swn1Addr := fmt.Sprintf("%s:%d", swn1.PeerIpv4, swn1.GrpcServerPort)
	conn1, err := grpc.Dial(swn1Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to swn1 gRPC: %v", err)
	}
	defer conn1.Close()

	swnclient1 := pb.NewSWNBusClient(conn1)
	log.Println("connected to swn1 gRPC")

	// connect to swn2 gRPC
	swn2Addr := fmt.Sprintf("%s:%d", swn2.PeerIpv4, swn2.GrpcServerPort)
	conn2, err := grpc.Dial(swn2Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to swn2 gRPC: %v", err)
	}
	defer conn2.Close()

	swnclient2 := pb.NewSWNBusClient(conn2)
	log.Println("connected to swn2 gRPC")

	done := make(chan bool, 1)

	go consumer(swnclient2, done)

	// [cwn1 -> swn1] -> swn2
	stream, err := swnclient1.LocalDistributeEvents(context.Background())
	if err != nil {
		log.Fatalf("failed to create stream from LocalDistributeEvents: %v", err)
	}

	swn2MultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%v/tcp/%v/p2p/%v", swn2.PeerIpv4, swn2.TransportPort, swn2.PeerId))
	if err != nil {
		log.Fatalf("failed to create multiaddr for swn2: %v", err)
	}

	event := &pb.Event{
		Id:      1234,
		Type:    pb.EventType_REQ,
		Dest:    &pb.Destination{Addr: swn2MultiAddr.Bytes()},
		Lexicon: &pb.LexiconUri{Uri: "/chat/message/send"},
		Data:    []byte{0x1},
	}

	if err := stream.Send(event); err != nil {
		log.Fatalf("failed to send event: %v", err)
	}

	log.Println("closing consumer")
	<-done
}

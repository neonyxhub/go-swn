package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	peers      *debugPeers
	swn1, swn2 debugPeerInfo
)

func fetchPeers() {
	if _, err := os.Stat("debug.yml"); os.IsNotExist(err) {
		log.Fatal("debug.yml is missing")
	}

	data, err := os.ReadFile("debug.yml")
	if err != nil {
		log.Fatalf("failed to read debug.yml: %v", err)
	}

	peers := &debugPeers{}
	err = yaml.Unmarshal(data, peers)
	if err != nil {
		log.Fatalf("failed to unmarshal debug.yml: %v", err)
	}

	if len(peers.Peers) != 2 {
		log.Fatal("should be 2 peers (swn1, swn2) in debug.yml")
	}

	for _, p := range peers.Peers {
		if p.GrpcServerPort == 8081 {
			swn1 = p
		} else if p.GrpcServerPort == 8082 {
			swn2 = p
		}
	}
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
	conn1, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to swn1 gRPC: %v", err)
	}
	defer conn1.Close()

	swnclient1 := pb.NewSWNBusClient(conn1)
	log.Println("connected to swn1 gRPC")

	// connect to swn2 gRPC
	conn2, err := grpc.Dial(":8082", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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
		Dest:    &pb.Destination{Addr: swn2MultiAddr.Bytes()},
		Lexicon: &pb.LexiconUri{Uri: "uri-123"},
		Data:    []byte("data-123"),
	}

	if err := stream.Send(event); err != nil {
		log.Fatalf("failed to send event: %v", err)
	}

	log.Println("closing consumer")
	<-done
}

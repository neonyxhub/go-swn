package swn

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const debugYml = "test/e2e/testdata/debug.yml"

// structure to be saved to debug.yml when debug mode is true
type debugPeerInfo struct {
	GrpcServerPort int    `yaml:"grpc_server_port"`
	PeerIpv4       string `yaml:"peer_ipv4"`
	PeerId         string `yaml:"peer_id"`
	TransportPort  string `yaml:"transport_port"`
}

type debugPeers struct {
	Peers []debugPeerInfo
}

func (s *SWN) DebugSavePeerInfo() error {
	port, err := s.Peer.GetTransportPort("tcp")
	if err != nil {
		return err
	}

	peer := debugPeerInfo{
		GrpcServerPort: s.GrpcServer.GetPort(),
		PeerIpv4:       s.Peer.GetIpv4(),
		PeerId:         s.ID().String(),
		TransportPort:  port,
	}

	Log.Info("saving debug info",
		zap.Int("gRPC_port", peer.GrpcServerPort),
		zap.String("gRPC IPv4", peer.PeerIpv4),
		zap.String("peerId", peer.PeerId),
		zap.String("peer tcp port", peer.TransportPort),
	)

	if _, err := os.Stat(debugYml); os.IsNotExist(err) {
		data, err := yaml.Marshal(&debugPeers{Peers: []debugPeerInfo{peer}})
		if err != nil {
			return err
		}
		err = os.WriteFile(debugYml, data, 0606)
		return err
	}

	// append to existing peers in debug.yml
	data, err := os.ReadFile(debugYml)
	if err != nil {
		return err
	}

	peers := &debugPeers{}
	err = yaml.Unmarshal(data, peers)
	if err != nil {
		return err
	}

	duplicate := false

	for _, p := range peers.Peers {
		if p.PeerId == peer.PeerId {
			duplicate = true
		}
	}

	if !duplicate {
		peers.Peers = append(peers.Peers, peer)

		data, err := yaml.Marshal(peers)
		if err != nil {
			return err
		}

		err = os.WriteFile(debugYml, data, 0606)
		return err
	}

	return nil
}

func (s *SWN) DebugDeletePeerInfo() error {
	peer := debugPeerInfo{
		PeerId: s.ID().String(),
	}

	data, err := os.ReadFile(debugYml)
	if err != nil {
		return err
	}

	peers := &debugPeers{}
	err = yaml.Unmarshal(data, peers)
	if err != nil {
		return err
	}

	for i, p := range peers.Peers {
		if p.PeerId == peer.PeerId {
			Log.Info("deleting debug info")

			peers.Peers = append(peers.Peers[:i], peers.Peers[i+1:]...)
			data, err := yaml.Marshal(peers)
			if err != nil {
				return err
			}

			err = os.WriteFile(debugYml, data, 0606)
			return err
		}
	}

	return nil
}

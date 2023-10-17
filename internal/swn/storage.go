package swn

import (
	auth_pb "go.neonyx.io/go-swn/internal/swn/pb"
	"google.golang.org/protobuf/proto"
)

const (
	dbRootKey = "/swn"
)

// Get device auth described in auth_model.proto
// LevelDB path: /swn/devPrvKey
func (s *SWN) GetDeviceAuth() (*auth_pb.DeviceAuth, error) {
	key := s.Ds.NewKey(dbRootKey, "devPrvKey")
	raw, err := s.Ds.Get(key, nil)
	if err != nil {
		return nil, err
	}

	deviceAuth := &auth_pb.DeviceAuth{}
	err = proto.Unmarshal(raw, deviceAuth)
	if err != nil {
		return nil, err
	}

	return deviceAuth, nil
}

// Save device auth described in auth_model.proto
// LevelDB path: /swn/devPrvKey
func (s *SWN) SaveDeviceAuth() error {
	key := s.Ds.NewKey(dbRootKey, string(s.Device.Id), "devPrvKey")

	deviceAuth := &auth_pb.DeviceAuth{
		PrivKey: s.Device.GetPrivKeyRaw(),
	}

	val, err := proto.Marshal(deviceAuth)
	if err != nil {
		return err
	}

	if err = s.Ds.Put(key, val, nil); err != nil {
		return err
	}

	return nil
}

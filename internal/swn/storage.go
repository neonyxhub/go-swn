package swn

import (
	"github.com/syndtr/goleveldb/leveldb"
	auth_pb "go.neonyx.io/go-swn/internal/swn/pb"
	"google.golang.org/protobuf/proto"
)

const (
	dbRootKey = "/swn"
)

// Get device auth described in auth_model.proto
// LevelDB path: /swn/prvkey
func (s *SWN) GetDeviceAuth() error {
	key := s.Ds.NewKey(dbRootKey, "prvkey")
	raw, err := s.Ds.Get(key, nil)
	if err != nil {
		return err
	}

	deviceAuth := &auth_pb.DeviceAuth{}
	err = proto.Unmarshal(raw, deviceAuth)
	if err != nil {
		return err
	}

	s.Device.ParsePrivKeyRaw(deviceAuth.PrivKey)

	return nil
}

// Save device auth described in auth_model.proto
// LevelDB path: /swn/prvkey
func (s *SWN) SaveDeviceAuth() error {
	key := s.Ds.NewKey(dbRootKey, "prvkey")

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

func (s *SWN) CheckDeviceId() error {
	if err := s.GetDeviceAuth(); err == leveldb.ErrNotFound {
		s.Log.Info("Generating a new device keypair")
		if err := s.Device.GenKeyPair(); err != nil {
			return err
		}

		if err = s.Device.GenDeviceId(); err != nil {
			return err
		}

		if err = s.SaveDeviceAuth(); err != nil {
			return err
		}

		return nil
	} else if err != nil {
		return err
	}

	s.Log.Info("Read an existing keypair")
	if err := s.Device.GenDeviceId(); err != nil {
		return err
	}

	return nil
}

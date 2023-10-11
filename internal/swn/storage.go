package swn

import (
	"crypto/sha256"

	"github.com/libp2p/go-libp2p/core/crypto"
	auth_pb "go.neonyx.io/go-swn/internal/swn/pb"
	"google.golang.org/protobuf/proto"
)

const (
	authDBkey = "/swnauth/"
)

// Get row for device id, described in auth/pb
func (s *SWN) GetAuthInfo(deviceId *auth_pb.DeviceId) (*auth_pb.AuthInfo, error) {
	key := s.Ds.NewKey(authDBkey, string(deviceId.DeviceId))
	raw, err := s.Ds.Get(key, nil)
	if err != nil {
		return nil, err
	}

	authInfo := auth_pb.AuthInfo{}

	err = proto.Unmarshal(raw, &authInfo)

	return &authInfo, err
}

// Get row for peer id, described in auth/pb
func (s *SWN) GetAuthInfoFromPeerId(peerId *auth_pb.PeerId) (*auth_pb.AuthInfo, error) {
	deviceId, err := s.GetDeviceForPeerId(peerId)

	if err != nil {
		return nil, err
	}
	return s.GetAuthInfo(deviceId)
}

// Get device id for peer id
func (s *SWN) GetDeviceForPeerId(peerId *auth_pb.PeerId) (*auth_pb.DeviceId, error) {
	return &auth_pb.DeviceId{}, nil
}

// Get device id for account id
func (s *SWN) GetDeviceForAcc(accountId *auth_pb.AccID) (*auth_pb.DeviceId, error) {
	return &auth_pb.DeviceId{}, nil
}

// Get row for account id, described in auth/pb
func (s *SWN) GetAuthInfoFromAccID(accountId *auth_pb.AccID) (*auth_pb.AuthInfo, error) {
	deviceId, err := s.GetDeviceForAcc(accountId)

	if err != nil {
		return nil, err
	}
	return s.GetAuthInfo(deviceId)
}

// Save row, described in auth/pb
func (s *SWN) SaveAuthInfo(devicePub crypto.PubKey, myPriv crypto.PrivKey) (*auth_pb.DeviceId, error) {
	hash := sha256.New()
	rawDeviceId, _ := devicePub.Raw()
	hash.Write(rawDeviceId)

	deviceId := hash.Sum(nil)[:12]
	key := s.Ds.NewKey(authDBkey, string(deviceId))

	otherDevicePubkey, _ := devicePub.Raw()
	myDevicePrivateKey, _ := myPriv.Raw()

	authInfo := &auth_pb.AuthInfo{
		OtherDevicePubKey:  otherDevicePubkey,
		MyDevicePrivateKey: myDevicePrivateKey,
	}

	data, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}

	err = s.Ds.Put(key, data, nil)

	return &auth_pb.DeviceId{DeviceId: deviceId}, err
}

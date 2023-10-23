// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: pkg/bus/pb/bus_api.proto

package pb

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EventType int32

const (
	EventType_REQ  EventType = 0
	EventType_RESP EventType = 1
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "REQ",
		1: "RESP",
	}
	EventType_value = map[string]int32{
		"REQ":  0,
		"RESP": 1,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_bus_pb_bus_api_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_pkg_bus_pb_bus_api_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{0}
}

type AccountAuthType int32

const (
	AccountAuthType_SIG AccountAuthType = 0
)

// Enum value maps for AccountAuthType.
var (
	AccountAuthType_name = map[int32]string{
		0: "SIG",
	}
	AccountAuthType_value = map[string]int32{
		"SIG": 0,
	}
)

func (x AccountAuthType) Enum() *AccountAuthType {
	p := new(AccountAuthType)
	*p = x
	return p
}

func (x AccountAuthType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccountAuthType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_bus_pb_bus_api_proto_enumTypes[1].Descriptor()
}

func (AccountAuthType) Type() protoreflect.EnumType {
	return &file_pkg_bus_pb_bus_api_proto_enumTypes[1]
}

func (x AccountAuthType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccountAuthType.Descriptor instead.
func (AccountAuthType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{1}
}

type Destination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr []byte `protobuf:"bytes,1,opt,name=Addr,proto3" json:"Addr,omitempty"`
}

func (x *Destination) Reset() {
	*x = Destination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Destination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Destination) ProtoMessage() {}

func (x *Destination) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Destination.ProtoReflect.Descriptor instead.
func (*Destination) Descriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{0}
}

func (x *Destination) GetAddr() []byte {
	if x != nil {
		return x.Addr
	}
	return nil
}

type LexiconUri struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uri string `protobuf:"bytes,1,opt,name=Uri,proto3" json:"Uri,omitempty"`
}

func (x *LexiconUri) Reset() {
	*x = LexiconUri{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LexiconUri) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LexiconUri) ProtoMessage() {}

func (x *LexiconUri) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LexiconUri.ProtoReflect.Descriptor instead.
func (*LexiconUri) Descriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{1}
}

func (x *LexiconUri) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

type AccountAuth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  AccountAuthType `protobuf:"varint,1,opt,name=Type,proto3,enum=bus_api.pb.AccountAuthType" json:"Type,omitempty"`
	Token []byte          `protobuf:"bytes,2,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *AccountAuth) Reset() {
	*x = AccountAuth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountAuth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountAuth) ProtoMessage() {}

func (x *AccountAuth) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountAuth.ProtoReflect.Descriptor instead.
func (*AccountAuth) Descriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{2}
}

func (x *AccountAuth) GetType() AccountAuthType {
	if x != nil {
		return x.Type
	}
	return AccountAuthType_SIG
}

func (x *AccountAuth) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64       `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Type        EventType    `protobuf:"varint,2,opt,name=Type,proto3,enum=bus_api.pb.EventType" json:"Type,omitempty"`
	Dest        *Destination `protobuf:"bytes,3,opt,name=Dest,proto3" json:"Dest,omitempty"`
	Lexicon     *LexiconUri  `protobuf:"bytes,4,opt,name=Lexicon,proto3" json:"Lexicon,omitempty"`
	Data        []byte       `protobuf:"bytes,5,opt,name=Data,proto3" json:"Data,omitempty"`
	AccountAuth *AccountAuth `protobuf:"bytes,6,opt,name=AccountAuth,proto3" json:"AccountAuth,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{3}
}

func (x *Event) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Event) GetType() EventType {
	if x != nil {
		return x.Type
	}
	return EventType_REQ
}

func (x *Event) GetDest() *Destination {
	if x != nil {
		return x.Dest
	}
	return nil
}

func (x *Event) GetLexicon() *LexiconUri {
	if x != nil {
		return x.Lexicon
	}
	return nil
}

func (x *Event) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Event) GetAccountAuth() *AccountAuth {
	if x != nil {
		return x.AccountAuth
	}
	return nil
}

type Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []byte `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *Peer) Reset() {
	*x = Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{4}
}

func (x *Peer) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

type ListenEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListenEventsRequest) Reset() {
	*x = ListenEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListenEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenEventsRequest) ProtoMessage() {}

func (x *ListenEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenEventsRequest.ProtoReflect.Descriptor instead.
func (*ListenEventsRequest) Descriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{5}
}

type StreamEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StreamEventsResponse) Reset() {
	*x = StreamEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamEventsResponse) ProtoMessage() {}

func (x *StreamEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_bus_pb_bus_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamEventsResponse.ProtoReflect.Descriptor instead.
func (*StreamEventsResponse) Descriptor() ([]byte, []int) {
	return file_pkg_bus_pb_bus_api_proto_rawDescGZIP(), []int{6}
}

var File_pkg_bus_pb_bus_api_proto protoreflect.FileDescriptor

var file_pkg_bus_pb_bus_api_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x6b, 0x67, 0x2f, 0x62, 0x75, 0x73, 0x2f, 0x70, 0x62, 0x2f, 0x62, 0x75, 0x73,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x62, 0x75, 0x73, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x21, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x41, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x41, 0x64, 0x64, 0x72, 0x22, 0x1e, 0x0a, 0x0a, 0x4c, 0x65, 0x78, 0x69, 0x63, 0x6f,
	0x6e, 0x55, 0x72, 0x69, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x72, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x55, 0x72, 0x69, 0x22, 0x54, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x41, 0x75, 0x74, 0x68, 0x12, 0x2f, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62,
	0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x75, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xf0, 0x01, 0x0a,
	0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x62, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x2b, 0x0a, 0x04, 0x44, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x44, 0x65, 0x73, 0x74, 0x12, 0x30,
	0x0a, 0x07, 0x4c, 0x65, 0x78, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x65, 0x78,
	0x69, 0x63, 0x6f, 0x6e, 0x55, 0x72, 0x69, 0x52, 0x07, 0x4c, 0x65, 0x78, 0x69, 0x63, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x39, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41,
	0x75, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x62, 0x75, 0x73, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x75,
	0x74, 0x68, 0x52, 0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x75, 0x74, 0x68, 0x22,
	0x16, 0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x02, 0x49, 0x64, 0x22, 0x15, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x16,
	0x0a, 0x14, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x1e, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x52, 0x45, 0x51, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x52, 0x45, 0x53, 0x50, 0x10, 0x01, 0x2a, 0x1a, 0x0a, 0x0f, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x41, 0x75, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x49, 0x47,
	0x10, 0x00, 0x32, 0xe2, 0x01, 0x0a, 0x06, 0x53, 0x57, 0x4e, 0x42, 0x75, 0x73, 0x12, 0x50, 0x0a,
	0x15, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x11, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69,
	0x2e, 0x70, 0x62, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x20, 0x2e, 0x62, 0x75, 0x73, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12,
	0x4b, 0x0a, 0x11, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x46, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x1f, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x62, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x30, 0x01, 0x12, 0x39, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x62, 0x75, 0x73, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62,
	0x2e, 0x50, 0x65, 0x65, 0x72, 0x22, 0x00, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x6f, 0x2e, 0x6e, 0x65,
	0x6f, 0x6e, 0x79, 0x78, 0x2e, 0x69, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x77, 0x6e, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x62, 0x75, 0x73, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_pkg_bus_pb_bus_api_proto_rawDescOnce sync.Once
	file_pkg_bus_pb_bus_api_proto_rawDescData = file_pkg_bus_pb_bus_api_proto_rawDesc
)

func file_pkg_bus_pb_bus_api_proto_rawDescGZIP() []byte {
	file_pkg_bus_pb_bus_api_proto_rawDescOnce.Do(func() {
		file_pkg_bus_pb_bus_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_bus_pb_bus_api_proto_rawDescData)
	})
	return file_pkg_bus_pb_bus_api_proto_rawDescData
}

var file_pkg_bus_pb_bus_api_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_pkg_bus_pb_bus_api_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pkg_bus_pb_bus_api_proto_goTypes = []interface{}{
	(EventType)(0),               // 0: bus_api.pb.EventType
	(AccountAuthType)(0),         // 1: bus_api.pb.AccountAuthType
	(*Destination)(nil),          // 2: bus_api.pb.Destination
	(*LexiconUri)(nil),           // 3: bus_api.pb.LexiconUri
	(*AccountAuth)(nil),          // 4: bus_api.pb.AccountAuth
	(*Event)(nil),                // 5: bus_api.pb.Event
	(*Peer)(nil),                 // 6: bus_api.pb.Peer
	(*ListenEventsRequest)(nil),  // 7: bus_api.pb.ListenEventsRequest
	(*StreamEventsResponse)(nil), // 8: bus_api.pb.StreamEventsResponse
	(*empty.Empty)(nil),          // 9: google.protobuf.Empty
}
var file_pkg_bus_pb_bus_api_proto_depIdxs = []int32{
	1, // 0: bus_api.pb.AccountAuth.Type:type_name -> bus_api.pb.AccountAuthType
	0, // 1: bus_api.pb.Event.Type:type_name -> bus_api.pb.EventType
	2, // 2: bus_api.pb.Event.Dest:type_name -> bus_api.pb.Destination
	3, // 3: bus_api.pb.Event.Lexicon:type_name -> bus_api.pb.LexiconUri
	4, // 4: bus_api.pb.Event.AccountAuth:type_name -> bus_api.pb.AccountAuth
	5, // 5: bus_api.pb.SWNBus.LocalDistributeEvents:input_type -> bus_api.pb.Event
	7, // 6: bus_api.pb.SWNBus.LocalFunnelEvents:input_type -> bus_api.pb.ListenEventsRequest
	9, // 7: bus_api.pb.SWNBus.GetPeerInfo:input_type -> google.protobuf.Empty
	8, // 8: bus_api.pb.SWNBus.LocalDistributeEvents:output_type -> bus_api.pb.StreamEventsResponse
	5, // 9: bus_api.pb.SWNBus.LocalFunnelEvents:output_type -> bus_api.pb.Event
	6, // 10: bus_api.pb.SWNBus.GetPeerInfo:output_type -> bus_api.pb.Peer
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pkg_bus_pb_bus_api_proto_init() }
func file_pkg_bus_pb_bus_api_proto_init() {
	if File_pkg_bus_pb_bus_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_bus_pb_bus_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Destination); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_bus_pb_bus_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LexiconUri); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_bus_pb_bus_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountAuth); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_bus_pb_bus_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_bus_pb_bus_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Peer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_bus_pb_bus_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListenEventsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_bus_pb_bus_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamEventsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_bus_pb_bus_api_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_bus_pb_bus_api_proto_goTypes,
		DependencyIndexes: file_pkg_bus_pb_bus_api_proto_depIdxs,
		EnumInfos:         file_pkg_bus_pb_bus_api_proto_enumTypes,
		MessageInfos:      file_pkg_bus_pb_bus_api_proto_msgTypes,
	}.Build()
	File_pkg_bus_pb_bus_api_proto = out.File
	file_pkg_bus_pb_bus_api_proto_rawDesc = nil
	file_pkg_bus_pb_bus_api_proto_goTypes = nil
	file_pkg_bus_pb_bus_api_proto_depIdxs = nil
}

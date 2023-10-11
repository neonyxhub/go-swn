// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: pkg/bus/pb/bus_api.proto

package pb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SWNBus_LocalDistributeEvents_FullMethodName = "/bus_api.pb.SWNBus/LocalDistributeEvents"
	SWNBus_LocalFunnelEvents_FullMethodName     = "/bus_api.pb.SWNBus/LocalFunnelEvents"
	SWNBus_GetPeerId_FullMethodName             = "/bus_api.pb.SWNBus/GetPeerId"
)

// SWNBusClient is the client API for SWNBus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SWNBusClient interface {
	// swn here is the server, which listens to events from cwn
	LocalDistributeEvents(ctx context.Context, opts ...grpc.CallOption) (SWNBus_LocalDistributeEventsClient, error)
	// swn here is the server, which gives events to CWN
	LocalFunnelEvents(ctx context.Context, in *ListenEventsRequest, opts ...grpc.CallOption) (SWNBus_LocalFunnelEventsClient, error)
	GetPeerId(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Peer, error)
}

type sWNBusClient struct {
	cc grpc.ClientConnInterface
}

func NewSWNBusClient(cc grpc.ClientConnInterface) SWNBusClient {
	return &sWNBusClient{cc}
}

func (c *sWNBusClient) LocalDistributeEvents(ctx context.Context, opts ...grpc.CallOption) (SWNBus_LocalDistributeEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &SWNBus_ServiceDesc.Streams[0], SWNBus_LocalDistributeEvents_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &sWNBusLocalDistributeEventsClient{stream}
	return x, nil
}

type SWNBus_LocalDistributeEventsClient interface {
	Send(*Event) error
	CloseAndRecv() (*StreamEventsResponse, error)
	grpc.ClientStream
}

type sWNBusLocalDistributeEventsClient struct {
	grpc.ClientStream
}

func (x *sWNBusLocalDistributeEventsClient) Send(m *Event) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sWNBusLocalDistributeEventsClient) CloseAndRecv() (*StreamEventsResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StreamEventsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sWNBusClient) LocalFunnelEvents(ctx context.Context, in *ListenEventsRequest, opts ...grpc.CallOption) (SWNBus_LocalFunnelEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &SWNBus_ServiceDesc.Streams[1], SWNBus_LocalFunnelEvents_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &sWNBusLocalFunnelEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SWNBus_LocalFunnelEventsClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type sWNBusLocalFunnelEventsClient struct {
	grpc.ClientStream
}

func (x *sWNBusLocalFunnelEventsClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sWNBusClient) GetPeerId(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Peer, error) {
	out := new(Peer)
	err := c.cc.Invoke(ctx, SWNBus_GetPeerId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SWNBusServer is the server API for SWNBus service.
// All implementations must embed UnimplementedSWNBusServer
// for forward compatibility
type SWNBusServer interface {
	// swn here is the server, which listens to events from cwn
	LocalDistributeEvents(SWNBus_LocalDistributeEventsServer) error
	// swn here is the server, which gives events to CWN
	LocalFunnelEvents(*ListenEventsRequest, SWNBus_LocalFunnelEventsServer) error
	GetPeerId(context.Context, *empty.Empty) (*Peer, error)
	mustEmbedUnimplementedSWNBusServer()
}

// UnimplementedSWNBusServer must be embedded to have forward compatible implementations.
type UnimplementedSWNBusServer struct {
}

func (UnimplementedSWNBusServer) LocalDistributeEvents(SWNBus_LocalDistributeEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method LocalDistributeEvents not implemented")
}
func (UnimplementedSWNBusServer) LocalFunnelEvents(*ListenEventsRequest, SWNBus_LocalFunnelEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method LocalFunnelEvents not implemented")
}
func (UnimplementedSWNBusServer) GetPeerId(context.Context, *empty.Empty) (*Peer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeerId not implemented")
}
func (UnimplementedSWNBusServer) mustEmbedUnimplementedSWNBusServer() {}

// UnsafeSWNBusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SWNBusServer will
// result in compilation errors.
type UnsafeSWNBusServer interface {
	mustEmbedUnimplementedSWNBusServer()
}

func RegisterSWNBusServer(s grpc.ServiceRegistrar, srv SWNBusServer) {
	s.RegisterService(&SWNBus_ServiceDesc, srv)
}

func _SWNBus_LocalDistributeEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SWNBusServer).LocalDistributeEvents(&sWNBusLocalDistributeEventsServer{stream})
}

type SWNBus_LocalDistributeEventsServer interface {
	SendAndClose(*StreamEventsResponse) error
	Recv() (*Event, error)
	grpc.ServerStream
}

type sWNBusLocalDistributeEventsServer struct {
	grpc.ServerStream
}

func (x *sWNBusLocalDistributeEventsServer) SendAndClose(m *StreamEventsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sWNBusLocalDistributeEventsServer) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SWNBus_LocalFunnelEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenEventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SWNBusServer).LocalFunnelEvents(m, &sWNBusLocalFunnelEventsServer{stream})
}

type SWNBus_LocalFunnelEventsServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type sWNBusLocalFunnelEventsServer struct {
	grpc.ServerStream
}

func (x *sWNBusLocalFunnelEventsServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func _SWNBus_GetPeerId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SWNBusServer).GetPeerId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SWNBus_GetPeerId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SWNBusServer).GetPeerId(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SWNBus_ServiceDesc is the grpc.ServiceDesc for SWNBus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SWNBus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bus_api.pb.SWNBus",
	HandlerType: (*SWNBusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPeerId",
			Handler:    _SWNBus_GetPeerId_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "LocalDistributeEvents",
			Handler:       _SWNBus_LocalDistributeEvents_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "LocalFunnelEvents",
			Handler:       _SWNBus_LocalFunnelEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/bus/pb/bus_api.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package jointrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// JointRPCClient is the client API for JointRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JointRPCClient interface {
	Call(ctx context.Context, in *JSONRPCCallRequest, opts ...grpc.CallOption) (*JSONRPCCallResult, error)
	Notify(ctx context.Context, in *JSONRPCNotifyRequest, opts ...grpc.CallOption) (*JSONRPCNotifyResponse, error)
	ListMethods(ctx context.Context, in *ListMethodsRequest, opts ...grpc.CallOption) (*ListMethodsResponse, error)
	ListDelegates(ctx context.Context, in *ListDelegatesRequest, opts ...grpc.CallOption) (*ListDelegatesResponse, error)
	// request/response dual streams
	Live(ctx context.Context, opts ...grpc.CallOption) (JointRPC_LiveClient, error)
}

type jointRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewJointRPCClient(cc grpc.ClientConnInterface) JointRPCClient {
	return &jointRPCClient{cc}
}

func (c *jointRPCClient) Call(ctx context.Context, in *JSONRPCCallRequest, opts ...grpc.CallOption) (*JSONRPCCallResult, error) {
	out := new(JSONRPCCallResult)
	err := c.cc.Invoke(ctx, "/JointRPC/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jointRPCClient) Notify(ctx context.Context, in *JSONRPCNotifyRequest, opts ...grpc.CallOption) (*JSONRPCNotifyResponse, error) {
	out := new(JSONRPCNotifyResponse)
	err := c.cc.Invoke(ctx, "/JointRPC/Notify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jointRPCClient) ListMethods(ctx context.Context, in *ListMethodsRequest, opts ...grpc.CallOption) (*ListMethodsResponse, error) {
	out := new(ListMethodsResponse)
	err := c.cc.Invoke(ctx, "/JointRPC/ListMethods", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jointRPCClient) ListDelegates(ctx context.Context, in *ListDelegatesRequest, opts ...grpc.CallOption) (*ListDelegatesResponse, error) {
	out := new(ListDelegatesResponse)
	err := c.cc.Invoke(ctx, "/JointRPC/ListDelegates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jointRPCClient) Live(ctx context.Context, opts ...grpc.CallOption) (JointRPC_LiveClient, error) {
	stream, err := c.cc.NewStream(ctx, &_JointRPC_serviceDesc.Streams[0], "/JointRPC/Live", opts...)
	if err != nil {
		return nil, err
	}
	x := &jointRPCLiveClient{stream}
	return x, nil
}

type JointRPC_LiveClient interface {
	Send(*JSONRPCEnvolope) error
	Recv() (*JSONRPCEnvolope, error)
	grpc.ClientStream
}

type jointRPCLiveClient struct {
	grpc.ClientStream
}

func (x *jointRPCLiveClient) Send(m *JSONRPCEnvolope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *jointRPCLiveClient) Recv() (*JSONRPCEnvolope, error) {
	m := new(JSONRPCEnvolope)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JointRPCServer is the server API for JointRPC service.
// All implementations must embed UnimplementedJointRPCServer
// for forward compatibility
type JointRPCServer interface {
	Call(context.Context, *JSONRPCCallRequest) (*JSONRPCCallResult, error)
	Notify(context.Context, *JSONRPCNotifyRequest) (*JSONRPCNotifyResponse, error)
	ListMethods(context.Context, *ListMethodsRequest) (*ListMethodsResponse, error)
	ListDelegates(context.Context, *ListDelegatesRequest) (*ListDelegatesResponse, error)
	// request/response dual streams
	Live(JointRPC_LiveServer) error
	mustEmbedUnimplementedJointRPCServer()
}

// UnimplementedJointRPCServer must be embedded to have forward compatible implementations.
type UnimplementedJointRPCServer struct {
}

func (UnimplementedJointRPCServer) Call(context.Context, *JSONRPCCallRequest) (*JSONRPCCallResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (UnimplementedJointRPCServer) Notify(context.Context, *JSONRPCNotifyRequest) (*JSONRPCNotifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
}
func (UnimplementedJointRPCServer) ListMethods(context.Context, *ListMethodsRequest) (*ListMethodsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMethods not implemented")
}
func (UnimplementedJointRPCServer) ListDelegates(context.Context, *ListDelegatesRequest) (*ListDelegatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDelegates not implemented")
}
func (UnimplementedJointRPCServer) Live(JointRPC_LiveServer) error {
	return status.Errorf(codes.Unimplemented, "method Live not implemented")
}
func (UnimplementedJointRPCServer) mustEmbedUnimplementedJointRPCServer() {}

// UnsafeJointRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JointRPCServer will
// result in compilation errors.
type UnsafeJointRPCServer interface {
	mustEmbedUnimplementedJointRPCServer()
}

func RegisterJointRPCServer(s *grpc.Server, srv JointRPCServer) {
	s.RegisterService(&_JointRPC_serviceDesc, srv)
}

func _JointRPC_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JSONRPCCallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JointRPCServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JointRPC/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JointRPCServer).Call(ctx, req.(*JSONRPCCallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JointRPC_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JSONRPCNotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JointRPCServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JointRPC/Notify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JointRPCServer).Notify(ctx, req.(*JSONRPCNotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JointRPC_ListMethods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMethodsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JointRPCServer).ListMethods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JointRPC/ListMethods",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JointRPCServer).ListMethods(ctx, req.(*ListMethodsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JointRPC_ListDelegates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDelegatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JointRPCServer).ListDelegates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JointRPC/ListDelegates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JointRPCServer).ListDelegates(ctx, req.(*ListDelegatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JointRPC_Live_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(JointRPCServer).Live(&jointRPCLiveServer{stream})
}

type JointRPC_LiveServer interface {
	Send(*JSONRPCEnvolope) error
	Recv() (*JSONRPCEnvolope, error)
	grpc.ServerStream
}

type jointRPCLiveServer struct {
	grpc.ServerStream
}

func (x *jointRPCLiveServer) Send(m *JSONRPCEnvolope) error {
	return x.ServerStream.SendMsg(m)
}

func (x *jointRPCLiveServer) Recv() (*JSONRPCEnvolope, error) {
	m := new(JSONRPCEnvolope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _JointRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "JointRPC",
	HandlerType: (*JointRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _JointRPC_Call_Handler,
		},
		{
			MethodName: "Notify",
			Handler:    _JointRPC_Notify_Handler,
		},
		{
			MethodName: "ListMethods",
			Handler:    _JointRPC_ListMethods_Handler,
		},
		{
			MethodName: "ListDelegates",
			Handler:    _JointRPC_ListDelegates_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Live",
			Handler:       _JointRPC_Live_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "jointrpc.proto",
}

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
	// state stream
	SubscribeState(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (JointRPC_SubscribeStateClient, error)
	// request/response dual streams
	Worker(ctx context.Context, opts ...grpc.CallOption) (JointRPC_WorkerClient, error)
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

func (c *jointRPCClient) SubscribeState(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (JointRPC_SubscribeStateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_JointRPC_serviceDesc.Streams[0], "/JointRPC/SubscribeState", opts...)
	if err != nil {
		return nil, err
	}
	x := &jointRPCSubscribeStateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type JointRPC_SubscribeStateClient interface {
	Recv() (*SubscribeStateResponse, error)
	grpc.ClientStream
}

type jointRPCSubscribeStateClient struct {
	grpc.ClientStream
}

func (x *jointRPCSubscribeStateClient) Recv() (*SubscribeStateResponse, error) {
	m := new(SubscribeStateResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *jointRPCClient) Worker(ctx context.Context, opts ...grpc.CallOption) (JointRPC_WorkerClient, error) {
	stream, err := c.cc.NewStream(ctx, &_JointRPC_serviceDesc.Streams[1], "/JointRPC/Worker", opts...)
	if err != nil {
		return nil, err
	}
	x := &jointRPCWorkerClient{stream}
	return x, nil
}

type JointRPC_WorkerClient interface {
	Send(*JSONRPCEnvolope) error
	Recv() (*JSONRPCEnvolope, error)
	grpc.ClientStream
}

type jointRPCWorkerClient struct {
	grpc.ClientStream
}

func (x *jointRPCWorkerClient) Send(m *JSONRPCEnvolope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *jointRPCWorkerClient) Recv() (*JSONRPCEnvolope, error) {
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
	// state stream
	SubscribeState(*AuthRequest, JointRPC_SubscribeStateServer) error
	// request/response dual streams
	Worker(JointRPC_WorkerServer) error
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
func (UnimplementedJointRPCServer) SubscribeState(*AuthRequest, JointRPC_SubscribeStateServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeState not implemented")
}
func (UnimplementedJointRPCServer) Worker(JointRPC_WorkerServer) error {
	return status.Errorf(codes.Unimplemented, "method Worker not implemented")
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

func _JointRPC_SubscribeState_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AuthRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JointRPCServer).SubscribeState(m, &jointRPCSubscribeStateServer{stream})
}

type JointRPC_SubscribeStateServer interface {
	Send(*SubscribeStateResponse) error
	grpc.ServerStream
}

type jointRPCSubscribeStateServer struct {
	grpc.ServerStream
}

func (x *jointRPCSubscribeStateServer) Send(m *SubscribeStateResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _JointRPC_Worker_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(JointRPCServer).Worker(&jointRPCWorkerServer{stream})
}

type JointRPC_WorkerServer interface {
	Send(*JSONRPCEnvolope) error
	Recv() (*JSONRPCEnvolope, error)
	grpc.ServerStream
}

type jointRPCWorkerServer struct {
	grpc.ServerStream
}

func (x *jointRPCWorkerServer) Send(m *JSONRPCEnvolope) error {
	return x.ServerStream.SendMsg(m)
}

func (x *jointRPCWorkerServer) Recv() (*JSONRPCEnvolope, error) {
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
			StreamName:    "SubscribeState",
			Handler:       _JointRPC_SubscribeState_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Worker",
			Handler:       _JointRPC_Worker_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "jointrpc.proto",
}

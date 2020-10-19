// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package tube

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// JSONRPCTubeClient is the client API for JSONRPCTube service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JSONRPCTubeClient interface {
	Call(ctx context.Context, in *JSONRPCRequest, opts ...grpc.CallOption) (*JSONRPCResult, error)
	Handle(ctx context.Context, opts ...grpc.CallOption) (JSONRPCTube_HandleClient, error)
}

type jSONRPCTubeClient struct {
	cc grpc.ClientConnInterface
}

func NewJSONRPCTubeClient(cc grpc.ClientConnInterface) JSONRPCTubeClient {
	return &jSONRPCTubeClient{cc}
}

func (c *jSONRPCTubeClient) Call(ctx context.Context, in *JSONRPCRequest, opts ...grpc.CallOption) (*JSONRPCResult, error) {
	out := new(JSONRPCResult)
	err := c.cc.Invoke(ctx, "/JSONRPCTube/call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jSONRPCTubeClient) Handle(ctx context.Context, opts ...grpc.CallOption) (JSONRPCTube_HandleClient, error) {
	stream, err := c.cc.NewStream(ctx, &_JSONRPCTube_serviceDesc.Streams[0], "/JSONRPCTube/handle", opts...)
	if err != nil {
		return nil, err
	}
	x := &jSONRPCTubeHandleClient{stream}
	return x, nil
}

type JSONRPCTube_HandleClient interface {
	Send(*JSONRPCResult) error
	Recv() (*JSONRPCRequest, error)
	grpc.ClientStream
}

type jSONRPCTubeHandleClient struct {
	grpc.ClientStream
}

func (x *jSONRPCTubeHandleClient) Send(m *JSONRPCResult) error {
	return x.ClientStream.SendMsg(m)
}

func (x *jSONRPCTubeHandleClient) Recv() (*JSONRPCRequest, error) {
	m := new(JSONRPCRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JSONRPCTubeServer is the server API for JSONRPCTube service.
// All implementations must embed UnimplementedJSONRPCTubeServer
// for forward compatibility
type JSONRPCTubeServer interface {
	Call(context.Context, *JSONRPCRequest) (*JSONRPCResult, error)
	Handle(JSONRPCTube_HandleServer) error
	mustEmbedUnimplementedJSONRPCTubeServer()
}

// UnimplementedJSONRPCTubeServer must be embedded to have forward compatible implementations.
type UnimplementedJSONRPCTubeServer struct {
}

func (UnimplementedJSONRPCTubeServer) Call(context.Context, *JSONRPCRequest) (*JSONRPCResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (UnimplementedJSONRPCTubeServer) Handle(JSONRPCTube_HandleServer) error {
	return status.Errorf(codes.Unimplemented, "method Handle not implemented")
}
func (UnimplementedJSONRPCTubeServer) mustEmbedUnimplementedJSONRPCTubeServer() {}

// UnsafeJSONRPCTubeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JSONRPCTubeServer will
// result in compilation errors.
type UnsafeJSONRPCTubeServer interface {
	mustEmbedUnimplementedJSONRPCTubeServer()
}

func RegisterJSONRPCTubeServer(s *grpc.Server, srv JSONRPCTubeServer) {
	s.RegisterService(&_JSONRPCTube_serviceDesc, srv)
}

func _JSONRPCTube_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JSONRPCRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JSONRPCTubeServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JSONRPCTube/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JSONRPCTubeServer).Call(ctx, req.(*JSONRPCRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JSONRPCTube_Handle_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(JSONRPCTubeServer).Handle(&jSONRPCTubeHandleServer{stream})
}

type JSONRPCTube_HandleServer interface {
	Send(*JSONRPCRequest) error
	Recv() (*JSONRPCResult, error)
	grpc.ServerStream
}

type jSONRPCTubeHandleServer struct {
	grpc.ServerStream
}

func (x *jSONRPCTubeHandleServer) Send(m *JSONRPCRequest) error {
	return x.ServerStream.SendMsg(m)
}

func (x *jSONRPCTubeHandleServer) Recv() (*JSONRPCResult, error) {
	m := new(JSONRPCResult)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _JSONRPCTube_serviceDesc = grpc.ServiceDesc{
	ServiceName: "JSONRPCTube",
	HandlerType: (*JSONRPCTubeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "call",
			Handler:    _JSONRPCTube_Call_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "handle",
			Handler:       _JSONRPCTube_Handle_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "tube.proto",
}

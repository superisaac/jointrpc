package server

import (
	"context"
	"math"
	//"flag"
	log "github.com/sirupsen/logrus"
	"net"
	//"os"
	//"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	datadir "github.com/superisaac/jointrpc/datadir"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	peer "google.golang.org/grpc/peer"
)

func ServerContext(rootCtx context.Context, factory *rpcrouter.RouterFactory) context.Context {
	if factory == nil {
		factory = rpcrouter.NewRouterFactory("server11")
	}
	aCtx := misc.NewBinder(rootCtx).Bind("routerfactory", factory).Context()
	return aCtx
}

func StartGRPCServer(rootCtx context.Context, bind string, opts ...grpc.ServerOption) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	} else {
		log.Debugf("entry server listen at %s", bind)
	}

	if r := rootCtx.Value("routerfactory"); r == nil {
		// no router attached, spawn a context with default router and cfg
		rootCtx = ServerContext(rootCtx, nil)
	}

	factory := rpcrouter.RouterFactoryFromContext(rootCtx)
	go factory.Start(rootCtx)

	cfg := factory.Config

	opts = append(opts,
		grpc.UnaryInterceptor(
			unaryBindContext(factory, cfg)),
		grpc.StreamInterceptor(
			streamBindContext(factory, cfg)),
		grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.WriteBufferSize(1024000),
		grpc.ReadBufferSize(1024000),
		grpc.ReadBufferSize(1024000),
	)
	grpcServer := grpc.NewServer(opts...)

	serverCtx, cancelServer := context.WithCancel(rootCtx)
	defer cancelServer()

	go func() {
		for {
			<-serverCtx.Done()
			log.Debugf("gRPC Server %s stops", bind)
			grpcServer.Stop()
			return
		}
	}()

	intf.RegisterJointRPCServer(grpcServer, NewJointRPCServer())
	grpcServer.Serve(lis)
}

func unaryBindContext(factory *rpcrouter.RouterFactory, cfg *datadir.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		b := misc.NewBinder(ctx)
		b = b.Bind("routerfactory", factory).Bind("config", cfg)
		if remotePeer, ok := peer.FromContext(ctx); ok {
			b = b.Bind("remoteAddress", remotePeer.Addr)
		} else {
			log.Warnf("fail to get the remote address")
		}

		cCtx := b.Context()
		h, err := handler(cCtx, req)
		return h, err
	}
}

func streamBindContext(factory *rpcrouter.RouterFactory, cfg *datadir.Config) grpc.StreamServerInterceptor {
	return func(srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {

		rootCtx := ss.Context()
		b := misc.NewBinder(ss.Context())
		b = b.Bind("routerfactory", factory).Bind("config", cfg)
		if remotePeer, ok := peer.FromContext(rootCtx); ok {
			b = b.Bind("remoteAddress", remotePeer.Addr)
		} else {
			log.Warnf("fail to get the remote address")
		}
		wrappedStream := grpc_middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = b.Context()
		return handler(srv, wrappedStream)
	}
}

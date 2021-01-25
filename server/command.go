package server

import (
	"context"
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
)

func ServerContext(rootCtx context.Context, router *rpcrouter.Router) context.Context {
	if router == nil {
		router = rpcrouter.NewRouter("server")
	}
	aCtx := misc.NewBinder(rootCtx).Bind("router", router).Context()
	return aCtx
}

func StartGRPCServer(rootCtx context.Context, bind string, opts ...grpc.ServerOption) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	} else {
		log.Debugf("entry server listen at %s", bind)
	}

	if r := rootCtx.Value("router"); r == nil {
		// no router attached, spawn a context with default router and cfg
		rootCtx = ServerContext(rootCtx, nil)
	}

	router := rpcrouter.RouterFromContext(rootCtx)
	go router.Start(rootCtx)

	cfg := router.Config

	//go handler.StartBuiltinHandlerManager(rootCtx)

	opts = append(opts,
		grpc.UnaryInterceptor(
			unaryBindContext(router, cfg)),
		grpc.StreamInterceptor(
			streamBindContext(router, cfg)))
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

func unaryBindContext(router *rpcrouter.Router, cfg *datadir.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		b := misc.NewBinder(ctx)
		cCtx := b.Bind("router", router).Bind("config", cfg).Context()
		h, err := handler(cCtx, req)
		return h, err
	}
}

func streamBindContext(router *rpcrouter.Router, cfg *datadir.Config) grpc.StreamServerInterceptor {
	return func(srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {

		b := misc.NewBinder(ss.Context())
		wrappedStream := grpc_middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = b.Bind("router", router).Bind("config", cfg).Context()
		return handler(srv, wrappedStream)
	}
}

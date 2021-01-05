package server

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	//"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	datadir "github.com/superisaac/jointrpc/datadir"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	mirror "github.com/superisaac/jointrpc/mirror"
	"github.com/superisaac/jointrpc/rpcrouter"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
)

func CommandStartServer() {
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	pBind := serverFlags.String("b", "", "The grpc server address and port")
	pDatadir := serverFlags.String("d", "", "The datadir to store configs")
	pCertFile := serverFlags.String("cert", "", "tls cert file")
	pKeyFile := serverFlags.String("key", "", "tls key file")
	//httpBind := serverFlags.String("httpd", "127.0.0.1:50056", "http address and port")

	serverFlags.Parse(os.Args[2:])
	if *pDatadir != "" {
		datadir.SetDatadir(*pDatadir)
	}

	cfg := datadir.NewConfig()
	cfg.ParseDatadir()

	//go StartHTTPd(*httpBind)
	var opts []grpc.ServerOption
	// server bind
	bind := *pBind
	if bind == "" {
		bind = cfg.Server.Bind
	}

	// tls settings
	certFile := *pCertFile
	if certFile == "" {
		certFile = cfg.Server.TLS.CertFile
	}

	keyFile := *pKeyFile
	if keyFile == "" {
		keyFile = cfg.Server.TLS.KeyFile
	}

	if certFile != "" || keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	StartServer(context.Background(), bind, cfg, opts...)
}

func StartServer(rootCtx context.Context, bind string, cfg *datadir.Config, opts ...grpc.ServerOption) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	} else {
		log.Debugf("entry server listen at %s", bind)
	}

	if cfg == nil {
		cfg = datadir.NewConfig()
	}

	router := rpcrouter.NewRouter("grpc_server")
	go router.Start(rootCtx)

	aCtx := context.WithValue(rootCtx, "config", cfg)
	aCtx = context.WithValue(aCtx, "router", router)

	handler.StartBuiltinHandlerManager(aCtx)

	mirror.StartMirrorsForPeers(aCtx)

	opts = append(opts,
		grpc.UnaryInterceptor(
			unaryBindContext(router, cfg)),
		grpc.StreamInterceptor(
			streamBindContext(router, cfg)))

	grpcServer := grpc.NewServer(opts...)
	intf.RegisterJointRPCServer(grpcServer, NewJointRPCServer())
	grpcServer.Serve(lis)
}

func unaryBindContext(router *rpcrouter.Router, cfg *datadir.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		rCtx := context.WithValue(ctx, "router", router)
		cCtx := context.WithValue(rCtx, "config", cfg)
		h, err := handler(cCtx, req)
		return h, err
	}
}

func streamBindContext(router *rpcrouter.Router, cfg *datadir.Config) grpc.StreamServerInterceptor {
	return func(srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		rCtx := context.WithValue(ss.Context(), "router", router)
		cCtx := context.WithValue(rCtx, "config", cfg)
		wrappedStream := grpc_middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = cCtx
		return handler(srv, wrappedStream)
	}
}

package server

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	//"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	mirror "github.com/superisaac/rpctube/cluster/mirror"
	datadir "github.com/superisaac/rpctube/datadir"
	intf "github.com/superisaac/rpctube/intf/tube"
	tube "github.com/superisaac/rpctube/tube"
	handler "github.com/superisaac/rpctube/tube/handler"
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

	//go StartHTTPd(*httpBind)
	cfg := datadir.GetConfig()
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
	StartServer(context.Background(), bind, opts...)
}

func StartServer(rootCtx context.Context, bind string, opts ...grpc.ServerOption) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	} else {
		log.Debugf("entry server listen at %s", bind)
	}

	//tube.Tube().Start(rootCtx)

	router := tube.NewRouter("grpc_server")
	go router.Start(rootCtx)
	ctxWithRouter := context.WithValue(rootCtx, "router", router)

	handler.StartBuiltinHandlerManager(ctxWithRouter)

	mirror.StartMirrorsForPeers(ctxWithRouter)

	opts = append(opts, grpc.UnaryInterceptor(
		unaryRouterAssigner(router)))

	opts = append(opts, grpc.StreamInterceptor(
		streamRouterAssigner(router)))

	grpcServer := grpc.NewServer(opts...)
	intf.RegisterJSONRPCTubeServer(grpcServer, NewJSONRPCTubeServer())
	grpcServer.Serve(lis)
}

func unaryRouterAssigner(router *tube.Router) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		rCtx := context.WithValue(ctx, "router", router)
		h, err := handler(rCtx, req)
		return h, err
	}
}

func streamRouterAssigner(router *tube.Router) grpc.StreamServerInterceptor {
	return func(srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		rCtx := context.WithValue(ss.Context(), "router", router)
		wrappedStream := grpc_middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = rCtx
		return handler(srv, wrappedStream)
	}
}

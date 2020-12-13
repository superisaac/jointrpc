package server

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	//"fmt"
	datadir "github.com/superisaac/rpctube/datadir"
	intf "github.com/superisaac/rpctube/intf/tube"
	tube "github.com/superisaac/rpctube/tube"
	grpc "google.golang.org/grpc"
)

func CommandStartServer() {
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	pBind := serverFlags.String("b", "", "The grpc server address and port")
	pDatadir := serverFlags.String("d", "", "The datadir to store configs")
	//httpBind := serverFlags.String("httpd", "127.0.0.1:50056", "http address and port")

	serverFlags.Parse(os.Args[2:])
	if *pDatadir != "" {
		datadir.SetDatadir(*pDatadir)
	}

	//go StartHTTPd(*httpBind)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	StartServer(ctx, *pBind)
}

func StartServer(ctx context.Context, bind string) {
	cfg := datadir.GetConfig()
	if bind == "" {
		bind = cfg.Server.Bind
	}
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Panicf("failed to listen: %w", err)
	} else {
		log.Infof("entry server listen at %s", bind)
	}

	tube.Tube().Start(ctx)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	//hello.RegisterHelloServer(grpcServer, s)
	intf.RegisterJSONRPCTubeServer(grpcServer, NewJSONRPCTubeServer())
	grpcServer.Serve(lis)
}

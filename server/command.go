package server

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	//"fmt"
	datadir "github.com/superisaac/rpctube/datadir"
	intf "github.com/superisaac/rpctube/intf/tube"
	tube "github.com/superisaac/rpctube/tube"
	grpc "google.golang.org/grpc"
)

func CommandStartServer() {
	datadir.GetConfig()
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	bind := serverFlags.String("bind", "127.0.0.1:50055", "The grpc server address and port")
	//httpBind := serverFlags.String("httpd", "127.0.0.1:50056", "http address and port")

	serverFlags.Parse(os.Args[2:])
	//lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bind, *port)

	//go StartHTTPd(*httpBind)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	StartServer(ctx, *bind)
}

func StartServer(ctx context.Context, bind string) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("entry server listen at %s", bind)
	}

	tube.Tube().Start(ctx)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	//hello.RegisterHelloServer(grpcServer, s)
	intf.RegisterJSONRPCTubeServer(grpcServer, NewJSONRPCTubeServer())
	grpcServer.Serve(lis)
}

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

func StartEntrypoint() {
	datadir.GetConfig()
	entryCmd := flag.NewFlagSet("node", flag.ExitOnError)
	bind := entryCmd.String("bind", "127.0.0.1:50055", "The grpc server address and port")
	httpBind := entryCmd.String("httpd", "127.0.0.1:50056", "http address and port")

	entryCmd.Parse(os.Args[2:])
	//lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bind, *port))
	lis, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("entry server listen at %s", *bind)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tube.Tube().Start(ctx)

	go StartHTTPd(*httpBind)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	//hello.RegisterHelloServer(grpcServer, s)
	intf.RegisterJSONRPCTubeServer(grpcServer, NewJSONRPCTubeServer())
	grpcServer.Serve(lis)
}

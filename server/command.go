package server

import (
	"os"
	context "context"
	"flag"
	"fmt"
	intf "github.com/superisaac/rpctube/intf/tube"
	tube "github.com/superisaac/rpctube/tube"
	grpc "google.golang.org/grpc"
	"log"
	"net"
)

func StartEntrypoint() {
	entryCmd := flag.NewFlagSet("entry", flag.ExitOnError)
	port := entryCmd.Int("port", 50055, "The server port")

	entryCmd.Parse(os.Args[2:])
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	tube.Tube().Start(ctx)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	//hello.RegisterHelloServer(grpcServer, s)
	intf.RegisterJSONRPCTubeServer(grpcServer, NewJSONRPCTubeServer())
	grpcServer.Serve(lis)
}


package server

import (
	context "context"
	"flag"
	"fmt"
	intf "github.com/superisaac/rpctube/intf/tube"
	tube "github.com/superisaac/rpctube/tube"
	grpc "google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func StartEntrypoint() {
	entryCmd := flag.NewFlagSet("entry", flag.ExitOnError)
	port := entryCmd.Int("port", 50055, "The server port")
	bind := entryCmd.String("bind", "localhost", "The server bind address")

	entryCmd.Parse(os.Args[2:])
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bind, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("entry server listen at %s:%d", *bind, *port)
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

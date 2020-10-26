package main

import (
	context "context"
	"flag"
	"fmt"
	hello "github.com/superisaac/rpctube/intf/hello"
	intf "github.com/superisaac/rpctube/intf/tube"
	server "github.com/superisaac/rpctube/server"
	tube "github.com/superisaac/rpctube/tube"
	"google.golang.org/grpc"
	"log"
	"net"
)

type HelloServer struct {
	hello.UnimplementedHelloServer
}

func (self *HelloServer) Greeting(context context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	resp := &hello.HelloResponse{Msg: "www"}
	return resp, nil
}

var (
	port = flag.Int("port", 50055, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tube.Tube().Start(ctx)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	s := &HelloServer{}
	hello.RegisterHelloServer(grpcServer, s)
	intf.RegisterJSONRPCTubeServer(grpcServer, server.NewJSONRPCTubeServer())
	grpcServer.Serve(lis)

}

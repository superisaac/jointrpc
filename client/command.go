package client

import (
	"os"
	"flag"
	"fmt"
	"log"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	//simplejson "github.com/bitly/go-simplejson"	
	grpc "google.golang.org/grpc"
)

func printHelp() {
	fmt.Println("method params...")
}

func CommandCallRPC() error {
	serverAddress := flag.String("address", "localhost:50055", "the tube server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	if flag.NArg() < 1 {
		printHelp()
		os.Exit(1)
	}

	args := flag.Args()
	method := args[0]
	clParams := args[1:len(args)]

	params, err := jsonrpc.GuessJsonArray(clParams)
	if err != nil {
		return err
	}
	
	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	
	c := intf.NewJSONRPCTubeClient(conn)
	res, err := CallRPC(c, method, params)
	if err != nil {
		return nil
	}

	repr, err := res.EncodePretty()
	if err != nil {
		return nil
	}
	fmt.Printf("%s\n", repr)
	return nil
}

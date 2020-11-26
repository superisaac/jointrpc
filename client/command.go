package client

import (
	"os"
	"context"
	"flag"
	"fmt"
	"log"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	simplejson "github.com/bitly/go-simplejson"	
	grpc "google.golang.org/grpc"
)

func printHelp() {
	fmt.Println("method params...")
}

func CallRPC() error {
	serverAddress := flag.String("address", "localhost:50055", "the tube server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	if flag.NArg() < 1 {
		printHelp()
		os.Exit(1)
	}

	args := flag.Args()
	method := args[0]
	params := args[1:len(args)]

	parsedParams, err := jsonrpc.GuessJsonArray(params)
	if err != nil {
		return err
	}
	paramsJson := simplejson.New()
	paramsJson.SetPath(nil, parsedParams)

	paramsStr, err := jsonrpc.MarshalJson(paramsJson)
	if err != nil {
		return err
	}
	

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	
	c := intf.NewJSONRPCTubeClient(conn)

	req := &intf.JSONRPCRequest{
		Id: 1,
		Method: method,
		Params: paramsStr}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.Call(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("call res %v", res)
	return nil
}

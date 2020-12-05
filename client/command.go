package client

import (
	"flag"
	"fmt"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	"log"
	"os"
	//simplejson "github.com/bitly/go-simplejson"
	grpc "google.golang.org/grpc"
)

func printHelp() {
	fmt.Println("method params...")
}

func CommandCallRPC() {
	callFlags := flag.NewFlagSet("rpc", flag.ExitOnError)

	serverAddress := callFlags.String("address", "localhost:50055", "the tube server address")
	callFlags.Parse(os.Args[2:])
	log.Printf("dial server %s", *serverAddress)

	if callFlags.NArg() < 1 {
		printHelp()
		os.Exit(1)
	}

	args := callFlags.Args()
	method := args[0]
	clParams := args[1:len(args)]

	params, err := jsonrpc.GuessJsonArray(clParams)
	if err != nil {
		panic(err)
	}

	err = RunCallRPC(*serverAddress, method, params)
	if err != nil {
		panic(err)
	}
}

func RunCallRPC(serverAddress string, method string, params []interface{}) error {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c := intf.NewJSONRPCTubeClient(conn)
	res, err := CallRPC(c, method, params)
	if err != nil {
		return err
	}

	repr, err := res.EncodePretty()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", repr)
	//repr = res.MustString()
	//fmt.Printf("%s\n", repr)
	return nil
}

// Call ListMethods
func CommandListMethods() {
	listMethodsFlags := flag.NewFlagSet("listmethods", flag.ExitOnError)

	serverAddress := listMethodsFlags.String("address", "localhost:50055", "the tube server address")
	listMethodsFlags.Parse(os.Args[2:])
	log.Printf("dial server %s", *serverAddress)
	err := RunListMethods(*serverAddress)
	if err != nil {
		panic(err)
	}
}

func RunListMethods(serverAddress string) error {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c := intf.NewJSONRPCTubeClient(conn)
	methods, err := ListMethods(c)
	if err != nil {
		return nil
	}

	for _, m := range methods {
		fmt.Printf("%s\n", m)
	}
	return nil
}

// Example FIFO
func CommandExampleFIFO() {
	examFlags := flag.NewFlagSet("example.fifo", flag.ExitOnError)
	serverAddress := examFlags.String("address", "localhost:50055", "the tube server address")
	examFlags.Parse(os.Args[2:])
	log.Printf("dial server %s", *serverAddress)
	err := StartExampleFIFO(*serverAddress)
	if err != nil {
		panic(err)
	}
}

func StartExampleFIFO(serverAddress string) error {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := intf.NewJSONRPCTubeClient(conn)

	err = ExampleFIFO(c)
	if err != nil {
		return err
	}
	return nil
}

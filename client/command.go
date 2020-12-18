package client

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	//intf "github.com/superisaac/rpctube/intf/tube"
	handler "github.com/superisaac/rpctube/handler"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	//example "github.com/superisaac/rpctube/client/example"
	//grpc "google.golang.org/grpc"
)

func printHelp() {
	fmt.Println("method params...")
}

func tryGetServerAddress(serverAddress string) string {
	if serverAddress == "" {
		serverAddress = os.Getenv("TUBE_CONNECT")
	}

	if serverAddress == "" {
		serverAddress = "localhost:50055"
	}
	return serverAddress
}

// Send Notify
func CommandSendNotify() {
	callFlags := flag.NewFlagSet("notify", flag.ExitOnError)
	pAddress := callFlags.String("c", "", "the server address to connect, default 127.0.0.1:50055")
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")

	callFlags.Parse(os.Args[2:])

	serverAddress := tryGetServerAddress(*pAddress)

	log.Infof("dial server %s", serverAddress)

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

	err = RunSendNotify(serverAddress, method, params, *pBroadcast)
	if err != nil {
		panic(err)
	}
}

func RunSendNotify(serverAddress string, method string, params []interface{}, broadcast bool) error {
	client := NewRPCClient(serverAddress)
	err := client.Connect()
	if err != nil {
		return err
	}
	err = client.SendNotify(method, params, broadcast)
	if err != nil {
		return err
	}
	return nil
}

// Call RPC
func CommandCallRPC(subcmd string) {
	callFlags := flag.NewFlagSet(subcmd, flag.ExitOnError)
	pAddress := callFlags.String("c", "", "the server address to connect, default 127.0.0.1:50055")

	callFlags.Parse(os.Args[2:])

	serverAddress := tryGetServerAddress(*pAddress)

	log.Infof("dial server %s", serverAddress)

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

	err = RunCallRPC(serverAddress, method, params)
	if err != nil {
		panic(err)
	}
}

func RunCallRPC(serverAddress string, method string, params []interface{}) error {
	client := NewRPCClient(serverAddress)
	err := client.Connect()
	if err != nil {
		return err
	}
	res, err := client.CallRPC(method, params)
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
	pAddress := listMethodsFlags.String("c", "", "the tube server address")
	listMethodsFlags.Parse(os.Args[2:])

	serverAddress := tryGetServerAddress(*pAddress)

	log.Infof("dial server %s", serverAddress)
	err := RunListMethods(serverAddress)
	if err != nil {
		panic(err)
	}
}

func RunListMethods(serverAddress string) error {
	client := NewRPCClient(serverAddress)
	err := client.Connect()
	if err != nil {
		return err
	}
	methodInfos, err := client.ListMethods()
	if err != nil {
		return nil
	}

	fmt.Printf("available methods:\n")
	for _, minfo := range methodInfos {
		fmt.Printf("  %s\t%s\n", minfo.Name, minfo.Help)
	}
	return nil
}

// Watch notify
func CommandWatch() {
	subFlags := flag.NewFlagSet("watchnotify", flag.ExitOnError)
	pAddress := subFlags.String("c", "", "the tube server address")
	subFlags.Parse(os.Args[2:])
	serverAddress := tryGetServerAddress(*pAddress)

	notifyNames := subFlags.Args()
	if len(notifyNames) < 1 {
		panic(errors.New("No notify methods specified to watch"))
	}
	rpcClient := NewRPCClient(serverAddress)

	for _, notifyName := range notifyNames {
		rpcClient.On(notifyName, func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
			repr, err := req.Message.EncodePretty()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", repr)
			return nil, nil
		})
	}
	err := rpcClient.Connect()
	if err != nil {
		panic(err)
	}
	err = rpcClient.RunHandlers()
	if err != nil {
		panic(err)
	}
}

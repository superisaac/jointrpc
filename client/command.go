package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	//"strings"
	simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	handler "github.com/superisaac/rpctube/tube/handler"
	"os"
	//example "github.com/superisaac/rpctube/client/example"
	//grpc "google.golang.org/grpc"
)

func printHelp() {
	fmt.Println("method params...")
}

// Send Notify
func CommandSendNotify() {
	callFlags := flag.NewFlagSet("notify", flag.ExitOnError)
	pAddress := callFlags.String("c", "", "the server address to connect, default 127.0.0.1:50055")
	pCertFile := callFlags.String("cert", "", "the cert file, default empty")
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")

	callFlags.Parse(os.Args[2:])

	serverAddress, certFile := TryGetServerSettings(*pAddress, *pCertFile)

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

	err = RunSendNotify(serverAddress, certFile, method, params, *pBroadcast)
	if err != nil {
		panic(err)
	}
}

func RunSendNotify(serverAddress string, certFile string, method string, params []interface{}, broadcast bool) error {
	client := NewRPCClient(serverAddress, certFile)
	err := client.Connect()
	if err != nil {
		return err
	}
	err = client.SendNotify(context.Background(), method, params, broadcast)
	if err != nil {
		return err
	}
	return nil
}

// Call RPC
func CommandCallRPC(subcmd string) {
	callFlags := flag.NewFlagSet(subcmd, flag.ExitOnError)
	pAddress := callFlags.String("c", "", "the server address to connect, default 127.0.0.1:50055")
	pCertFile := callFlags.String("cert", "", "the cert file, default empty")

	callFlags.Parse(os.Args[2:])

	serverAddress, certFile := TryGetServerSettings(*pAddress, *pCertFile)

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

	err = RunCallRPC(serverAddress, certFile, method, params)
	if err != nil {
		panic(err)
	}
}

func RunCallRPC(serverAddress string, certFile string, method string, params []interface{}) error {
	client := NewRPCClient(serverAddress, certFile)
	err := client.Connect()
	if err != nil {
		return err
	}
	res, err := client.CallRPC(context.Background(), method, params)
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
	pCertFile := listMethodsFlags.String("cert", "", "the cert file, default empty")
	listMethodsFlags.Parse(os.Args[2:])

	serverAddress, certFile := TryGetServerSettings(*pAddress, *pCertFile)

	log.Infof("dial server %s", serverAddress)
	err := RunListMethods(serverAddress, certFile)
	if err != nil {
		panic(err)
	}
}

func RunListMethods(serverAddress string, certFile string) error {
	client := NewRPCClient(serverAddress, certFile)
	err := client.Connect()
	if err != nil {
		return err
	}
	methodInfos, err := client.ListMethods(context.Background())
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
	pCertFile := subFlags.String("cert", "", "the cert file, default empty")
	subFlags.Parse(os.Args[2:])

	notifyNames := subFlags.Args()
	if len(notifyNames) < 1 {
		panic(errors.New("No notify methods specified to watch"))
	}
	serverAddress, certFile := TryGetServerSettings(*pAddress, *pCertFile)
	rpcClient := NewRPCClient(serverAddress, certFile)

	for _, notifyName := range notifyNames {
		rpcClient.On(notifyName, func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
			msg := req.MsgVec.Msg
			repr, err := msg.EncodePretty()
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

	err = rpcClient.Handle(context.Background())
	if err != nil {
		panic(err)
	}
}

// Watch methods update
func CommandWatchMethods() {
	subFlags := flag.NewFlagSet("watchmethods", flag.ExitOnError)
	pAddress := subFlags.String("c", "", "the tube server address")
	pCertFile := subFlags.String("cert", "", "the cert file, default empty")
	pVerbose := subFlags.Bool("verbose", false, "show method info")
	subFlags.Parse(os.Args[2:])

	serverAddress, certFile := TryGetServerSettings(*pAddress, *pCertFile)
	rpcClient := NewRPCClient(serverAddress, certFile)

	rpcClient.Connect()

	ch, err := rpcClient.WatchMethods(context.Background())
	if err != nil {
		panic(err)
	}
	for {
		select {
		case update, ok := <-ch:
			if !ok {
				return
			}
			if *pVerbose {
				printMethodInfos(update)
			} else {
				printMethodNames(update)
			}
		}
	}
}

func printMethodInfos(update []*intf.MethodInfo) {
	var arr [](map[string](interface{}))
	for _, info := range update {
		mapInfo := map[string](interface{}){
			"name": info.Name,
		}
		if info.Help != "" {
			mapInfo["help"] = info.Help
		}
		if info.SchemaJson != "" {
			schemaJson, err := simplejson.NewJson([]byte(info.SchemaJson))
			if err != nil {
				panic(err)
			}
			mapInfo["schema"] = schemaJson.Interface()
		}
		// TODO: schema
		arr = append(arr, mapInfo)
	}
	jarr := simplejson.New()
	jarr.SetPath(nil, arr)
	repr, err := jarr.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(repr))
}

func printMethodNames(update []*intf.MethodInfo) {
	var arr []string
	for _, info := range update {
		arr = append(arr, info.Name)
	}
	jarr := simplejson.New()
	jarr.SetPath(nil, arr)
	repr, err := jarr.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(repr))
}

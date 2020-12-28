package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	//"strings"
	simplejson "github.com/bitly/go-simplejson"
	//log "github.com/sirupsen/logrus"
	//intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
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
	serverFlag := NewServerFlag(callFlags)
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")

	callFlags.Parse(os.Args[2:])

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

	err = RunSendNotify(serverFlag.Get(), method, params, *pBroadcast)
	if err != nil {
		panic(err)
	}
}

func RunSendNotify(serverEntry ServerEntry, method string, params []interface{}, broadcast bool) error {
	client := NewRPCClient(serverEntry)
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
	serverFlag := NewServerFlag(callFlags)
	callFlags.Parse(os.Args[2:])

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

	err = RunCallRPC(serverFlag.Get(), method, params)
	if err != nil {
		panic(err)
	}
}

func RunCallRPC(serverEntry ServerEntry, method string, params []interface{}) error {
	client := NewRPCClient(serverEntry)
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
	serverFlag := NewServerFlag(listMethodsFlags)
	listMethodsFlags.Parse(os.Args[2:])

	err := RunListMethods(serverFlag.Get())
	if err != nil {
		panic(err)
	}
}

func RunListMethods(serverEntry ServerEntry) error {
	client := NewRPCClient(serverEntry)
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
	serverFlag := NewServerFlag(subFlags)
	subFlags.Parse(os.Args[2:])

	notifyNames := subFlags.Args()
	if len(notifyNames) < 1 {
		panic(errors.New("No notify methods specified to watch"))
	}

	rpcClient := NewRPCClient(serverFlag.Get())

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
func CommandWatchState() {
	subFlags := flag.NewFlagSet("watchstate", flag.ExitOnError)
	serverFlag := NewServerFlag(subFlags)
	pVerbose := subFlags.Bool("verbose", false, "show method info")
	subFlags.Parse(os.Args[2:])

	rpcClient := NewRPCClient(serverFlag.Get())

	if *pVerbose {
		rpcClient.OnStateChange(printMethodInfos)
	} else {
		rpcClient.OnStateChange(printMethodNames)
	}

	rpcClient.Connect()
	rpcClient.Handle(context.Background())
}

func printMethodInfos(state *tube.TubeState) {
	var arr [](map[string](interface{}))
	for _, info := range state.Methods {
		mapInfo := map[string](interface{}){
			"name": info.Name,
		}
		if info.Help != "" {
			mapInfo["help"] = info.Help
		}

		if info.Schema() != nil {
			mapInfo["schema"] = info.Schema().RebuildType()
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

func printMethodNames(state *tube.TubeState) {
	var arr []string
	for _, info := range state.Methods {
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

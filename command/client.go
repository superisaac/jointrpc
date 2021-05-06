package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	//"strings"

	simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"os"
	//example "github.com/superisaac/jointrpc/client/example"
	//grpc "google.golang.org/grpc"
)

// Send Notify
func CommandSendNotify() {
	callFlags := flag.NewFlagSet("notify", flag.ExitOnError)
	serverFlag := client.NewServerFlag(callFlags)
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")
	pTraceId := callFlags.String("traceid", "", "trace id during the workflow")

	callFlags.Parse(os.Args[2:])

	if callFlags.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "method params...\n")
		os.Exit(1)
	}

	args := callFlags.Args()
	method := args[0]
	clParams := args[1:len(args)]

	params, err := jsonrpc.GuessJsonArray(clParams)
	if err != nil {
		panic(err)
	}

	err = RunSendNotify(serverFlag.Get(), method, params,
		client.WithBroadcast(*pBroadcast), client.WithTraceId(*pTraceId))
	if err != nil {
		panic(err)
	}
}

func RunSendNotify(serverEntry client.ServerEntry, method string, params []interface{}, opts ...client.CallOptionFunc) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}
	err = c.SendNotify(context.Background(), method, params, opts...)
	if err != nil {
		return err
	}
	return nil
}

// Call RPC
func CommandCallRPC(subcmd string) {
	callFlags := flag.NewFlagSet(subcmd, flag.ExitOnError)
	serverFlag := client.NewServerFlag(callFlags)
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")
	pTraceId := callFlags.String("traceid", "", "trace id during the workflow")

	callFlags.Parse(os.Args[2:])
	// TODO, check the sanity agains traceId

	if callFlags.NArg() < 1 {
		fmt.Println("method params...")
		os.Exit(1)
	}

	args := callFlags.Args()
	method := args[0]
	clParams := args[1:len(args)]

	params, err := jsonrpc.GuessJsonArray(clParams)
	if err != nil {
		panic(err)
	}

	err = RunCallRPC(serverFlag.Get(), method, params,
		client.WithBroadcast(*pBroadcast), client.WithTraceId(*pTraceId))
	if err != nil {
		panic(err)
	}
}

func RunCallRPC(serverEntry client.ServerEntry, method string, params []interface{}, opts ...client.CallOptionFunc) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}

	res, err := c.CallRPC(context.Background(), method, params, opts...)
	if err != nil {
		return err
	}
	log.Debugf("result got trace id %s", res.TraceId())
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
	aFlags := flag.NewFlagSet("methods", flag.ExitOnError)
	serverFlag := client.NewServerFlag(aFlags)
	aFlags.Parse(os.Args[2:])

	err := RunListMethods(serverFlag.Get())
	if err != nil {
		panic(err)
	}
}

func RunListMethods(serverEntry client.ServerEntry) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}
	methodInfos, err := c.ListMethods(context.Background())
	if err != nil {
		return nil
	}

	fmt.Printf("available methods:\n")
	for _, minfo := range methodInfos {
		fmt.Printf("  %s\t%s\n", minfo.Name, minfo.Help)
	}
	return nil
}

// Call ListDelegates
func CommandListDelegates() {
	aFlags := flag.NewFlagSet("delegates", flag.ExitOnError)
	serverFlag := client.NewServerFlag(aFlags)
	aFlags.Parse(os.Args[2:])

	err := RunListDelegates(serverFlag.Get())
	if err != nil {
		panic(err)
	}
}

func RunListDelegates(serverEntry client.ServerEntry) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}
	delegates, err := c.ListDelegates(context.Background())
	if err != nil {
		return nil
	}

	fmt.Printf("delegated methods:\n")
	for _, name := range delegates {
		fmt.Printf("  %s\n", name)
	}
	return nil
}

// Watch notify
func CommandWatch() {
	subFlags := flag.NewFlagSet("watchnotify", flag.ExitOnError)
	serverFlag := client.NewServerFlag(subFlags)
	subFlags.Parse(os.Args[2:])

	notifyNames := subFlags.Args()
	if len(notifyNames) < 1 {
		panic(errors.New("No notify methods specified to watch"))
	}

	rpcClient := client.NewRPCClient(serverFlag.Get())

	disp := dispatch.NewDispatcher()

	for _, notifyName := range notifyNames {
		disp.On(notifyName, func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
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

	err = rpcClient.Worker(context.Background(), disp)
	if err != nil {
		panic(err)
	}
}

// Watch methods update
func CommandWatchState() {
	subFlags := flag.NewFlagSet("watchstate", flag.ExitOnError)
	serverFlag := client.NewServerFlag(subFlags)
	pVerbose := subFlags.Bool("verbose", false, "show method info")
	subFlags.Parse(os.Args[2:])

	rpcClient := client.NewRPCClient(serverFlag.Get())
	stateDisp := dispatch.NewStateDispatcher()

	if *pVerbose {
		stateDisp.OnStateChange(printMethodInfos)
	} else {
		stateDisp.OnStateChange(printMethodNames)
	}

	rpcClient.Connect()
	rpcClient.SubscribeState(context.Background(), stateDisp)
}

func printMethodInfos(state *rpcrouter.ServerState) {
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

func printMethodNames(state *rpcrouter.ServerState) {
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

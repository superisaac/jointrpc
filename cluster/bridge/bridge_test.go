package bridge

import (
	//"fmt"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	client "github.com/superisaac/jointrpc/client"
	server "github.com/superisaac/jointrpc/server"	
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//mirror "github.com/superisaac/jointrpc/cluster/mirror"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
	"testing"
	"time"
)

const addSchema = `
{
  "type": "method",
  "params": [
    {
      "type": "number"
    },
    {
      "type": "number"
    }
  ],
  "returns": {
    "type": "number"
  }
}
`

func TestBridgeRun(t *testing.T) {
	assert := assert.New(t)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// start server1
	go server.StartServer(rootCtx, "localhost:10020", nil)
	// start server2
	go server.StartServer(rootCtx, "localhost:10021", nil)
	// start server3
	go server.StartServer(rootCtx, "localhost:10022", nil)

	time.Sleep(100 * time.Millisecond)

	serverEntries := []client.ServerEntry{{"h2c://localhost:10020", ""}, {"h2c://localhost:10021", ""}, {"h2c://localhost:10022", ""}}
		
	go StartNewBridge(rootCtx, serverEntries)

	// start client1, the serve of add2int()
	c1 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10020", ""})
	c1.On("add2int", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		a := jsonrpc.MustInt(params[0], "params[0]")
		b := jsonrpc.MustInt(params[1], "params[1]")
		return a + b, nil
	}, handler.WithSchema(addSchema))
	err := c1.Connect()
	assert.Nil(err)
	cCtx, cancelServo := context.WithCancel(context.Background())
	//defer cancelServo()
	go c1.Handle(cCtx)

	// start c2, the add2int() caller to server2
	time.Sleep(100 * time.Millisecond)
	c2 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10021", ""})
	err = c2.Connect()
	assert.Nil(err)

	// call rpc from server2 which delegates server1
	time.Sleep(100 * time.Millisecond)
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	delegates, err := c2.ListDelegates(ctx1)
	assert.Nil(err)
	assert.Equal([]string{"add2int"}, delegates)
	res, err := c2.CallRPC(ctx1, "add2int", [](interface{}){5, 6})
	assert.Nil(err)
	assert.True(res.IsResult())
	assert.Equal(json.Number("11"), res.MustResult())

	// start client3
	c3 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10022", ""})
	err = c3.Connect()
	assert.Nil(err)
	// call rpc from server3 which doesnot delegates server1
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	res3, err := c3.CallRPC(ctx2, "add2int", [](interface{}){15, 16})
	assert.Nil(err)
	assert.True(res3.IsResult())
	assert.Equal(json.Number("31"), res3.MustResult())	

	// close client serve, the add2int() provider
	cancelServo()
	time.Sleep(100 * time.Millisecond)
	d1, err := c2.ListDelegates(ctx1)
	assert.Nil(err)
	assert.Equal(0, len(d1))
	//assert.Equal([]string{}, d1)

	// call rpc from server2 which doesnot delegates server1
	ctx4, cancel4 := context.WithCancel(context.Background())
	defer cancel4()
	res4, err := c2.CallRPC(ctx4, "add2int", [](interface{}){15, 16})
	assert.Nil(err)
	assert.True(res4.IsError())
	errBody4, ok := res4.MustError().(map[string]interface{})
	assert.True(ok)
	assert.Equal(json.Number("404"), errBody4["code"])
	assert.Equal("method not found", errBody4["reason"])
}

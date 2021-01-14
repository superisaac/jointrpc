package bridge

import (
	//"fmt"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"

	client "github.com/superisaac/jointrpc/client"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	server "github.com/superisaac/jointrpc/server"
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

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

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

	time.Sleep(200 * time.Millisecond)

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
	time.Sleep(200 * time.Millisecond)
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
	res, err := c2.CallRPC(ctx1, "add2int", [](interface{}){5, 6},
		client.WithTraceId("trace1"))
	assert.Nil(err)
	assert.Equal("trace1", res.TraceId())
	assert.True(res.IsResult())
	assert.Equal(json.Number("11"), res.MustResult())

	// start client3
	c3 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10022", ""})
	err = c3.Connect()
	assert.Nil(err)
	// call rpc from server3 which doesnot delegates server1
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	res3, err := c3.CallRPC(ctx2, "add2int", [](interface{}){15, 16},
		client.WithTraceId("trace3"))
	assert.Nil(err)
	assert.Equal("trace3", res3.TraceId())
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
	res4, err := c2.CallRPC(ctx4, "add2int", [](interface{}){15, 16},
		client.WithTraceId("trace4"))
	assert.Nil(err)
	assert.Equal("trace4", res4.TraceId())
	assert.True(res4.IsError())
	errBody4, ok := res4.MustError().(map[string]interface{})
	assert.True(ok)
	assert.Equal(json.Number("404"), errBody4["code"])
	assert.Equal("method not found", errBody4["reason"])
}

func TestServerBreak(t *testing.T) {
	assert := assert.New(t)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start server1
	ctx1, cancelServer1 := context.WithCancel(rootCtx)
	go server.StartServer(ctx1, "localhost:10030", nil)
	// start server2
	go server.StartServer(rootCtx, "localhost:10031", nil)
	// start server3
	go server.StartServer(rootCtx, "localhost:10032", nil)

	time.Sleep(100 * time.Millisecond)

	serverEntries := []client.ServerEntry{{"h2c://localhost:10030", ""}, {"h2c://localhost:10031", ""}, {"h2c://localhost:10032", ""}}

	go StartNewBridge(rootCtx, serverEntries)

	// start client1, the serve of add2int()
	c1 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10030", ""})
	c1.On("add2int", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		a := jsonrpc.MustInt(params[0], "params[0]")
		b := jsonrpc.MustInt(params[1], "params[1]")
		return a + b, nil
	}, handler.WithSchema(addSchema))
	err := c1.Connect()
	assert.Nil(err)
	cCtx, cancelServo := context.WithCancel(context.Background())
	defer cancelServo()
	go c1.Handle(cCtx)
	// start c2, the add2int() caller to server2
	time.Sleep(100 * time.Millisecond)
	c2 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10031", ""})
	err = c2.Connect()
	assert.Nil(err)

	// call rpc from server2 which delegates server1
	time.Sleep(100 * time.Millisecond)
	ctx2, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	delegates, err := c2.ListDelegates(ctx1)
	assert.Nil(err)
	assert.Equal([]string{"add2int"}, delegates)
	res, err := c2.CallRPC(ctx2, "add2int", [](interface{}){5, 6},
		client.WithTraceId("trace11"))
	assert.Nil(err)
	assert.Equal("trace11", res.TraceId())
	assert.True(res.IsResult())
	assert.Equal(json.Number("11"), res.MustResult())
	// close server1
	cancelServer1()
	time.Sleep(100 * time.Millisecond)

	delegates1, err := c2.ListDelegates(context.Background())
	assert.Nil(err)
	assert.Nil(delegates1)

	ctx3, cancel3 := context.WithCancel(context.Background())
	defer cancel3()

	resf, err := c2.CallRPC(ctx3, "add2int", [](interface{}){8, 12},
		client.WithTraceId("tracebag"))
	assert.Nil(err)
	assert.Equal("tracebag", resf.TraceId())
	assert.True(resf.IsError())

	// start client3
	c3 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10032", ""})
	err = c3.Connect()
	assert.Nil(err)
	// call rpc from server3 which doesnot delegates server1
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	res3, err := c3.CallRPC(
		ctx2, "add2int", [](interface{}){15, 16},
		client.WithTraceId("trace3"))
	assert.Nil(err)
	assert.Equal("trace3", res3.TraceId())
	assert.True(res3.IsError())

	errBody3, ok := res3.MustError().(map[string]interface{})
	assert.True(ok)
	assert.Equal(json.Number("404"), errBody3["code"])
	assert.Equal("method not found", errBody3["reason"])

}

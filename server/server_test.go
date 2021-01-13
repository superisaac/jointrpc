package server

import (
	//"fmt"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	client "github.com/superisaac/jointrpc/client"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	assert := assert.New(t)

	root, cancelRoot := context.WithCancel(context.Background())
	defer cancelRoot()

	c1 := context.WithValue(root, "key1", "value1")
	assert.Equal("value1", c1.Value("key1"))
	c2, cancelC2 := context.WithCancel(c1)
	defer cancelC2()
	assert.Equal("value1", c2.Value("key1"))
}

func TestServerClientRound(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartServer(ctx, "127.0.0.1:10001", nil)

	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"h2c://127.0.0.1:10001", ""})
	err := c.Connect()
	assert.Nil(err)

	res, err := c.CallRPC(ctx, ".echo", [](interface{}){"nice"}, client.WithTraceId("trace1"))
	assert.Nil(err)
	assert.Equal("trace1", res.TraceId())

	assert.True(res.IsResult())
	m, ok := res.MustResult().(map[string]interface{})
	assert.True(ok)
	assert.Equal("nice", m["echo"])

	res1, err := c.CallRPC(ctx, ".echo", [](interface{}){1}, client.WithTraceId("trace2"))
	assert.Nil(err)
	assert.Equal("trace2", res1.TraceId())
	assert.True(res1.IsError())
	errbody, ok := res1.MustError().(map[string]interface{})
	assert.True(ok)
	assert.Equal("Validation Error: .params[0] data is not string", errbody["reason"])
}

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

func StartTestServe(rootCtx context.Context, serverUrl string, whoami string) {
	c := client.NewRPCClient(client.ServerEntry{serverUrl, ""})
	c.On("add2int", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		a := jsonrpc.MustInt(params[0], "params[0]")
		b := jsonrpc.MustInt(params[1], "params[1]")
		return a + b, nil
	}, handler.WithSchema(addSchema))

	c.On("whoami", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		return whoami, nil
	})
	c.Connect()
	c.Handle(rootCtx)
}

func TestClientAsServe(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go StartServer(ctx, "127.0.0.1:10002", nil)
	time.Sleep(100 * time.Millisecond)

	go StartTestServe(ctx, "h2c://127.0.0.1:10002", "testclent")
	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"h2c://127.0.0.1:10002", ""})
	err := c.Connect()
	assert.Nil(err)

	res, err := c.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace11"))
	assert.Nil(err)
	assert.Equal("trace11", res.TraceId())
	assert.True(res.IsResult())
	assert.Equal(json.Number("11"), res.MustResult())

}

func TestBroadcastRequest(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go StartServer(ctx, "127.0.0.1:10005", nil)
	time.Sleep(100 * time.Millisecond)

	go StartTestServe(ctx, "h2c://127.0.0.1:10005", "jack")
	time.Sleep(100 * time.Millisecond)
	go StartTestServe(ctx, "h2c://127.0.0.1:10005", "mike")
	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"h2c://127.0.0.1:10005", ""})
	err := c.Connect()
	assert.Nil(err)

	res, err := c.CallRPC(ctx, "whoami", [](interface{}){},
		client.WithBroadcast(true),
		client.WithTraceId("trace41"))
	assert.Nil(err)
	assert.Equal("trace41", res.TraceId())
	assert.True(res.IsResult())

	arr, ok := res.MustResult().([]interface{})
	assert.True(ok)
	assert.Equal(2, len(arr))

	r1, ok := arr[0].(map[string]interface{})
	assert.True(ok)
	who, ok := r1["result"]
	assert.True(ok)
	assert.True(who == "jack" || who == "mike")
}

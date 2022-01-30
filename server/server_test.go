package server

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	//"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/datadir"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonz"
	"io/ioutil"

	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestAdhocContext(t *testing.T) {
	assert := assert.New(t)

	root, cancelRoot := context.WithCancel(context.Background())
	defer cancelRoot()

	c1 := context.WithValue(root, "key1", "value1")
	assert.Equal("value1", c1.Value("key1"))
	c2, cancelC2 := context.WithCancel(c1)
	defer cancelC2()
	assert.Equal("value1", c2.Value("key1"))
}

const addSchema = `
{
  "type": "method",
  "params": ["number", "number"],
  "returns": "number"
}
`

func StartTestServe(rootCtx context.Context, serverUrl string, whoami string) {
	c := client.NewRPCClient(client.ServerEntry{serverUrl, ""})
	disp := dispatch.NewDispatcher()

	disp.On("add2int", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		a := jsonz.MustInt(params[0], "params[0]")
		b := jsonz.MustInt(params[1], "params[1]")
		return a + b, nil
	}, dispatch.WithSchema(addSchema))

	disp.On("fakeadd2int", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		return "not a number", nil
	}, dispatch.WithSchema(addSchema))

	disp.On("whoami", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		return whoami, nil
	})
	c.Connect()
	c.Live(rootCtx, disp)
}

func TestClientAsServe(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go StartGRPCServer(ctx, "127.0.0.1:10002")
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
	// invalid schema tests
	res1, err := c.CallRPC(ctx, "add2int", [](interface{}){5, "ss1"}, client.WithTraceId("trace13"))
	assert.Nil(err)
	assert.Equal("trace13", res1.TraceId())
	assert.True(res1.IsError())
	//errbody, ok := res1.MustError().(map[string](interface{}))
	errbody := res1.MustError()
	assert.Equal(10901, errbody.Code)
	assert.Equal("Validation Error: .params[1] data is not number", errbody.Message)

	// invalid schema tests
	res2, err := c.CallRPC(ctx, "fakeadd2int", [](interface{}){5, 8}, client.WithTraceId("trace73"))
	assert.Nil(err)
	assert.Equal("trace73", res2.TraceId())
	assert.True(res2.IsError())
	//errbody2, ok := res2.MustError().(map[string](interface{}))
	errbody2 := res2.MustError()
	assert.Equal(10901, errbody2.Code)
	assert.Equal("Validation Error: .result data is not number", errbody2.Message)
}
func TestClientAuth(t *testing.T) {
	assert := assert.New(t)

	ctx1, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	ctx := ServerContext(ctx1, nil)
	factory := rpcrouter.RouterFactoryFromContext(ctx)
	//router := factory.DefaultRouter()
	factory.Config.Authorizations = []datadir.BasicAuth{{Username: "abc", Password: "1111"}}

	go StartGRPCServer(ctx, "127.0.0.1:10092")
	time.Sleep(100 * time.Millisecond)

	go StartHTTPServer(ctx, "127.0.0.1:10093")

	go StartTestServe(ctx, "h2c://abc:1111@127.0.0.1:10092", "testclent")
	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"h2c://127.0.0.1:10092", ""})
	err := c.Connect()
	assert.Nil(err)

	_, err = c.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace11"))
	assert.NotNil(err)
	var statusErr *client.RPCStatusError
	ok := errors.As(err, &statusErr)
	assert.True(ok)
	assert.Equal(401, statusErr.Code)
	assert.Equal("auth failed", statusErr.Reason)

	c1 := client.NewRPCClient(client.ServerEntry{"h2c://abc:1111@127.0.0.1:10092", ""})
	err1 := c1.Connect()
	assert.Nil(err1)

	res1, err := c1.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace811"))
	assert.Nil(err)
	assert.Equal("trace811", res1.TraceId())
	assert.True(res1.IsResult())
	assert.Equal(json.Number("11"), res1.MustResult())

	c2 := client.NewRPCClient(client.ServerEntry{"http://abc:1111@127.0.0.1:10093", ""})
	res2, err := c2.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace811"))
	assert.Nil(err)
	assert.Equal("trace811", res2.TraceId())
	assert.True(res2.IsResult())
	assert.Equal(json.Number("11"), res2.MustResult())

	// bad password
	c3 := client.NewRPCClient(client.ServerEntry{"http://abc:1113@127.0.0.1:10093", ""})
	_, err3 := c3.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace811"))
	assert.NotNil(err3)
	assert.Equal("bad resp 401", err3.Error())

}

func TestClientAuthNamespace(t *testing.T) {
	assert := assert.New(t)

	ctx1, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	ctx := ServerContext(ctx1, nil)
	factory := rpcrouter.RouterFactoryFromContext(ctx)
	factory.Config.Authorizations = []datadir.BasicAuth{{Username: "abc", Password: "1111", Namespace: "a1"}}

	go StartGRPCServer(ctx, "127.0.0.1:10392")
	time.Sleep(100 * time.Millisecond)

	go StartHTTPServer(ctx, "127.0.0.1:10393")

	go StartTestServe(ctx, "h2c://abc:1111@127.0.0.1:10392", "testclent")
	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"h2c://127.0.0.1:10392", ""})
	err := c.Connect()
	assert.Nil(err)

	_, err = c.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace11"))
	assert.NotNil(err)
	var statusErr *client.RPCStatusError
	ok := errors.As(err, &statusErr)
	assert.True(ok)
	assert.Equal(401, statusErr.Code)
	assert.Equal("auth failed", statusErr.Reason)

	c1 := client.NewRPCClient(client.ServerEntry{"h2c://abc:1111@127.0.0.1:10392", ""})
	err1 := c1.Connect()
	assert.Nil(err1)

	res1, err := c1.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace811"))
	assert.Nil(err)
	assert.Equal("trace811", res1.TraceId())
	assert.True(res1.IsResult())
	assert.Equal(json.Number("11"), res1.MustResult())

	c2 := client.NewRPCClient(client.ServerEntry{"http://abc:1111@127.0.0.1:10393", ""})
	res2, err := c2.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace811"))
	assert.Nil(err)
	assert.Equal("trace811", res2.TraceId())
	assert.True(res2.IsResult())
	assert.Equal(json.Number("11"), res2.MustResult())

	// bad password
	c3 := client.NewRPCClient(client.ServerEntry{"http://abc:1113@127.0.0.1:10393", ""})
	_, err3 := c3.CallRPC(ctx, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace811"))
	assert.NotNil(err3)
	assert.Equal("bad resp 401", err3.Error())

}

func TestHTTPClient(t *testing.T) {
	assert := assert.New(t)

	ctx1, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	ctx := ServerContext(ctx1, nil)

	go StartGRPCServer(ctx, "127.0.0.1:10072")
	time.Sleep(100 * time.Millisecond)
	go StartHTTPServer(ctx, "127.0.0.1:10073")
	time.Sleep(100 * time.Millisecond)

	go StartTestServe(ctx, "http://127.0.0.1:10073", "testws")
	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"http://127.0.0.1:10073", ""})
	assert.True(c.IsHttp())
	assert.False(c.IsSecure())
	err := c.Connect()
	assert.Nil(err)

	res, err := c.CallRPC(ctx, "add2int", [](interface{}){9, 36}, client.WithTraceId("trace11"))
	assert.Nil(err)
	assert.Equal("trace11", res.TraceId())
	assert.True(res.IsResult())
	assert.Equal(json.Number("45"), res.MustResult())

	err = c.SendNotify(ctx, "nosuchmethod", [](interface{}){}, client.WithTraceId("trace31"))
	assert.Nil(err)
}

func TestBroadcastRequest(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go StartGRPCServer(ctx, "127.0.0.1:10005")
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

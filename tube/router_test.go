package tube

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	"testing"
)

// implements ConnT
type TestConnT struct {
	ch MsgChannel
}

func (self TestConnT) RecvChannel() MsgChannel {
	return self.ch
}

func (self TestConnT) CanBroadcast() bool {
	return true
}

func TestJoinConn(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()

	cid := CID(1002)
	ch := make(MsgChannel, 100)
	router.Join(cid, ch)
	assert.Equal(1, len(router.ConnMap))

	router.RegisterMethod(cid, "abc")
	methods := router.GetMethods(cid)
	assert.Equal(1, len(methods))
	assert.Equal("abc", methods[0])
}

func TestRouteMessage(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()

	cid := CID(1002)
	ch := make(MsgChannel, 100)
	router.Join(cid, ch)
	assert.Equal(1, len(router.ConnMap))
	router.RegisterMethod(cid, "abc")

	cid1 := CID(1003)
	ch1 := make(MsgChannel, 100)
	router.Join(cid1, ch1)

	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`

	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.Nil(err)
	assert.Equal(json.Number("100002"), msg.Id)
	router.RouteMessage(msg, cid1)

	rcvmsg := <-ch
	assert.Equal(msg.Id, rcvmsg.Id)
	assert.True(rcvmsg.IsRequest())
}

func TestRouteRoutine(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()
	ctx, cancel := context.WithCancel(context.Background())
	router.Start(ctx)
	defer cancel()

	cid := CID(1002)
	ch := make(MsgChannel, 100)

	router.ChJoin <- CmdJoin{RecvChannel: ch, ConnId: cid}
	router.ChRegister <- CmdRegister{ConnId: cid, Method: "abc"}

	cid1 := CID(1003)
	ch1 := make(MsgChannel, 100)
	router.ChJoin <- CmdJoin{RecvChannel: ch1, ConnId: cid1}

	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`

	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.Nil(err)
	assert.Equal(json.Number("100002"), msg.Id)

	router.ChMsg <- CmdMsg{Msg: msg, FromConnId: cid1}

	rcvmsg := <-ch
	assert.Equal(msg.Id, rcvmsg.Id)
	assert.True(rcvmsg.IsRequest())
}

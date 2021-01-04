package joint

import (
	//"fmt"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"testing"
	"time"
)

// implements ConnT
/*type TestConnT struct {
	ch MsgChannel
}

func (self TestConnT) RecvChannel() MsgChannel {
	return self.ch
}

func (self TestConnT) Location() MethodLocation {
	return Location_Local
}*/

func TestJoinConn(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter("test_join_conn")

	//cid := CID(1002)
	//ch := make(MsgChannel, 100)
	conn := router.Join() //cid, ch)
	assert.Equal(1, len(router.connMap))

	router.UpdateServeMethods(conn, []MethodInfo{{"abc", "method abc", "", nil}})
	methods := conn.GetMethods()
	assert.Equal(1, len(methods))
	assert.Equal("abc", methods[0])
}

func TestRouteMessage(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter("test_message")

	//cid := CID(1002)
	//ch := make(MsgChannel, 100)
	conn := router.Join()

	assert.Equal(1, len(router.connMap))
	router.UpdateServeMethods(conn, []MethodInfo{{"abc", "", "", nil}, {"def", "", "", nil}})

	methods := router.GetAllMethods()
	assert.Equal([]string{"abc", "def"}, methods)

	localMethods := router.GetLocalMethods()
	assert.Equal(2, len(localMethods))
	assert.Equal("abc", localMethods[0].Name)
	assert.Equal("def", localMethods[1].Name)

	_ = router.Join()

	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`

	msg, err := jsonrpc.ParseBytes([]byte(j1))
	assert.Nil(err)
	assert.Equal(json.Number("100002"), msg.MustId())
	router.RouteMessage(CmdMsg{MsgVec{msg, conn.ConnId}, false})

	rcvmsg := <-conn.RecvChannel
	assert.Equal(msg.MustId(), rcvmsg.Msg.MustId())
	assert.True(rcvmsg.Msg.IsRequest())
}

func TestRouteRoutine(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter("test_route_routine")
	ctx, cancel := context.WithCancel(context.Background())
	router.Start(ctx)
	defer cancel()

	time.Sleep(1 * time.Second)

	conn := router.Join()
	cid := conn.ConnId
	ch := conn.RecvChannel
	router.ChServe <- CmdServe{ConnId: cid, Methods: []MethodInfo{{"abc", "method abc", "", nil}}}

	conn1 := router.Join()
	cid1 := conn1.ConnId

	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`
	time.Sleep(1 * time.Second)
	msg, err := jsonrpc.ParseBytes([]byte(j1))
	assert.Nil(err)
	assert.Equal(json.Number("100002"), msg.MustId())

	router.ChMsg <- CmdMsg{MsgVec{msg, cid1}, false}

	//fmt.Printf("will rcv %v\n", ch)
	rcvmsg := <-ch
	//fmt.Printf("recved %v\n", rcvmsg)

	assert.Equal(msg.MustId(), rcvmsg.Msg.MustId())
	assert.True(rcvmsg.Msg.IsRequest())
}

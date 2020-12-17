package tube

import (
	//"fmt"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
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
	router := NewRouter()

	//cid := CID(1002)
	//ch := make(MsgChannel, 100)
	conn := router.Join() //cid, ch)
	assert.Equal(1, len(router.ConnMap))

	router.UpdateMethods(conn, []MethodInfo{{"abc", "method abc", false}})
	methods := conn.GetMethods()
	assert.Equal(1, len(methods))
	assert.Equal("abc", methods[0])
}

func TestRouteMessage(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()

	//cid := CID(1002)
	//ch := make(MsgChannel, 100)
	conn := router.Join()

	assert.Equal(1, len(router.ConnMap))
	router.UpdateMethods(conn, []MethodInfo{{"abc", "", false}, {"def", "", true}})

	methods := router.GetAllMethods()
	assert.Equal([]string{"abc", "def"}, methods)

	localMethods := router.GetLocalMethods()
	assert.Equal([]MethodInfo{{"abc", "", false}}, localMethods)

	_ = router.Join()

	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`

	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.Nil(err)
	assert.Equal(json.Number("100002"), msg.Id)
	router.RouteMessage(CmdMsg{msg, conn.ConnId, false})

	rcvmsg := <-conn.RecvChannel
	assert.Equal(msg.Id, rcvmsg.Id)
	assert.True(rcvmsg.IsRequest())
}

func TestRouteRoutine(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()
	ctx, cancel := context.WithCancel(context.Background())
	router.Start(ctx)
	defer cancel()

	time.Sleep(1 * time.Second)

	conn := router.Join()
	cid := conn.ConnId
	ch := conn.RecvChannel
	router.ChUpdate <- CmdUpdate{ConnId: cid, Methods: []MethodInfo{{"abc", "method abc", false}}}

	conn1 := router.Join()
	cid1 := conn1.ConnId

	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`
	time.Sleep(1 * time.Second)
	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.Nil(err)
	assert.Equal(json.Number("100002"), msg.Id)

	router.ChMsg <- CmdMsg{Msg: msg, FromConnId: cid1}

	//fmt.Printf("will rcv %v\n", ch)
	rcvmsg := <-ch
	//fmt.Printf("recved %v\n", rcvmsg)

	assert.Equal(msg.Id, rcvmsg.Id)
	assert.True(rcvmsg.IsRequest())
}

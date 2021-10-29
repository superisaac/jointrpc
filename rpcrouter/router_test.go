package rpcrouter

import (
	//"fmt"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestJoinConn(t *testing.T) {
	assert := assert.New(t)
	factory := NewRouterFactory("test_join_conn")

	router := factory.DefaultRouter()
	conn := router.Join() //cid, ch)
	assert.Equal(1, len(router.connMap))

	router.UpdateServeMethods(conn, []MethodInfo{{"abc", "method abc", "", nil}})
	methods := conn.GetMethods()
	assert.Equal(1, len(methods))
	assert.Equal("abc", methods[0])
}

func TestRouteMessage(t *testing.T) {
	assert := assert.New(t)
	router := NewRouterFactory("test_message").DefaultRouter()

	//cid := CID(1002)
	//ch := make(MsgChannel, 100)
	conn := router.Join()

	assert.Equal(1, len(router.connMap))
	router.UpdateServeMethods(conn, []MethodInfo{{"abc", "", "", nil}, {"def", "", "", nil}})

	methods := router.GetMethodNames()
	assert.Equal([]string{"abc", "def"}, methods)

	localMethods := router.GetMethods()
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
	assert.False(router.Started())

	router.deliverMessage(CmdMsg{
		MsgVec: MsgVec{
			Msg:        msg,
			Namespace:  router.Name(),
			FromConnId: conn.ConnId},
		ChRes: conn.RecvChannel})
	rcvmsg := <-conn.RecvChannel
	assert.True(rcvmsg.Msg.IsRequest())
	assert.Equal("abc", rcvmsg.Msg.MustMethod())
}

func TestRouteRoutine(t *testing.T) {
	assert := assert.New(t)
	factory := NewRouterFactory("test_route_routine")
	router := factory.DefaultRouter()
	factory.EnsureStart(context.Background())
	defer factory.Stop()

	time.Sleep(100 * time.Millisecond)

	conn := router.Join()
	cid := conn.ConnId
	ch := conn.RecvChannel
	router.ChMethods <- CmdMethods{
		Namespace: router.Name(),
		ConnId:    cid,
		Methods:   []MethodInfo{{"abc", "method abc", "", nil}},
	}

	conn1 := router.Join()
	cid1 := conn1.ConnId
	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`
	time.Sleep(100 * time.Millisecond)
	msg, err := jsonrpc.ParseBytes([]byte(j1))
	assert.Nil(err)
	assert.Equal(json.Number("100002"), msg.MustId())

	router.deliverRequest(MsgVec{
		Msg:        msg,
		Namespace:  router.Name(),
		FromConnId: cid1,
		ToConnId:   cid,
	}, 0, ch)
	rcvmsg := <-ch
	//assert.Equal(msg.MustId(), rcvmsg.Msg.MustId())
	assert.True(rcvmsg.Msg.IsRequest())
	assert.Equal("abc", rcvmsg.Msg.MustMethod())

	// wrong target id
	conn2 := router.Join()
	cid2 := conn2.ConnId
	j2 := `{
"id": 100003,
"method": "abc",
"params": [2, 6]
}`
	msg2, err := jsonrpc.ParseBytes([]byte(j2))
	assert.Nil(err)
	assert.Equal(json.Number("100003"), msg2.MustId())
	router.deliverRequest(MsgVec{
		Msg:        msg2,
		FromConnId: cid2,
		ToConnId:   CID(int(cid) + 100),
	}, 0, conn2.RecvChannel)

	rcvmsg2 := <-conn2.RecvChannel
	assert.Equal(msg2.MustId(), rcvmsg2.Msg.MustId())
	assert.True(rcvmsg2.Msg.IsError())

}

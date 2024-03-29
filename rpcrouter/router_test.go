package rpcrouter

import (
	//"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(router, conn.router)

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

	msg, err := jsonz.ParseBytes([]byte(j1))
	assert.Nil(err)
	assert.True(msg.IsRequest())
	assert.Equal(100002, msg.MustId())
	assert.False(router.Started())

	//router.deliverMessage(CmdMsg{
	chRes := make(MsgChannel, 1)
	//router.MsgInput() <- CmdMsg{
	router.relayMessage(CmdMsg{
		Msg:       msg,
		Namespace: router.Name(),
		ChRes:     chRes})
	cmdMsg := <-conn.MsgInput()
	assert.True(cmdMsg.Msg.IsRequest())
	assert.Equal("abc", cmdMsg.Msg.MustMethod())
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
	chRes := make(MsgChannel, 100)
	router.ChMethods <- CmdMethods{
		Namespace: router.Name(),
		ConnId:    cid,
		Methods:   []MethodInfo{{"abc", "method abc", "", nil}},
	}

	_ = router.Join()
	j1 := `{
"id": 100002,
"method": "abc",
"params": [1, 3]
}`
	time.Sleep(100 * time.Millisecond)
	msg, err := jsonz.ParseBytes([]byte(j1))
	assert.Nil(err)
	assert.Equal(100002, msg.MustId())

	conn.handleRequest(factory.startCtx, CmdMsg{
		Msg:       msg,
		Namespace: router.Name(),
		Timeout:   0,
		ChRes:     chRes,
		ConnId:    cid,
	})
	rcvmsg := <-conn.MsgOutput()
	//assert.Equal(msg.MustId(), rcvmsg.Msg.MustId())
	assert.True(rcvmsg.Msg.IsRequest())
	assert.Equal("abc", rcvmsg.Msg.MustMethod())

	// wrong target id
	_ = router.Join()
	j2 := `{
"id": 100003,
"method": "abc",
"params": [2, 6]
}`
	msg2, err := jsonz.ParseBytes([]byte(j2))
	assert.Nil(err)
	assert.Equal(100003, msg2.MustId())
	chRes2 := make(MsgChannel, 100)
	router.relayMessage(CmdMsg{
		Msg:     msg2,
		Timeout: 0,
		ChRes:   chRes2,
		ConnId:  CID(int(cid) + 100),
	})
	rcvmsg2 := <-chRes2
	assert.Equal(msg2.MustId(), rcvmsg2.Msg.MustId())
	assert.True(rcvmsg2.Msg.IsError())

}

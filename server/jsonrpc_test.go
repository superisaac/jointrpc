package server

import (
	//"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	//simplejson "github.com/bitly/go-simplejson"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

const INT_100 = jsonrpc.UID(100)

func TestReqConvert(t *testing.T) {
	assert := assert.New(t)

	j1 := `{
"version": "2.0",
"id": 100, 
"method": "testAgain",
"params": [3, "hello", "nice"]
}`
	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.Nil(err)
	//assert.Equal(json.Number("100"), msg.Id)
	assert.Equal(INT_100, msg.Id)

	req, err := MessageToRequest(msg)
	assert.Nil(err)

	assert.Equal(int64(100), req.Id)
	assert.Equal("testAgain", req.Method)
	assert.Equal(`[3,"hello","nice"]`, req.Params)

	msg1, err := RequestToMessage(req)
	assert.Nil(err)

	assert.True(msg1.IsRequest())
	assert.Equal(INT_100, msg1.Id)
}

func TestNotifyConvert(t *testing.T) {
	assert := assert.New(t)

	j1 := `{
"version": "2.0",
"method": "testAgain",
"params": [3, "hello", "nice"]
}`
	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.True(msg.IsNotify())
	assert.Nil(err)
	assert.Equal(jsonrpc.UID(0), msg.Id)

	notify, err := MessageToRequest(msg)
	assert.Nil(err)

	assert.Equal(int64(0), notify.Id)
	assert.Equal("testAgain", notify.Method)
	assert.Equal(`[3,"hello","nice"]`, notify.Params)

	msg1, err := RequestToMessage(notify)
	assert.Nil(err)

	assert.True(msg1.IsNotify())
	assert.Equal(jsonrpc.UID(0), msg1.Id)
}

func TestResultConvert(t *testing.T) {
	assert := assert.New(t)

	j1 := `{
"version": "2.0",
"id": 100, 
"result": "ok"
}`
	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.Nil(err)
	assert.Equal(INT_100, msg.Id)

	_, err = MessageToRequest(msg)
	assert.Equal("msg is neither request nor notify", err.Error())

	res, err := MessageToResult(msg)
	assert.Nil(err)

	assert.Equal(int64(100), res.Id)
	assert.Equal("\"ok\"", res.GetOk())

	msg1, err := ResultToMessage(res)
	assert.Nil(err)

	assert.True(msg1.IsResult())
	assert.Equal(INT_100, msg1.Id)
}

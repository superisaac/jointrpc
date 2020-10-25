package server

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	//simplejson "github.com/bitly/go-simplejson"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func TestConvertBetween(t *testing.T) {
	assert := assert.New(t)

	j1 := `{
"id": 100, 
"method": "testAgain",
"params": [3, "hello", "nice"]
}`
	msg, err := jsonrpc.ParseMessage([]byte(j1))
	assert.Nil(err)
	assert.Equal(msg.Id, json.Number("100"))
}

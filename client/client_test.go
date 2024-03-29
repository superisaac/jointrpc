package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"net/url"
	"testing"
)

func TestServerUrl(t *testing.T) {
	assert := assert.New(t)
	ustr := "h2c://localhost:8999#cert=/tmp/p1.cert&ff=kk"
	u, err := url.Parse(ustr)
	assert.Nil(err)

	assert.Equal("cert=/tmp/p1.cert&ff=kk", u.Fragment)

	v, err := url.ParseQuery(u.Fragment)
	assert.Nil(err)
	assert.Equal("", v.Get("opp"))
	assert.Equal("/tmp/p1.cert", v.Get("cert"))
}

func TestClientProps(t *testing.T) {
	assert := assert.New(t)

	c := NewRPCClient(ServerEntry{ServerUrl: "https://127.0.0.1:8899"})

	assert.Equal("wss://127.0.0.1:8899/ws", c.WebsocketUrlString())

	assert.True(c.IsSecure())
	assert.True(c.IsHttp())
}

func TestMapstructure(t *testing.T) {
	assert := assert.New(t)

	minfo := rpcrouter.MethodInfo{Name: "aaa",
		Help:       "a help text",
		SchemaJson: "{}"}

	infoDict := make(map[string]interface{})
	err := misc.DecodeStruct(minfo, &infoDict)
	assert.Nil(err)
	assert.Equal(minfo.Name, infoDict["name"])
	assert.Equal(minfo.Help, infoDict["help"])
	assert.Equal(minfo.SchemaJson, infoDict["schema"])

	var newminfo rpcrouter.MethodInfo
	err = misc.DecodeStruct(infoDict, &newminfo)
	assert.Nil(err)
	assert.Equal(newminfo.Name, infoDict["name"])
	assert.Equal(newminfo.Help, infoDict["help"])
	assert.Equal(newminfo.SchemaJson, infoDict["schema"])

}

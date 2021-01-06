package client

import (
	"github.com/stretchr/testify/assert"
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

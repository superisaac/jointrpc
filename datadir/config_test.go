package datadir

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestBasicAuth(t *testing.T) {
	assert := assert.New(t)

	bauth := BasicAuth{Username: "aaa", Password: "1111", AllowedSources: []string{"192.168.0.0/16", "172.18.1.0/24"}}
	err := bauth.validateValues()
	assert.Nil(err)
	assert.Equal(2, len(bauth.allowedIPNets))
	assert.Equal("192.168.0.0/16", bauth.allowedIPNets[0].String())
	assert.Equal("172.18.1.0/24", bauth.allowedIPNets[1].String())
	assert.True(bauth.Authorize("aaa", "1111", "192.168.2.3"))

	bauth = BasicAuth{Username: "aaa", Password: "1111", AllowedSources: []string{"192.168.1.3/32"}}
	err = bauth.validateValues()
	assert.Nil(err)
	assert.Equal(1, len(bauth.allowedIPNets))
	assert.Equal("192.168.1.3/32", bauth.allowedIPNets[0].String())
	assert.False(bauth.Authorize("aaa", "1111", "192.168.2.3"))

	var a []string = nil
	assert.Equal(0, len(a))
}

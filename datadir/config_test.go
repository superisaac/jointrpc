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

	bauth := BasicAuth{Username: "aaa", Password: "1111", AllowedCIDR: "192.168.0.0/16"}
	err := bauth.validateValues()
	assert.Nil(err)
	assert.Equal("192.168.0.0", bauth.cidrIP.String())
	assert.True(bauth.Authorize("aaa", "1111", "192.168.2.3"))

	bauth = BasicAuth{Username: "aaa", Password: "1111", AllowedCIDR: "192.168.1.3/32"}
	err = bauth.validateValues()
	assert.Nil(err)
	assert.Equal("192.168.1.3", bauth.cidrIP.String())
	assert.Equal("192.168.1.3/32", bauth.cidrIPNet.String())
	assert.False(bauth.Authorize("aaa", "1111", "192.168.2.3"))
}

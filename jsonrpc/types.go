package jsonrpc

import (
	"github.com/bitly/go-simplejson"
)

//type CID uint64
type UID uint64

type RPCMessage struct {
	Initialized bool
	//FromConnId  CID
	Id          interface{}
	Method      string
	Params      *simplejson.Json
	Result      *simplejson.Json
	Error       *simplejson.Json
	Raw         *simplejson.Json
}

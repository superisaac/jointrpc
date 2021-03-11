package jsonrpc

import (
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
)

//type CID uint64
//type UID uint64
type UID string

/*type RPCMessage struct {
	Initialized bool
	//FromConnId  CID
	Id interface{}
	//Id     UID
	Method string
	Params *simplejson.Json
	Result *simplejson.Json
	Error  *simplejson.Json
	Raw    *simplejson.Json
} */

type RPCError struct {
	Code      int
	Reason    string
	Retryable bool
}

const (
	MTRequest = 1
	MTNotify  = 2
	MTResult  = 3
	MTError   = 4
)

type IMessage interface {
	IsRequest() bool
	IsNotify() bool
	IsRequestOrNotify() bool
	IsResult() bool
	IsError() bool
	IsResultOrError() bool

	EncodePretty() (string, error)
	Interface() interface{}
	MustString() string
	Bytes() ([]byte, error)

	// TraceId
	SetTraceId(traceId string)
	TraceId() string

	// upvote
	MustId() interface{}
	MustMethod() string
	MustParams() []interface{}
	MustResult() interface{}
	MustError() interface{}

	Log() *log.Entry
}

type BaseMessage struct {
	messageType int
	raw         *simplejson.Json
	traceId     string
}

type RequestMessage struct {
	BaseMessage
	Id     interface{}
	Method string
	Params []interface{}
}

type NotifyMessage struct {
	BaseMessage
	Method string
	Params []interface{}
}

type ResultMessage struct {
	BaseMessage
	Id     interface{}
	Result interface{}
}

type ErrorMessage struct {
	BaseMessage
	Id    interface{}
	Error interface{}
}

package jsonrpc

import (
	"github.com/bitly/go-simplejson"
)

//type CID uint64
//type UID uint64
type UID string

type RPCMessage struct {
	Initialized bool
	//FromConnId  CID
	Id interface{}
	//Id     UID
	Method string
	Params *simplejson.Json
	Result *simplejson.Json
	Error  *simplejson.Json
	Raw    *simplejson.Json
}

type RPCError struct {
	Code      int
	Reason    string
	Retryable bool
}

// Schema builder
type SchemaBuildError struct {
	info string
}

type SchemaBuilder struct {
}

// Schema validator
type SchemaValidator struct {
	paths     []string
	hint      string
	errorPath string
}

type ErrorPos struct {
	paths []string
	hint  string
}

type Schema interface {
	// returns the generated
	Type() string
	Scan(validator *SchemaValidator, data interface{}) *ErrorPos
}

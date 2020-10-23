package tube

import (
	//"github.com/gorilla/websocket"
	"errors"
	"sync"
	"time"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"	
)


// 5 seconds
const (
	DefaultRequestTimeout time.Duration = 1000000 * 5
	
	IntentLocal string = "local"
)

var (
	ErrNotNotify = errors.New("json message is not notify")
)

// Commands
type MsgChannel chan *jsonrpc.RPCMessage

// Pending Struct
type PendingKey struct {
	ConnId jsonrpc.CID
	MsgId  interface{}
}

type PendingValue struct {
	ConnId jsonrpc.CID
	Expire time.Time
}

// IConn
type IConn interface {
	RecvChannel() MsgChannel
	CanBroadcast() bool
}

type Router struct {
	// channels
	routerLock    *sync.RWMutex
	MethodConnMap map[string]([]jsonrpc.CID)
	ConnMethodMap map[jsonrpc.CID]([]string)

	ConnMap    map[jsonrpc.CID](IConn)
	PendingMap map[PendingKey]PendingValue
}

// An ConnActor manage a websocket connection and handles incoming messages

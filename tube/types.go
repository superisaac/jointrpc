package tube

import (
	//"github.com/gorilla/websocket"
	"errors"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	"sync"
	"time"
)

// 5 seconds
const (
	DefaultRequestTimeout time.Duration = 1000000 * 5

	IntentLocal string = "local"
)

var (
	ErrNotNotify = errors.New("json message is not notify")
)

type CID uint64

// Commands
type MsgChannel chan *jsonrpc.RPCMessage

// Pending Struct
type PendingKey struct {
	ConnId CID
	MsgId  interface{}
}

type PendingValue struct {
	ConnId CID
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
	MethodConnMap map[string]([]CID)
	ConnMethodMap map[CID]([]string)

	ConnMap    map[CID](IConn)
	PendingMap map[PendingKey]PendingValue
}

// An ConnActor manage a websocket connection and handles incoming messages

package rpcrouter

import (
	//"github.com/gorilla/websocket"
	"errors"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"net"
	"sync"
	"time"
)

// 5 seconds
const (
	DefaultRequestTimeout time.Duration = 1000000 * 5

	//IntentLocal string = "local"
)

var (
	ErrNotNotify             = errors.New("json message is not notify")
	ErrRequestNotifyRequired = errors.New("only request and notify message accepted")
)

type CID uint64

const ZeroCID = CID(0)

// Commands
type MsgVec struct {
	Msg          jsonrpc.IMessage
	TraceId      string
	FromConnId   CID
	TargetConnId CID
}
type MsgChannel chan MsgVec

// Pending Struct
type PendingKey struct {
	ConnId CID
	MsgId  interface{}
}

type PendingValue struct {
	ConnId CID
	Expire time.Time
}

type MethodInfo struct {
	Name       string
	Help       string
	SchemaJson string
	schemaObj  schema.Schema
}

// tube state
type TubeState struct {
	Methods []MethodInfo
}

// Connect Struct
type ConnT struct {
	ConnId      CID
	PeerAddr    net.Addr
	RecvChannel MsgChannel

	ServeMethods    map[string]MethodInfo
	DelegateMethods map[string]bool

	AsFallback bool
	watchState bool

	stateChannel chan *TubeState
}

type MethodDesc struct {
	Conn *ConnT
	Info MethodInfo
}

type MethodDelegation struct {
	Conn *ConnT
	Name string // method name
}

// Method update watcher
type MethodInfoMap map[string](interface{})

// Channel commands
type CmdMsg struct {
	MsgVec    MsgVec
	Broadcast bool
}

// type CmdJoin struct {
// 	ConnId      CID
// 	RecvChannel MsgChannel
// }

// type CmdLeave struct {
// 	ConnId CID
// }

type CmdServe struct {
	ConnId  CID
	Methods []MethodInfo
}

type CmdDelegate struct {
	ConnId      CID
	MethodNames []string
}

type Router struct {
	// channels
	name            string
	routerLock      *sync.RWMutex
	methodConnMap   map[string]([]MethodDesc)
	delegateConnMap map[string]([]MethodDelegation)

	fallbackConns []*ConnT

	connMap    map[CID](*ConnT)
	pendingMap map[PendingKey]PendingValue

	// channels
	ChMsg chan CmdMsg
	//ChJoin     chan CmdJoin
	//ChLeave  chan CmdLeave
	ChServe    chan CmdServe
	ChDelegate chan CmdDelegate

	methodsSig string
}

// type TubeT struct {
// 	Router *Router
// }

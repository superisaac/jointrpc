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

const (
	// 5 seconds
	DefaultRequestTimeout time.Duration = time.Second * 5
)

var (
	ErrNotNotify             = errors.New("json message is not notify")
	ErrRequestNotifyRequired = errors.New("only request and notify message accepted")
)

type CID uint64

const ZeroCID = CID(0)

// Commands
type MsgVec struct {
	Msg        jsonrpc.IMessage
	FromConnId CID
	ToConnId   CID
}
type MsgChannel chan MsgVec

// Pending Struct
type PendingT struct {
	//OrigMsgVec MsgVec
	ReqMsg     *jsonrpc.RequestMessage
	FromConnId CID
	ToConnId   CID
	Expire     time.Time
}

type MethodInfo struct {
	Name       string
	Help       string
	SchemaJson string
	schemaObj  *schema.MethodSchema
}

// tube state
type ServerState struct {
	Methods []MethodInfo
}

// Connect Struct
type ConnT struct {
	ConnId      CID
	publicId    string
	PeerAddr    net.Addr
	RecvChannel MsgChannel

	ServeMethods    map[string]MethodInfo
	DelegateMethods map[string]bool

	AsFallback bool
	watchState bool

	stateChannel chan *ServerState
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
	MsgVec  MsgVec
	Timeout time.Duration
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

	connMap       map[CID](*ConnT)
	publicConnMap map[string](*ConnT)

	pendingRequests map[interface{}]PendingT

	// channels
	chMsg      chan CmdMsg
	ChServe    chan CmdServe
	ChDelegate chan CmdDelegate

	methodsSig string

	// flags
	ValidateSchema bool
}

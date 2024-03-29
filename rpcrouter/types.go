package rpcrouter

import (
	"context"
	"github.com/pkg/errors"
	datadir "github.com/superisaac/jointrpc/datadir"
	"github.com/superisaac/jsonz"
	schema "github.com/superisaac/jsonz/schema"
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
// type MsgVec struct {
// 	Msg       jsonz.Message
// 	Namespace string
// }

type MethodInfo struct {
	Name       string `json:"name"`
	Help       string `json:"help"`
	SchemaJson string `json:"schema"`
	schemaObj  *schema.MethodSchema
}

// tube state
type ServerState struct {
	Methods []MethodInfo `json:"methods"`
}

// Connect Struct
type ConnPending struct {
	cmdMsg CmdMsg
	Expire time.Time
}

type ConnT struct {
	ConnId    CID
	Namespace string
	PeerAddr  net.Addr

	router *Router

	ServeMethods    map[string]MethodInfo
	DelegateMethods map[string]bool

	watchState bool

	stateChannel chan *ServerState

	lastPing time.Time

	msgOutput MsgChannel
	msgInput  MsgChannel //chan CmdMsg
	pendings  map[interface{}]ConnPending
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
type CmdRet struct {
	Ok bool
}

type CmdJoin struct {
	Conn  *ConnT
	ChRet chan CmdRet
}

type CmdLeave struct {
	Conn  *ConnT
	ChRet chan CmdRet
}

type RetSelectConn struct {
	Method string
	ConnId CID
	Conn   *ConnT
	Found  bool
}
type CmdSelectConn struct {
	Method string
	ConnId CID
	ChRet  chan RetSelectConn
}

type CmdMsg struct {
	Msg       jsonz.Message
	Namespace string
	ConnId    CID
	Timeout   time.Duration
	ChRes     MsgChannel
}
type MsgChannel chan CmdMsg

type CmdMethods struct {
	Namespace string
	ConnId    CID
	Methods   []MethodInfo
}

type CmdDelegates struct {
	Namespace   string
	ConnId      CID
	MethodNames []string
}

type Router struct {
	name    string
	factory *RouterFactory
	//routerLock      *sync.RWMutex
	//pendingLock     *sync.RWMutex
	methodConnMap   map[string]([]MethodDesc)
	methodsSig      string
	connMap         map[CID](*ConnT)
	delegateConnMap map[string]([]MethodDelegation)
	//pendingRequests map[interface{}]PendingT
	//started         bool
	startCtx   context.Context
	cancelFunc func()

	// channels
	ChJoin       chan CmdJoin
	ChLeave      chan CmdLeave
	chRouteMsg   MsgChannel
	chSelectConn chan CmdSelectConn
	ChMethods    chan CmdMethods
	ChDelegates  chan CmdDelegates
}

type RouterFactory struct {
	// channels
	name string

	routerMap sync.Map

	// flags
	Config *datadir.Config
	//started        bool
	startCtx   context.Context
	cancelFunc func()
}

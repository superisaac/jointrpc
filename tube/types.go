package tube

import (
	//"github.com/gorilla/websocket"
	"errors"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	schema "github.com/superisaac/rpctube/jsonrpc/schema"
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

// Commands
type MsgVec struct {
	Msg        jsonrpc.IMessage
	FromConnId CID
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
	Name      string
	Help      string
	Schema    schema.Schema
	Delegated bool
}

// Connect Struct
type ConnT struct {
	ConnId      CID
	PeerAddr    net.Addr
	RecvChannel MsgChannel
	Methods     map[string]MethodInfo
	AsFallback  bool
}

type MethodDesc struct {
	//ConnId  CID
	Conn *ConnT
	Info MethodInfo
}

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

type CmdUpdate struct {
	ConnId  CID
	Methods []MethodInfo
}

type Router struct {
	// channels
	routerLock    *sync.RWMutex
	methodConnMap map[string]([]MethodDesc)
	fallbackConns []*ConnT

	connMap    map[CID](*ConnT)
	pendingMap map[PendingKey]PendingValue

	// channels
	ChMsg chan CmdMsg
	//ChJoin     chan CmdJoin
	//ChLeave  chan CmdLeave
	ChUpdate chan CmdUpdate

	localMethodsSig string
}

type TubeT struct {
	Router *Router
}

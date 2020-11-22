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

	//IntentLocal string = "local"
)

var (
	ErrNotNotify = errors.New("json message is not notify")
)

type CID uint64

type MethodLocation int32

const (
	Location_Local  = 0
	Location_Remote = 1
)

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

// Connect Struct
type ConnT struct {
	ConnId      CID
	RecvChannel MsgChannel
	Methods     map[string]bool
}

type MethodDesc struct {
	//ConnId  CID
	Conn     *ConnT
	Location MethodLocation
}

// Channel commands
type CmdMsg struct {
	Msg        *jsonrpc.RPCMessage
	FromConnId CID
}

type CmdJoin struct {
	ConnId      CID
	RecvChannel MsgChannel
}

type CmdLeave struct {
	ConnId CID
}

type CmdReg struct {
	ConnId   CID
	Method   string
	Location MethodLocation
}

type CmdUnreg struct {
	ConnId CID
	Method string
}

type Router struct {
	// channels
	routerLock    *sync.RWMutex
	MethodConnMap map[string]([]MethodDesc)
	//ConnMethodMap map[CID]([]string)

	ConnMap    map[CID](*ConnT)
	PendingMap map[PendingKey]PendingValue

	// channels
	ChMsg chan CmdMsg
	//ChJoin     chan CmdJoin
	ChLeave chan CmdLeave
	ChReg   chan CmdReg
	ChUnreg chan CmdUnreg
}

type TubeT struct {
	Router *Router
}

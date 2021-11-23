package builtin

import (
	"context"
	//"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/misc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonrpc"
)

type BuiltinService struct {
	disp     *dispatch.Dispatcher
	chResult chan dispatch.ResultT
	//router *rpcrouter.Router
	conn *rpcrouter.ConnT
	done chan error
}

func (self BuiltinService) Name() string {
	return "builtin"
}

func (self BuiltinService) CanRun(rootCtx context.Context) bool {
	return true
}

func (self *BuiltinService) Start(rootCtx context.Context) error {
	self.Init(rootCtx)

	factory := rpcrouter.RouterFactoryFromContext(rootCtx)
	commonRouter := factory.CommonRouter()

	ctx, cancel := context.WithCancel(rootCtx)
	defer func() {
		cancel()
		log.Debug("buildin dispatcher context canceled")
	}()

	self.conn = commonRouter.Join()
	// self.conn = rpcrouter.NewConn()
	// commonRouter.ChJoin <- rpcrouter.CmdJoin{Conn: self.conn}

	defer func() {
		log.Debugf("conn %d leave router", self.conn.ConnId)
		//commonRouter.Leave(self.conn)
		commonRouter.ChLeave <- rpcrouter.CmdLeave{Conn: self.conn}
		self.conn = nil
	}()

	self.declareMethods(factory)

	senderCtx, cancelSender := context.WithCancel(ctx)
	defer cancelSender()
	dispatch.SenderLoop(senderCtx, self, self.conn, self.chResult)
	return nil
}

func (self *BuiltinService) requestReceived(ctx context.Context, cmdMsg rpcrouter.CmdMsg) {
	msg := cmdMsg.Msg
	if msg.IsRequestOrNotify() {
		self.disp.Feed(ctx, cmdMsg, self.chResult)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
}

const (
	echoSchema = `{
"type": "method",
 "params": ["string"],
 "returns": "string"
}`
)

func (self *BuiltinService) Init(rootCtx context.Context) *BuiltinService {
	factory := rpcrouter.RouterFactoryFromContext(rootCtx)

	self.disp = dispatch.NewDispatcher()
	self.chResult = make(chan dispatch.ResultT, misc.DefaultChanSize())
	self.done = make(chan error, 10)

	self.disp.OnTyped("_listMethods", func(req *dispatch.RPCRequest) ([]rpcrouter.MethodInfo, error) {
		router := factory.Get(req.CmdMsg.Namespace)
		minfos := router.GetMethods()
		return minfos, nil
	})

	self.disp.OnTyped("_echo", func(req *dispatch.RPCRequest, text string) (string, error) {
		return text, nil
	}, dispatch.WithSchema(echoSchema))

	self.disp.OnChange(func() {
		self.declareMethods(factory)
	})
	return self
}

func (self *BuiltinService) declareMethods(factory *rpcrouter.RouterFactory) {
	if self.conn != nil {
		minfos := self.disp.GetMethodInfos()
		ns := factory.CommonRouter().Name()
		cmdMethods := rpcrouter.CmdMethods{
			Namespace: ns,
			ConnId:    self.conn.ConnId,
			Methods:   minfos,
		}
		factory.Get(cmdMethods.Namespace).ChMethods <- cmdMethods
	}
}

func NewBuiltinService() *BuiltinService {
	return new(BuiltinService)
}

func (self BuiltinService) SendMessage(ctx context.Context, msg jsonrpc.IMessage) error {
	factory := rpcrouter.RouterFactoryFromContext(ctx)
	commonRouter := factory.CommonRouter()
	self.conn.MsgInput() <- rpcrouter.CmdMsg{
		Msg:       msg,
		Namespace: commonRouter.Name(),
	}
	return nil
}

func (self BuiltinService) SendCmdMsg(ctx context.Context, cmdMsg rpcrouter.CmdMsg) error {
	msg := cmdMsg.Msg
	if msg.IsRequestOrNotify() {
		self.disp.Feed(ctx, cmdMsg, self.chResult)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
	return nil
}

func (self BuiltinService) Done() chan error {
	return self.done
}

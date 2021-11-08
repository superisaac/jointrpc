package builtin

import (
	"context"
	//"fmt"
	log "github.com/sirupsen/logrus"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"time"

	//schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/dispatch"
	//service "github.com/superisaac/jointrpc/service"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
)

type BuiltinService struct {
	disp     *dispatch.Dispatcher
	chResult chan dispatch.ResultT
	//router *rpcrouter.Router
	conn *rpcrouter.ConnT
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

	for {
		select {
		case <-ctx.Done():
			log.Debugf("builtin handlers, context done")
			return nil
		case <-time.After(3 * time.Second):
			self.conn.ClearPendings()
		case cmdMsg, ok := <-self.conn.MsgOutput():
			if !ok {
				log.Debugf("recv channel colosed, leave")
				return nil
			}
			//timeoutCtx, _ := context.WithTimeout(rootCtx, 10 * time.Second)
			self.requestReceived(ctx, cmdMsg)
		case cmdMsg, ok := <-self.conn.MsgInput():
			if !ok {
				log.Debugf("MsgInput() closed")
				return nil
			}
			err := self.conn.HandleRouteMessage(ctx, cmdMsg)
			if err != nil {
				panic(err)
			}
		case result, ok := <-self.chResult:
			if !ok {
				log.Infof("result channel closed, return")
				return nil
			}

			self.conn.MsgInput() <- rpcrouter.CmdMsg{
				Msg:       result.ResMsg,
				Namespace: commonRouter.Name(),
			}
		}
	}
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

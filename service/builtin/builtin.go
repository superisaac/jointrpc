package builtin

import (
	"context"
	"fmt"
	//"time"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"

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

	defer func() {
		log.Debugf("conn %d leave router", self.conn.ConnId)
		commonRouter.Leave(self.conn)
		self.conn = nil
	}()

	self.declareMethods(factory)

	for {
		select {
		case <-ctx.Done():
			log.Debugf("builtin handlers, context done")
			return nil
		case msgvec, ok := <-self.conn.RecvChannel:
			if !ok {
				log.Debugf("recv channel colosed, leave")
				return nil
			}
			//timeoutCtx, _ := context.WithTimeout(rootCtx, 10 * time.Second)
			self.requestReceived(ctx, msgvec)
		case result, ok := <-self.chResult:
			if !ok {
				log.Infof("result channel closed, return")
				return nil
			}
			commonRouter.DeliverResultOrError(
				rpcrouter.MsgVec{
					Msg:        result.ResMsg,
					Namespace:  commonRouter.Name(),
					FromConnId: self.conn.ConnId,
				})
		}
	}
	return nil
}

func (self *BuiltinService) requestReceived(ctx context.Context, msgvec rpcrouter.MsgVec) {
	msg := msgvec.Msg
	if msg.IsRequest() || msg.IsNotify() {
		self.disp.Feed(ctx, msgvec, self.chResult)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
}

const (
	echoSchema = `{
"type": "method",
 "params": [{"type": "string"}]
}`
)

func (self *BuiltinService) Init(rootCtx context.Context) *BuiltinService {
	factory := rpcrouter.RouterFactoryFromContext(rootCtx)

	self.disp = dispatch.NewDispatcher()
	self.chResult = make(chan dispatch.ResultT, misc.DefaultChanSize())

	self.disp.On("_listMethods", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		fmt.Printf("list methods, %s\n", req.MsgVec.Namespace)
		router := factory.Get(req.MsgVec.Namespace)
		minfos := router.GetMethods()

		minfos = append(minfos, factory.CommonRouter().GetMethods()...)
		arr := make([](rpcrouter.MethodInfoMap), 0)
		for _, minfo := range minfos {
			arr = append(arr, minfo.ToMap())
		}
		return arr, nil
	})

	self.disp.On("_echo", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		if len(params) < 1 {
			return nil, &jsonrpc.RPCError{Code: 400, Message: "len params should be at least 1"}
		}
		msg, ok := params[0].(string)
		if !ok {
			return nil, &jsonrpc.RPCError{Code: 400, Message: "string params required"}
		}
		return map[string]string{"echo": msg}, nil
	}, dispatch.WithSchema(echoSchema))

	self.disp.OnChange(func() {
		self.declareMethods(factory)
	})
	return self
}

func (self *BuiltinService) declareMethods(factory *rpcrouter.RouterFactory) {
	if self.conn != nil {
		minfos := self.disp.GetMethodInfos()
		cmdMethods := rpcrouter.CmdMethods{
			Namespace: factory.CommonRouter().Name(),
			ConnId:    self.conn.ConnId,
			Methods:   minfos,
		}
		factory.ChMethods <- cmdMethods
	}
}

func NewBuiltinService() *BuiltinService {
	return new(BuiltinService)
}

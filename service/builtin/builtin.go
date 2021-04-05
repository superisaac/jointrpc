package builtin

import (
	"context"
	//"fmt"
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
	disp   *dispatch.Dispatcher
	router *rpcrouter.Router
	conn   *rpcrouter.ConnT
}

func (self BuiltinService) Name() string {
	return "builtin"
}

func (self BuiltinService) CanRun(rootCtx context.Context) bool {
	return true
}

func (self *BuiltinService) Start(rootCtx context.Context) error {
	self.router = rpcrouter.RouterFromContext(rootCtx)
	ctx, cancel := context.WithCancel(rootCtx)
	defer func() {
		cancel()
		log.Debug("buildin dispatcher context canceled")
	}()

	self.conn = self.router.Join(false)

	defer func() {
		log.Debugf("conn %d leave router", self.conn.ConnId)
		self.router.Leave(self.conn)
		self.conn = nil
	}()

	self.declareMethods()

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
			self.requestReceived(msgvec)
		case resmsg, ok := <-self.disp.ChResult:
			if !ok {
				log.Infof("result channel closed, return")
				return nil
			}
			self.router.DeliverResultOrError(
				rpcrouter.MsgVec{
					Msg:        resmsg,
					FromConnId: self.conn.ConnId,
				})
		}
	}
	return nil
}

func (self *BuiltinService) requestReceived(msgvec rpcrouter.MsgVec) {
	msg := msgvec.Msg
	if msg.IsRequest() || msg.IsNotify() {
		self.disp.HandleRequestMessage(msgvec)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
}

const (
	echoSchema = `{"type": "method", "params": [{"type": "string"}]}`
)

func (self *BuiltinService) Init() *BuiltinService {
	misc.Assert(self.router == nil, "already initited")
	self.disp = dispatch.NewDispatcher()

	self.disp.On("_listMethods", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		minfos := self.router.GetMethods()

		arr := make([](rpcrouter.MethodInfoMap), 0)
		for _, minfo := range minfos {
			arr = append(arr, minfo.ToMap())
		}
		return arr, nil
	})

	self.disp.On("_echo", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		if len(params) < 1 {
			return nil, &jsonrpc.RPCError{Code: 400, Reason: "len params should be at least 1"}
		}
		msg, ok := params[0].(string)
		if !ok {
			return nil, &jsonrpc.RPCError{Code: 400, Reason: "string params required"}
		}
		return map[string]string{"echo": msg}, nil
	}, dispatch.WithSchema(echoSchema))

	self.disp.OnChange(func() {
		self.declareMethods()
	})
	return self
}

func (self *BuiltinService) declareMethods() {
	if self.conn != nil {
		minfos := make([]rpcrouter.MethodInfo, 0)
		for m, info := range self.disp.MethodHandlers {
			minfo := rpcrouter.MethodInfo{
				Name:       m,
				Help:       info.Help,
				SchemaJson: info.SchemaJson,
			}
			minfos = append(minfos, minfo)
		}
		cmdServe := rpcrouter.CmdServe{ConnId: self.conn.ConnId, Methods: minfos}
		self.router.ChServe <- cmdServe
	}
}

func NewBuiltinService() *BuiltinService {
	return new(BuiltinService).Init()
}

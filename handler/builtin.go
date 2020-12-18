package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
)

type BuiltinHandlerManager struct {
	HandlerManager
	conn *tube.ConnT
}

func (self *BuiltinHandlerManager) Start(ctx context.Context) {
	go func() {
		router := tube.Tube().Router
		self.conn = router.Join()

		defer func() {
			router.Leave(self.conn)
			self.conn = nil
		}()

		self.updateMethods()

		for {
			select {
			case <-ctx.Done():
				log.Debugf("builtin handlers, context done")
				return
			case msg, ok := <-self.conn.RecvChannel:
				if !ok {
					return
				}
				self.messageReceived(msg)
			case resmsg, ok := <-self.ChResultMsg:
				if !ok {
					return
				}
				router.ChMsg <- tube.CmdMsg{Msg: resmsg, FromConnId: self.conn.ConnId}
			}
		}
	}()
}

func (self *BuiltinHandlerManager) messageReceived(msg *jsonrpc.RPCMessage) {
	if msg.IsRequest() || msg.IsNotify() {
		self.HandleRequestMessage(msg)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
}

func (self *BuiltinHandlerManager) Init() *BuiltinHandlerManager {
	self.InitHandlerManager()
	self.On(".listMethods", func(req *RPCRequest, params []interface{}) (interface{}, error) {
		minfos := tube.Tube().Router.GetLocalMethods()
		arr := make([](tube.MethodInfoMap), 0)
		for _, minfo := range minfos {
			arr = append(arr, minfo.ToMap())
		}
		return arr, nil
	})

	self.On(".broadcast", func(req *RPCRequest, params []interface{}) (interface{}, error) {
		if len(params) < 1 {
			return nil, &jsonrpc.RPCError{Code: 400, Reason: "method must be specified"}
		}

		method, err := jsonrpc.ValidateString(params[0], "method")
		if err != nil {
			return nil, err
		}
		notify := jsonrpc.NewNotifyMessage(method, params[1:])

		connId := tube.CID(0)
		if self.conn != nil {
			connId = self.conn.ConnId
		}
		tube.Tube().Router.ChMsg <- tube.CmdMsg{Msg: notify, FromConnId: connId, Broadcast: true}
		return "ok", nil
	}, WithHelp("broadcast a notify to all receivers"))

	self.OnChange(func() {
		self.updateMethods()
	})
	return self
}

func (self *BuiltinHandlerManager) updateMethods() {
	if self.conn != nil {
		minfos := make([]tube.MethodInfo, 0)
		for m, info := range self.MethodHandlers {
			minfo := tube.MethodInfo{
				Name:      m,
				Help:      info.Help,
				Delegated: false,
			}
			minfos = append(minfos, minfo)
		}
		cmdUpdate := tube.CmdUpdate{ConnId: self.conn.ConnId, Methods: minfos}
		tube.Tube().Router.ChUpdate <- cmdUpdate
	}
}

var (
	builtin *BuiltinHandlerManager
)

func Builtin() *BuiltinHandlerManager {
	if builtin == nil {
		builtin = new(BuiltinHandlerManager).Init()
	}
	return builtin
}

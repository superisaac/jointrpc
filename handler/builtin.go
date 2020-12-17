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
					log.Debugf("builtin handlers, recv channel closed")
					return
				}
				self.messageReceived(msg)
			case resmsg, ok := <-self.ChResultMsg:
				if !ok {
					log.Debugf("builtin handlers, recv channel closed")
					return
				}
				log.Debugf("ch result msg %v", resmsg)
				router.ChMsg <- tube.CmdMsg{Msg: resmsg, FromConnId: self.conn.ConnId}
			}
		}
	}()
}

func (self *BuiltinHandlerManager) messageReceived(msg *jsonrpc.RPCMessage) {
	log.Debugf("message received %v", msg)
	if msg.IsRequest() || msg.IsNotify() {
		self.HandleRequestMessage(msg)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
}

func (self *BuiltinHandlerManager) Init() *BuiltinHandlerManager {
	self.InitHandlerManager()

	self.On(".listMethods", func(req *RPCRequest, params []interface{}) (interface{}, error) {
		log.Debugf("list methods %v", params)
		minfos := tube.Tube().Router.GetLocalMethods()
		arr := make([](tube.MethodInfoMap), 0)
		for _, minfo := range minfos {
			arr = append(arr, minfo.ToMap())
		}
		log.Debugf("got lcoal methods %v", arr)
		return arr, nil
	})

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

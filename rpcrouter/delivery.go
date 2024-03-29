package rpcrouter

import (
	"fmt"
	//log "github.com/sirupsen/logrus"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jsonz"
	//"time"
)

func (self CmdMsg) Res(res jsonz.Message) CmdMsg {
	return CmdMsg{Msg: res, Namespace: self.Namespace}
}

func (self *Router) PostMessage(cmdMsg CmdMsg) {
	// there are two paradigms two post methods
	// 1. send to router and let router relay the message to correspondend connections
	// 2. query the router goroutine for connections and send message to conn's msg output port, this way is a little bit faster but suffers the channel closing panic

	//self.chRouteMsg <- cmdMsg
	self.redirectMessage(cmdMsg)
}

func (self *Router) relayMessage(cmdMsg CmdMsg) {
	msg := cmdMsg.Msg
	misc.Assert(msg.IsRequestOrNotify(), "router only support request and notify")
	toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.ConnId)
	if found {
		toConn.MsgInput() <- cmdMsg
	} else if msg.IsRequest() {
		reqMsg, _ := msg.(*jsonz.RequestMessage)
		errMsg := jsonz.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		cmdMsg.ChRes <- cmdMsg.Res(errMsg)
	}
}

func (self *Router) redirectMessage(cmdMsg CmdMsg) {
	msg := cmdMsg.Msg
	misc.Assert(msg.IsRequestOrNotify(), "router only support request and notify")

	// query a connection by method name
	chRet := make(chan RetSelectConn, 1)
	self.chSelectConn <- CmdSelectConn{
		Method: msg.MustMethod(),
		ConnId: cmdMsg.ConnId,
		ChRet:  chRet,
	}

	if ret, ok := <-chRet; ok && ret.Found {
		ret.Conn.MsgInput() <- cmdMsg
	} else if msg.IsRequest() {
		reqMsg, _ := msg.(*jsonz.RequestMessage)
		errMsg := jsonz.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		cmdMsg.ChRes <- cmdMsg.Res(errMsg)
	} else {
		msg.Log().Warnf("fail to select connect")
	}
}

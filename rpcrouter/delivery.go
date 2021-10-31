package rpcrouter

import (
	"fmt"
	//log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	//"time"
)

func (self *Router) PostMessage(cmdMsg CmdMsg) {
	// there are two paradigms two post methods
	// 1. send to router and let router relay the message to correspondend connections
	// 2. query the router goroutine for connections and send message to conn's msg output port, this way is a little bit faster but suffers the channel closing panic

	//self.chRouteMsg <- cmdMsg
	self.redirectMessage(cmdMsg)
}

func (self *Router) relayMessage(cmdMsg CmdMsg) {
	msg := cmdMsg.MsgVec.Msg
	misc.Assert(msg.IsRequestOrNotify(), "router only support request and notify")
	toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.ConnId)
	if found {
		toConn.MsgInput() <- cmdMsg
	} else if msg.IsRequest() {
		reqMsg, _ := msg.(*jsonrpc.RequestMessage)
		errMsg := jsonrpc.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg}
		cmdMsg.ChRes <- errMsgVec
	}
}

func (self *Router) redirectMessage(cmdMsg CmdMsg) {
	msg := cmdMsg.MsgVec.Msg
	misc.Assert(msg.IsRequestOrNotify(), "router only support request and notify")

	chRet := make(chan RetSelectConn, 1)

	self.chSelectConn <- CmdSelectConn{
		Method: msg.MustMethod(),
		ConnId: cmdMsg.ConnId,
		ChRet:  chRet,
	}

	if ret, ok := <-chRet; ok && ret.Found {
		ret.Conn.MsgInput() <- cmdMsg
	} else if msg.IsRequest() {
		reqMsg, _ := msg.(*jsonrpc.RequestMessage)
		errMsg := jsonrpc.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg, Namespace: cmdMsg.MsgVec.Namespace}
		cmdMsg.ChRes <- errMsgVec
	} else {
		msg.Log().Warnf("fail to select connect")
	}
}

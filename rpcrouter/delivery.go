package rpcrouter

import (
	"fmt"
	//log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	//"time"
)

func (self *Router) relayMessage(cmdMsg CmdMsg) {
	msg := cmdMsg.MsgVec.Msg
	misc.Assert(msg.IsRequestOrNotify(), "router only support request and notify")
	toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.MsgVec.ToConnId)
	if found {
		toConn.ChRouteMsg <- cmdMsg
	} else if msg.IsRequest() {
		reqMsg, _ := msg.(*jsonrpc.RequestMessage)
		errMsg := jsonrpc.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg}
		cmdMsg.ChRes <- errMsgVec
	}
}

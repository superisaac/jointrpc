package rpcrouter

import (
	//"fmt"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/misc"
)

func NewConn() *ConnT {
	connId := NextCID()
	ch := make(MsgChannel, 10000)
	methods := make(map[string]MethodInfo)
	conn := &ConnT{ConnId: connId,
		RecvChannel:  ch,
		ServeMethods: methods,
	}
	return conn
}

func (self ConnT) Joined() bool {
	return self.Namespace != ""
}

func (self ConnT) GetMethods() []string {
	var keys []string
	for k := range self.ServeMethods {
		keys = append(keys, k)
	}
	return keys
}

func (self ConnT) ValidateRequestMsg(reqMsg *jsonrpc.RequestMessage) (bool, jsonrpc.IMessage) {
	if info, ok := self.ServeMethods[reqMsg.Method]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, jsonrpc.MessageInterface(reqMsg))
		if errPos != nil {
			errmsg := errPos.ToMessage(reqMsg)
			return false, errmsg
		}
	}
	return true, nil
}

func (self ConnT) ValidateNotifyMsg(notifyMsg *jsonrpc.NotifyMessage) (bool, error) {
	if info, ok := self.ServeMethods[notifyMsg.Method]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, jsonrpc.MessageInterface(notifyMsg))
		if errPos != nil {
			return false, errPos
		}
	}
	return true, nil
}

func (self ConnT) ValidateResultMsg(resMsg *jsonrpc.ResultMessage, reqMsg *jsonrpc.RequestMessage) (bool, jsonrpc.IMessage) {
	if info, ok := self.ServeMethods[reqMsg.Method]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, jsonrpc.MessageInterface(resMsg))
		if errPos != nil {
			errmsg := errPos.ToMessage(reqMsg)
			return false, errmsg
		}
	}
	return true, nil
}

func (self *ConnT) SetWatchState(w bool) {
	self.watchState = w
}

func (self *ConnT) StateChannel() chan *ServerState {
	if self.stateChannel == nil {
		self.stateChannel = make(chan *ServerState, misc.DefaultChanSize())
	}
	return self.stateChannel
}

func (self ConnT) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"namespace":  self.Namespace,
		"connid":     self.ConnId,
		"remoteaddr": self.PeerAddr,
	})
}

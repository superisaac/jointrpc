package rpcrouter

import (
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
)

func NewConn() *ConnT {
	connId := NextCID()
	ch := make(MsgChannel, 100)
	//chState := make(chan TuebState, 100)
	methods := make(map[string]MethodInfo)
	conn := &ConnT{ConnId: connId,
		RecvChannel:  ch,
		ServeMethods: methods,
		AsFallback:   false,
	}
	return conn
}

func (self ConnT) GetMethods() []string {
	var keys []string
	for k := range self.ServeMethods {
		keys = append(keys, k)
	}
	return keys
}

func (self ConnT) PublicId() string {
	return self.publicId
}

func (self ConnT) ValidateMsg(msg jsonrpc.IMessage) (bool, jsonrpc.IMessage) {
	if info, ok := self.ServeMethods[msg.MustMethod()]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, msg.Interface())
		if errPos != nil {
			if msg.IsRequest() {
				errmsg := errPos.ToMessage(msg)
				return false, errmsg
			} else {
				log.Warnf("validate error %s", errPos.Error())
				return false, nil
			}
		}
	}
	return true, nil
}

func (self *ConnT) SetWatchState(w bool) {
	self.watchState = w
}

func (self *ConnT) StateChannel() chan *ServerState {
	if self.stateChannel == nil {
		self.stateChannel = make(chan *ServerState, 100)
	}
	return self.stateChannel
}

func (self ConnT) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"conn_id":     self.ConnId,
		"public_id":   self.publicId,
		"remote_addr": self.PeerAddr,
	})
}

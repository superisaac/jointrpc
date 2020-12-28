package tube

import (
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	schema "github.com/superisaac/rpctube/jsonrpc/schema"
)

func NewConn() *ConnT {
	connId := NextCID()
	ch := make(MsgChannel, 100)
	//chState := make(chan TuebState, 100)
	methods := make(map[string]MethodInfo)
	conn := &ConnT{ConnId: connId,
		RecvChannel: ch,
		Methods:     methods,
		AsFallback:  false,
	}
	return conn
}

func (self ConnT) GetMethods() []string {
	var keys []string
	for k := range self.Methods {
		keys = append(keys, k)
	}
	return keys
}

func (self ConnT) ValidateMsg(msg jsonrpc.IMessage) (bool, jsonrpc.IMessage) {
	if info, ok := self.Methods[msg.MustMethod()]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, msg.Interface())
		if errPos != nil {
			if msg.IsRequest() {
				errmsg := errPos.ToMessage(msg.MustId())
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

func (self *ConnT) StateChannel() chan *TubeState {
	if self.stateChannel == nil {
		self.stateChannel = make(chan *TubeState, 100)
	}
	return self.stateChannel
}

package dispatch

import (
	//"errors"
	//log "github.com/sirupsen/logrus"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
)

func NewStateDispatcher() *StateDispatcher {
	disp := new(StateDispatcher)
	disp.stateHandlers = make([]StateHandlerFunc, 0)
	return disp
}

func (self *StateDispatcher) OnStateChange(stateChange StateHandlerFunc) {
	self.stateHandlers = append(self.stateHandlers, stateChange)
}

func (self *StateDispatcher) TriggerStateChange(state *rpcrouter.ServerState) {
	for _, f := range self.stateHandlers {
		f(state)
	}
}

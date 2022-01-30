package dispatch

import (
	//"errors"
	//log "github.com/sirupsen/logrus"
	//"github.com/superisaac/jsonz"
	//schema "github.com/superisaac/jsonz/schema"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
)

func NewStateListener() *StateListener {
	listener := new(StateListener)
	listener.stateHandlers = make([]StateHandlerFunc, 0)
	return listener
}

func (self *StateListener) OnStateChange(stateChange StateHandlerFunc) {
	self.stateHandlers = append(self.stateHandlers, stateChange)
}

func (self *StateListener) TriggerStateChange(state *rpcrouter.ServerState) {
	for _, f := range self.stateHandlers {
		f(state)
	}
}

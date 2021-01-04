package joint

import (
	log "github.com/sirupsen/logrus"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
)

// lazy style schema builder
func (self *MethodInfo) SchemaOrError() (schema.Schema, error) {
	if self.schemaObj == nil && self.SchemaJson != "" {
		builder := schema.NewSchemaBuilder()
		var err error
		self.schemaObj, err = builder.BuildBytes([]byte(self.SchemaJson))
		if err != nil {
			log.Warnf("error on building schema %s", self.SchemaJson)
			return nil, err
		}
	}
	return self.schemaObj, nil
}

func (self *MethodInfo) Schema() schema.Schema {
	s, err := self.SchemaOrError()
	if err != nil {
		panic(err)
	}
	return s
}

func (self *Router) NotifyStateChange() {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	self.notifyStateChange()
}

func (self *Router) notifyStateChange() {
	state := self.GetState()
	for _, conn := range self.connMap {
		if conn.watchState {
			conn.StateChannel() <- state
		}
	}
}

func (self *Router) GetState() *TubeState {
	state := &TubeState{Methods: self.getLocalMethods()}
	return state
}

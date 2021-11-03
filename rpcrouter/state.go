package rpcrouter

import (
	log "github.com/sirupsen/logrus"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
)

// lazy style schema builder
func (self *MethodInfo) SchemaOrError() (*schema.MethodSchema, error) {
	if self.schemaObj == nil && self.SchemaJson != "" {
		builder := schema.NewSchemaBuilder()
		var err error
		s, err := builder.BuildBytes([]byte(self.SchemaJson))
		if err != nil {
			log.Warnf("error on building schema %s", self.SchemaJson)
			return nil, err
		}
		if methodSchema, ok := s.(*schema.MethodSchema); ok {
			self.schemaObj = methodSchema
		} else {
			log.Warnf("is not method schema")
			return nil, schema.NewBuildError("method schema required", []string{})
		}
	}
	return self.schemaObj, nil
}

func (self *MethodInfo) Schema() *schema.MethodSchema {
	s, err := self.SchemaOrError()
	if err != nil {
		panic(err)
	}
	return s
}

func (self *Router) NotifyStateChange() {
	//self.rlock("notifystatechange")
	//defer self.runlock("notifystatechange")

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

func (self *Router) GetState() *ServerState {
	state := &ServerState{Methods: self.getMethods()}
	return state
}

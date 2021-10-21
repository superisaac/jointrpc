package client

import (
	//"context"
	//"errors"
	//"flag"
	//"fmt"
	//log "github.com/sirupsen/logrus"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"github.com/mitchellh/mapstructure"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/rpcrouter"
)

const (
	stateChangedSchema = `{
"type": "method",
"params": [{
   "type": "list",
   "items": {
     "type": "object",
     "properties": {
       "name": "string",
       "help": "string",
       "schema": "string" 
    },
    "requires": ["name"]
   }
}]
}`
)

func OnStateChanged(disp *dispatch.Dispatcher, stateListener *dispatch.StateListener) {
	disp.On("_state.changed",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			var serverState rpcrouter.ServerState
			err := mapstructure.Decode(params[0], &serverState)
			if err != nil {
				return nil, err
			}
			stateListener.TriggerStateChange(&serverState)
			return nil, nil
		}, dispatch.WithSchema(stateChangedSchema))
}

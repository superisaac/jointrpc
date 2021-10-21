package client

import (
	"context"
	"errors"
	//"flag"
	//"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"io"
	//"net/url"
	//"os"
	//"time"
	//server "github.com/superisaac/jointrpc/server"
	"github.com/mitchellh/mapstructure"
	"github.com/superisaac/jointrpc/dispatch"
	encoding "github.com/superisaac/jointrpc/encoding"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	//credentials "google.golang.org/grpc/credentials"
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

func (self *RPCClient) OldSubscribeState(rootCtx context.Context, stateListener *dispatch.StateListener) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	if self.stateStream != nil {
		return errors.New("state stream already exist")
	}

	authReq := &intf.AuthRequest{
		RequestId:  misc.NewUuid(),
		ClientAuth: self.ClientAuth(),
	}
	stream, err := self.grpcClient.SubscribeState(ctx, authReq, grpc_retry.WithMax(500))
	if err == io.EOF {
		log.Infof("cannot subscribe state")
		return nil
	} else if grpc.Code(err) == codes.Unavailable {
		log.Debugf("connect closed retrying")
		return nil
	} else if err != nil {
		log.Warnf("error on subscibe state %v", err)
		return err
	}

	self.stateStream = stream
	defer func() {
		self.stateStream = nil
	}()

	for {
		pac, err := self.stateStream.Recv()
		if err == io.EOF {
			log.Infof("state stream closed")
			return nil
		} else if grpc.Code(err) == codes.Unavailable {
			log.Debugf("state connect closed retrying")
			return nil
		} else if err != nil {
			log.Debugf("down pack error %+v %d", err, grpc.Code(err))
			return err
		}

		// Set connPublicId
		authResp := pac.GetAuthResponse()
		if authResp != nil {
			err := self.CheckStatus(authResp.Status, "SubscribeState.Auth")
			if err != nil {
				log.Warn(err.Error())
				return err
			}
			continue
		}

		istate := pac.GetState()
		if istate != nil {
			state := encoding.DecodeServerState(istate)
			stateListener.TriggerStateChange(state)
			continue
		}
	}

}

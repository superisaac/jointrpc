package server

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"net"
)

type StreamDispatcher struct {
	disp     *dispatch.Dispatcher
	authDisp *dispatch.Dispatcher
}

func NewStreamDispatcher() *StreamDispatcher {
	disp := dispatch.NewDispatcher()
	h := &StreamDispatcher{disp: disp}
	h.Init()
	return h
}

var connDisp *StreamDispatcher

func GetStreamDispatcher() *StreamDispatcher {
	if connDisp == nil {
		connDisp = NewStreamDispatcher()
	}
	return connDisp
}

const (
	declareMethodsSchema = `{
"type": "method",
"params": [{
  "type": "object",
  "properties": {
     "name": "string",
     "help": "string",
     "schema": "string"
    },
  "requires": ["name"]
  }],
"returns": "string"
}`

	declareDelegatesSchema = `{
"type": "method",
"params": [{
  "type": "list",
  "items": "string"
  }],
"returns": "string"
}`

	authorizeSchema = `{
"type": "method",
"params": ["string", "string"],
"returns": "string"
}`
)

func (self *StreamDispatcher) Init() {
	self.disp = dispatch.NewDispatcher()
	self.authDisp = dispatch.NewDispatcher()

	self.disp.On("_stream.declareMethods",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {

			conn, found := req.Data.(*rpcrouter.ConnT)
			if !found {
				return nil, jsonrpc.ParamsError("conn not found")
			}
			arr, ok := params[0].([]interface{})
			misc.Assert(ok, "params[0] is not an array")

			upMethods := make([]rpcrouter.MethodInfo, 0)
			var methodNames []string
			for _, infoDict := range arr {
				var minfo rpcrouter.MethodInfo
				err := mapstructure.Decode(infoDict, &minfo)
				if err != nil {
					return nil, err
				}
				if !jsonrpc.IsPublicMethod(minfo.Name) {
					conn.Log().WithFields(log.Fields{
						"rpc": "DeclareMethods",
					}).Warnf("%s is not valid public method name", minfo.Name)
					return nil, jsonrpc.ParamsError(fmt.Sprintf("method %s cannot prefix with .", minfo.Name))
				}
				methodNames = append(methodNames, minfo.Name)
				_, err = minfo.SchemaOrError()
				if err != nil {
					if buildError, ok := err.(*schema.SchemaBuildError); ok {
						// parse schema error
						conn.Log().WithFields(log.Fields{
							"rpc": "DeclareMethods",
						}).Warnf("error build schema %s, %+v", buildError.Error(), minfo)
						return nil, jsonrpc.ParamsError(fmt.Sprintf("build schema error %s", buildError.Error()))
					}
					return nil, err
				}
				upMethods = append(upMethods, minfo)
			}
			cmdMethods := rpcrouter.CmdMethods{
				Namespace: conn.Namespace,
				ConnId:    conn.ConnId,
				Methods:   upMethods,
			}

			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			factory.ChMethods <- cmdMethods
			return "ok", nil
		}, dispatch.WithSchema(declareMethodsSchema))

	self.disp.On("_stream.declareDelegates",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			conn, found := req.Data.(*rpcrouter.ConnT)
			if !found {
				return nil, jsonrpc.ParamsError("conn not found")
			}

			methodNames := jsonrpc.ConvertStringList(params[0])
			conn.Log().Infof("declared delegates %+v", methodNames)
			cmdDelegates := rpcrouter.CmdDelegates{
				Namespace:   conn.Namespace,
				ConnId:      conn.ConnId,
				MethodNames: methodNames,
			}
			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			factory.ChDelegates <- cmdDelegates
			return "ok", nil
		}, dispatch.WithSchema(declareDelegatesSchema))

	self.authDisp.On("_stream.authorize",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			conn, found := req.Data.(*rpcrouter.ConnT)
			if !found {
				return nil, jsonrpc.ParamsError("conn not found")
			}

			v := req.Context.Value("remoteAddress")
			remoteAddress := ""
			if v != nil {
				remoteAddr, isAddr := v.(net.Addr)
				misc.Assert(isAddr, "context value remoteAddress is not net.Addr")
				remoteAddress = remoteAddr.String()
			}

			username := jsonrpc.ConvertString(params[0])
			password := jsonrpc.ConvertString(params[1])

			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			cfg := factory.Config
			namespace := cfg.Authorize(username, password, remoteAddress)
			if namespace != "" {
				router := factory.Get(namespace)
				router.JoinConn(conn)
				conn.SetWatchState(true)
				return namespace, nil
			}
			return nil, jsonrpc.ErrAuthFailed
		}, dispatch.WithSchema(authorizeSchema))
} // end of Init()

func (self *StreamDispatcher) HandleMessage(ctx context.Context, msgvec rpcrouter.MsgVec, chResult chan dispatch.ResultT, conn *rpcrouter.ConnT) jsonrpc.IMessage {
	if !conn.Joined() {
		instRes := self.authDisp.Expect(ctx, msgvec, dispatch.WithRequestData(conn))
		return instRes
	} else {
		msg := msgvec.Msg
		if msg.IsRequestOrNotify() && self.disp.HasMethod(msg.MustMethod()) {
			self.disp.Feed(ctx, msgvec, chResult, dispatch.WithRequestData(conn))
		} else {
			factory := rpcrouter.RouterFactoryFromContext(ctx)
			router := factory.Get(conn.Namespace)
			router.DeliverMessage(rpcrouter.CmdMsg{MsgVec: msgvec})
		}
		return nil
	}
}

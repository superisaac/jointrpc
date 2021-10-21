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

type ConnDispatcher struct {
	disp     *dispatch.Dispatcher
	authDisp *dispatch.Dispatcher
}

func NewConnDispatcher() *ConnDispatcher {
	disp := dispatch.NewDispatcher()
	h := &ConnDispatcher{disp: disp}
	h.Init()
	return h
}

var connDisp *ConnDispatcher

func GetConnDispatcher() *ConnDispatcher {
	if connDisp == nil {
		connDisp = NewConnDispatcher()
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

func (self *ConnDispatcher) Init() {
	self.disp = dispatch.NewDispatcher()
	self.authDisp = dispatch.NewDispatcher()

	self.disp.On("_conn.declareMethods",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			router := factory.Get(req.MsgVec.Namespace)

			conn, found := router.GetConn(req.MsgVec.FromConnId)
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

			factory.ChMethods <- cmdMethods
			return "ok", nil
		}, dispatch.WithSchema(declareMethodsSchema))

	self.disp.On("_conn.declareDelegates",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			router := factory.Get(req.MsgVec.Namespace)

			conn, found := router.GetConn(req.MsgVec.FromConnId)
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
			factory.ChDelegates <- cmdDelegates
			return "ok", nil
		}, dispatch.WithSchema(declareDelegatesSchema))

	self.authDisp.On("_conn.Authorize",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
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
			if len(cfg.Authorizations) == 0 {
				return "default", nil
			}
			for _, bauth := range cfg.Authorizations {
				if bauth.Authorize(username, password, remoteAddress) {
					ns := bauth.Namespace
					if ns == "" {
						ns = "default"
					}
					return ns, nil
				}
			}
			return nil, jsonrpc.ErrAuthFailed
		}, dispatch.WithSchema(authorizeSchema))
} // end of Init()

func (self *ConnDispatcher) HandleRequest(ctx context.Context, msgvec rpcrouter.MsgVec, chResult chan dispatch.ResultT) bool {
	msg := msgvec.Msg
	if !msg.IsRequest() {
		return false
	}

	methodName := msg.MustMethod()
	if !self.disp.HasMethod(methodName) {
		return false
	}

	self.disp.Feed(ctx, msgvec, chResult)
	return true
}

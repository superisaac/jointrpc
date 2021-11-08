package server

import (
	"context"
	"fmt"
	//"time"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"net"
	"sync"
)

var (
	once       sync.Once
	streamDisp *StreamDispatcher
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

func GetStreamDispatcher() *StreamDispatcher {
	once.Do(func() {
		streamDisp = NewStreamDispatcher()
	})
	return streamDisp
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

	watchStateSchema = `{
"type": "method",
"params": [],
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

	self.disp.OnTyped("_stream.ping",
		func(req *dispatch.RPCRequest) (string, error) {
			req.CmdMsg.Msg.Log().Debugf("ping received")
			if conn, ok := req.Data.(*rpcrouter.ConnT); ok {
				conn.Touch()
			}
			return "pong", nil
		})

	self.disp.On("_stream.declareMethods",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {

			conn, found := req.Data.(*rpcrouter.ConnT)
			if !found {
				return nil, jsonrpc.ParamsError("conn not found")
			}
			conn.Log().Infof("call _stream.declareMethods")
			arr, ok := params[0].([]interface{})
			misc.Assert(ok, "params[0] is not an array")

			upMethods := make([]rpcrouter.MethodInfo, 0)
			var methodNames []string
			for _, infoDict := range arr {
				var minfo rpcrouter.MethodInfo
				err := misc.DecodeStruct(infoDict, &minfo)
				if err != nil {
					return nil, errors.Wrap(err, "misc.DecodeStruct()")
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
					var buildError *schema.SchemaBuildError
					if errors.As(err, &buildError) {
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
			factory.Get(conn.Namespace).ChMethods <- cmdMethods
			return "ok", nil
		}, dispatch.WithSchema(declareMethodsSchema))

	self.disp.On("_stream.declareDelegates",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			conn, found := req.Data.(*rpcrouter.ConnT)
			if !found {
				return nil, jsonrpc.ParamsError("conn not found")
			}

			methodNames := jsonrpc.ConvertStringList(params[0])
			conn.Log().Infof("call _stream.declareDelegates %+v", methodNames)
			cmdDelegates := rpcrouter.CmdDelegates{
				Namespace:   conn.Namespace,
				ConnId:      conn.ConnId,
				MethodNames: methodNames,
			}
			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			router := factory.Get(conn.Namespace)
			router.ChDelegates <- cmdDelegates

			misc.Assert(router.Started(), "router is not started")
			return "ok", nil
		}, dispatch.WithSchema(declareDelegatesSchema))

	self.disp.On("_stream.watchState",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			conn, found := req.Data.(*rpcrouter.ConnT)
			if !found {
				return nil, jsonrpc.ParamsError("conn not found")
			}
			conn.SetWatchState(true)
			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			router := factory.Get(conn.Namespace)
			go func() {
				state := router.GetState()
				conn.StateChannel() <- state
			}()
			return "ok", nil
		}, dispatch.WithSchema(watchStateSchema))

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

				chRet := make(chan rpcrouter.CmdRet, 1)
				router.ChJoin <- rpcrouter.CmdJoin{Conn: conn, ChRet: chRet}

				<-chRet
				conn.Log().Infof("joined to router %s", namespace)
				return namespace, nil
			}
			return nil, jsonrpc.ErrAuthFailed
		}, dispatch.WithSchema(authorizeSchema))
} // end of Init()

func (self *StreamDispatcher) HandleMessage(ctx context.Context, msg jsonrpc.IMessage, ns string, chResult chan dispatch.ResultT, conn *rpcrouter.ConnT, allowRequest bool) jsonrpc.IMessage {
	cmdMsg := rpcrouter.CmdMsg{Msg: msg, Namespace: ns}
	if !conn.Joined() {
		instRes := self.authDisp.Expect(ctx, cmdMsg, dispatch.WithRequestData(conn))
		return instRes
	} else {
		if msg.IsRequestOrNotify() && self.disp.HasMethod(msg.MustMethod()) {
			self.disp.Feed(ctx, cmdMsg, chResult, dispatch.WithRequestData(conn))
		} else if msg.IsRequest() && !allowRequest {
			instRes := jsonrpc.ErrNotAllowed.ToMessage(msg)
			return instRes
		} else if msg.IsNotify() && !allowRequest {
			msg.Log().Warnf("not allowed")
			return nil
		} else if msg.IsRequestOrNotify() {
			factory := rpcrouter.RouterFactoryFromContext(ctx)
			router := factory.Get(conn.Namespace)

			chRes := conn.MsgOutput()
			if msg.IsNotify() {
				chRes = nil
			}
			cmdMsg.ChRes = chRes
			router.PostMessage(cmdMsg)
		} else if msg.IsResultOrError() {
			// result and error don't need ChRes
			conn.MsgInput() <- cmdMsg
		}
		return nil
	}
}

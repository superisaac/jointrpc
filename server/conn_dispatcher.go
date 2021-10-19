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
		})

	self.disp.On("_conn.declareDelegates",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			factory := rpcrouter.RouterFactoryFromContext(req.Context)
			router := factory.Get(req.MsgVec.Namespace)

			conn, found := router.GetConn(req.MsgVec.FromConnId)
			if !found {
				return nil, jsonrpc.ParamsError("conn not found")
			}
			var arr []string
			if params[0] == nil {
				arr = make([]string, 0)
			} else {
				if iarr, ok := params[0].([]interface{}); ok {
					for _, item := range iarr {
						arr = append(arr, fmt.Sprintf("%s", item))
					}
				} else {
					req.MsgVec.Msg.Log().Warnf("params[0] is not array, =%+v", params[0])
					return nil, jsonrpc.ParamsError("params[0] is not an array")
				}
			}

			methodNames := make([]string, 0)

			for _, methodName := range arr {
				methodNames = append(methodNames, methodName)
			}

			conn.Log().Infof("declared delegates %+v", methodNames)
			cmdDelegates := rpcrouter.CmdDelegates{
				Namespace:   conn.Namespace,
				ConnId:      conn.ConnId,
				MethodNames: methodNames,
			}
			factory.ChDelegates <- cmdDelegates
			return "ok", nil
		})
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

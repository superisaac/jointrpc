package server

import (
	"fmt"
	"github.com/mitchellh/mapstructure"

	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/rpcrouter"
)

func handleDeclareMethods(factory *rpcrouter.RouterFactory, conn *rpcrouter.ConnT, req jsonrpc.IMessage) (interface{}, error) {
	params := req.MustParams()
	if len(params) != 1 {
		return nil, jsonrpc.ParamsError("invalid params length")
	}

	arr, ok := params[0].([]interface{})
	if !ok {
		return nil, jsonrpc.ParamsError("params[0] is not an array")
	}

	upMethods := make([]rpcrouter.MethodInfo, 0)

	for _, infoDict := range arr {
		var minfo rpcrouter.MethodInfo
		err := mapstructure.Decode(infoDict, &minfo)
		if err != nil {
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
}

func handleDeclareDelegates(factory *rpcrouter.RouterFactory, conn *rpcrouter.ConnT, req jsonrpc.IMessage) (interface{}, error) {
	params := req.MustParams()

	if len(params) != 1 {
		return nil, jsonrpc.ParamsError("invalid params length")
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
			req.Log().Warnf("params[0] is not array, =%+v", params[0])
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
}

func handleConnRequests(factory *rpcrouter.RouterFactory, conn *rpcrouter.ConnT, msg jsonrpc.IMessage) (jsonrpc.IMessage, error) {
	if !msg.IsRequest() {
		return nil, nil
	}

	var err error
	var r interface{}
	switch msg.MustMethod() {
	case "_conn.declareMethods":
		r, err = handleDeclareMethods(factory, conn, msg)
	case "_conn.declareDelegates":
		r, err = handleDeclareDelegates(factory, conn, msg)
	default:
		return nil, nil
	}

	if err != nil {
		if rpcError, ok := r.(*jsonrpc.RPCError); ok {
			errmsg := rpcError.ToMessage(msg)
			return errmsg, nil
		}
		return nil, err
	}
	resmsg := jsonrpc.NewResultMessage(msg, r, nil)
	return resmsg, nil

}

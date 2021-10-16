package server

import (
	//"fmt"
	"github.com/mitchellh/mapstructure"

	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/rpcrouter"
)

func handleDeclareMethods(factory *rpcrouter.RouterFactory, conn *rpcrouter.ConnT, req jsonrpc.IMessage) (res interface{}, err error) {
	params := req.MustParams()
	if len(params) != 1 {
		return nil, jsonrpc.ParamsError("invalid params length")
	}

	arr, ok := params[0].([]interface{})
	if !ok {
		return nil, jsonrpc.ParamsError("params[0] is not a array")
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

func handleConnRequests(factory *rpcrouter.RouterFactory, conn *rpcrouter.ConnT, msg jsonrpc.IMessage) (jsonrpc.IMessage, error) {
	if !msg.IsRequest() {
		return nil, nil
	}

	var err error
	var r interface{}
	switch msg.MustMethod() {
	case "_conn.declareMethods":
		r, err = handleDeclareMethods(factory, conn, msg)
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

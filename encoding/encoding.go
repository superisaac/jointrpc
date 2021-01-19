package encoding

import (
	//log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/rpcrouter"
)

// turn from tube's struct to protobuf message
func EncodeMethodInfo(minfo rpcrouter.MethodInfo) *intf.MethodInfo {
	intfInfo := &intf.MethodInfo{
		Name:       minfo.Name,
		Help:       minfo.Help,
		SchemaJson: minfo.SchemaJson,
	}
	return intfInfo
}

// turn from protobuf to tube's struct
func DecodeMethodInfo(iminfo *intf.MethodInfo) *rpcrouter.MethodInfo {
	return &rpcrouter.MethodInfo{
		Name:       iminfo.Name,
		Help:       iminfo.Help,
		SchemaJson: iminfo.SchemaJson,
	}
}

func EncodeServerState(state *rpcrouter.ServerState) *intf.ServerState {
	var resMethods []*intf.MethodInfo
	for _, minfo := range state.Methods {
		resMethods = append(resMethods, EncodeMethodInfo(minfo))
	}
	return &intf.ServerState{Methods: resMethods}
}

func DecodeServerState(iState *intf.ServerState) *rpcrouter.ServerState {
	var resMethods []rpcrouter.MethodInfo
	for _, info := range iState.Methods {
		minfo := DecodeMethodInfo(info)
		resMethods = append(resMethods, *minfo)
	}
	return &rpcrouter.ServerState{Methods: resMethods}
}

func MessageToEnvolope(msg jsonrpc.IMessage) *intf.JSONRPCEnvolope {
	return &intf.JSONRPCEnvolope{
		Body:    msg.MustString(),
		TraceId: msg.TraceId()}
}

func MessageFromEnvolope(envo *intf.JSONRPCEnvolope) (jsonrpc.IMessage, error) {
	msg, err := jsonrpc.ParseBytes([]byte(envo.Body))
	if err != nil {
		return nil, err
	}
	msg.SetTraceId(envo.TraceId)
	return msg, nil
}

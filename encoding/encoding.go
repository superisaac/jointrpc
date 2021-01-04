package encoding

import (
	//log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/joint"
)

// turn from tube's struct to protobuf message
func EncodeMethodInfo(minfo joint.MethodInfo) *intf.MethodInfo {
	intfInfo := &intf.MethodInfo{
		Name:       minfo.Name,
		Help:       minfo.Help,
		SchemaJson: minfo.SchemaJson,
	}
	return intfInfo
}

// turn from protobuf to tube's struct
func DecodeMethodInfo(iminfo *intf.MethodInfo) *joint.MethodInfo {
	return &joint.MethodInfo{
		Name:       iminfo.Name,
		Help:       iminfo.Help,
		SchemaJson: iminfo.SchemaJson,
	}
}

func EncodeTubeState(state *joint.TubeState) *intf.TubeState {
	var resMethods []*intf.MethodInfo
	for _, minfo := range state.Methods {
		resMethods = append(resMethods, EncodeMethodInfo(minfo))
	}
	return &intf.TubeState{Methods: resMethods}
}

func DecodeTubeState(iState *intf.TubeState) *joint.TubeState {
	var resMethods []joint.MethodInfo
	for _, info := range iState.Methods {
		minfo := DecodeMethodInfo(info)
		resMethods = append(resMethods, *minfo)
	}
	return &joint.TubeState{Methods: resMethods}
}

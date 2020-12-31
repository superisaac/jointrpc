package encoding

import (
	//log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	//schema "github.com/superisaac/rpctube/jsonrpc/schema"
	tube "github.com/superisaac/rpctube/tube"
)

// turn from tube's struct to protobuf message
func EncodeMethodInfo(minfo tube.MethodInfo) *intf.MethodInfo {
	intfInfo := &intf.MethodInfo{
		Name:       minfo.Name,
		Help:       minfo.Help,
		SchemaJson: minfo.SchemaJson,
	}
	return intfInfo
}

// turn from protobuf to tube's struct
func DecodeMethodInfo(iminfo *intf.MethodInfo) *tube.MethodInfo {
	return &tube.MethodInfo{
		Name:       iminfo.Name,
		Help:       iminfo.Help,
		SchemaJson: iminfo.SchemaJson,
	}
}

func EncodeTubeState(state *tube.TubeState) *intf.TubeState {
	var resMethods []*intf.MethodInfo
	for _, minfo := range state.Methods {
		resMethods = append(resMethods, EncodeMethodInfo(minfo))
	}
	return &intf.TubeState{Methods: resMethods}
}

func DecodeTubeState(iState *intf.TubeState) *tube.TubeState {
	var resMethods []tube.MethodInfo
	for _, info := range iState.Methods {
		minfo := DecodeMethodInfo(info)
		resMethods = append(resMethods, *minfo)
	}
	return &tube.TubeState{Methods: resMethods}
}

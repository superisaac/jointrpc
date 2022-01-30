package msgutil

import (
	//log "github.com/sirupsen/logrus"
	//"github.com/pkg/errors"
	//"fmt"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jsonz"
	//schema "github.com/superisaac/jsonz/schema"
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

func MessageToEnvolope(msg jsonz.Message) *intf.JSONRPCEnvolope {
	return &intf.JSONRPCEnvolope{
		Body: jsonz.MessageString(msg),
	}
}

func MessageFromEnvolope(envo *intf.JSONRPCEnvolope) (jsonz.Message, error) {

	msg, err := jsonz.ParseBytes([]byte(envo.Body))
	if err != nil {
		return nil, err
	}
	return msg, nil
}

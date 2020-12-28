package client

import (
	"context"
	//"io"
	//simplejson "github.com/bitly/go-simplejson"
	//log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	//grpc "google.golang.org/grpc"
	//codes "google.golang.org/grpc/codes"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	//server "github.com/superisaac/rpctube/server"
)

func (self *RPCClient) ListMethods(rootCtx context.Context) ([]*intf.MethodInfo, error) {
	req := &intf.ListMethodsRequest{}
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	res, err := self.tubeClient.ListMethods(ctx, req)
	if err != nil {
		return [](*intf.MethodInfo){}, err
	}

	return res.MethodInfos, nil
}

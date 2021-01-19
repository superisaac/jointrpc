package client

import (
	"context"
	"errors"
	"fmt"
	//"io"
	//simplejson "github.com/bitly/go-simplejson"
	//log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	//grpc "google.golang.org/grpc"
	//codes "google.golang.org/grpc/codes"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//server "github.com/superisaac/jointrpc/server"
)

func (self *RPCClient) ListMethods(rootCtx context.Context) ([]*intf.MethodInfo, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.ListMethodsRequest{}
	res, err := self.tubeClient.ListMethods(ctx, req)
	if err != nil {
		return [](*intf.MethodInfo){}, err
	}

	return res.Methods, nil
}

func (self *RPCClient) DeclareMethods(rootCtx context.Context, methodInfos [](*intf.MethodInfo)) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.DeclareMethodsRequest{
		ConnPublicId: self.connPublicId,
		Methods:      methodInfos}
	res, err := self.tubeClient.DeclareMethods(ctx, req)
	if err != nil {
		return err
	}
	if res.Error != nil && res.Error.Code != 0 {
		return errors.New(fmt.Sprintf("declare methods failed %d %s", res.Error.Code, res.Error.Reason))
	}
	return nil
}

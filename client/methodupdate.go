package client

import (
	"context"
	//"errors"
	//"fmt"
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
	req := &intf.ListMethodsRequest{Auth: self.ClientAuth()}
	res, err := self.grpcClient.ListMethods(ctx, req)
	if err != nil {
		return [](*intf.MethodInfo){}, err
	}

	if err := self.CheckStatus(res.Status, "ListMethods"); err != nil {
		return [](*intf.MethodInfo){}, err
	}

	return res.Methods, nil
}

// func (self *RPCClient) DeclareMethods(rootCtx context.Context, methodInfos [](*intf.MethodInfo)) error {
// 	ctx, cancel := context.WithCancel(rootCtx)
// 	defer cancel()
// 	req := &intf.DeclareMethodsRequest{
// 		Auth:         self.ClientAuth(),
// 		ConnPublicId: self.connPublicId,
// 		Methods:      methodInfos}
// 	res, err := self.rpcClient.DeclareMethods(ctx, req)
// 	if err != nil {
// 		return err
// 	}
// 	if err := self.CheckStatus(res.Status, "DeclareMethods"); err != nil {
// 		return err
// 	}
// 	return nil
// }

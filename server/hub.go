package server

import (
	context "context"
	//json "encoding/json"
	//"errors"
	//"fmt"
	intf "github.com/superisaac/rpctube/intf/tube"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
	//"time"
)

type MethodHub struct {
	intf.UnimplementedMethodHubServer
}

func (self *MethodHub) UpdateMethods(context context.Context, req *intf.MethodsDecl) (*intf.UpdateMethodsResponse, error) {
	tube.Hub().UpdateMethods(req.Entrypoint, req.Methods)
	res := &intf.UpdateMethodsResponse{Res: "ok"}
	return res, nil
}

func (self *MethodHub) SubscribeAllMethods(ctx context.Context, in *Empty, opts ...grpc.CallOption) (MethodHub_SubscribeAllMethodsClient, error) {
	
}

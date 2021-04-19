package client

import (
	"context"
	//"errors"
	//"fmt"
	"errors"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/misc"
)

func (self *RPCClient) ListDelegates(rootCtx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.ListDelegatesRequest{Auth: self.ClientAuth()}
	res, err := self.rpcClient.ListDelegates(ctx, req)
	if err != nil {
		return nil, err
	}
	err = self.CheckStatus(res.Status, "ListDelegates")
	if err != nil {
		return nil, err
	}

	return res.Delegates, nil
}

func (self *RPCClient) DeclareDelegates(rootCtx context.Context, methods []string) error {
	if self.workerStream == nil {
		return errors.New("worker stream not setup")
	}

	req := &intf.DeclareDelegatesRequest{Methods: methods, RequestId: misc.NewUuid()}
	payload := &intf.JointRPCUpPacket_DelegatesRequest{DelegatesRequest: req}
	uppac := &intf.JointRPCUpPacket{Payload: payload}
	self.DeliverUpPacket(uppac)
	return nil
}

// func (self *RPCClient) DeclareDelegates(rootCtx context.Context, methods []string) error {
// 	ctx, cancel := context.WithCancel(rootCtx)
// 	defer cancel()
// 	res, err := self.rpcClient.DeclareDelegates(ctx, req)
// 	if err != nil {
// 		return err
// 	}
// 	err = self.CheckStatus(res.Status, "DeclareDelegates")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

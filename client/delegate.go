package client

import (
	"context"
	//"errors"
	//"fmt"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
)

func (self *RPCClient) ListDelegates(rootCtx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.ListDelegatesRequest{Auth: self.ClientAuth()}
	res, err := self.tubeClient.ListDelegates(ctx, req)
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
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.DeclareDelegatesRequest{
		Auth:         self.ClientAuth(),
		ConnPublicId: self.connPublicId,
		Methods:      methods}
	res, err := self.tubeClient.DeclareDelegates(ctx, req)
	if err != nil {
		return err
	}
	err = self.CheckStatus(res.Status, "DeclareDelegates")
	if err != nil {
		return err
	}
	return nil
}

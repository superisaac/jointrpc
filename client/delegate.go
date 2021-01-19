package client

import (
	"context"
	"errors"
	"fmt"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
)

func (self *RPCClient) ListDelegates(rootCtx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.ListDelegatesRequest{}
	res, err := self.tubeClient.ListDelegates(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Delegates, err
}

func (self *RPCClient) DeclareDelegates(rootCtx context.Context, methods []string) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.DeclareDelegatesRequest{
		ConnPublicId: self.connPublicId,
		Methods:      methods}
	res, err := self.tubeClient.DeclareDelegates(ctx, req)
	if err != nil {
		return err
	}
	if res.Error != nil && res.Error.Code != 0 {
		return errors.New(fmt.Sprintf("declared failed %d %s", res.Error.Code, res.Error.Reason))
	}
	return nil
}

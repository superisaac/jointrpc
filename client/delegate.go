package client

import (
	"context"
	//"errors"
	//"fmt"
	"errors"
	//log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/misc"
)

func (self *RPCClient) ListDelegates(rootCtx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.ListDelegatesRequest{Auth: self.ClientAuth()}
	res, err := self.grpcClient.ListDelegates(ctx, req)
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

	reqId := misc.NewUuid()
	if methods == nil {
		methods = make([]string, 0)
	}
	params := [](interface{}){methods}

	reqmsg := jsonrpc.NewRequestMessage(reqId, "_conn.declareDelegates", params, nil)

	return self.CallInWire(rootCtx, reqmsg, func(res jsonrpc.IMessage) {
		res.Log().Debugf("declared delegates")
	})
}

package client

import (
	"context"
	//"fmt"
	"github.com/pkg/errors"
	//log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jsonz"
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
	if !self.connected {
		return errors.New("live stream not setup")
	}

	reqId := misc.NewUuid()
	if methods == nil {
		methods = make([]string, 0)
	}
	params := [](interface{}){methods}

	reqmsg := jsonz.NewRequestMessage(reqId, "_stream.declareDelegates", params)

	return self.LiveCall(rootCtx, reqmsg, func(res jsonz.Message) {
		res.Log().Debugf("declared delegates")
	})
}

package client

import (
	"context"
	"io"
	//simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
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

func (self *RPCClient) WatchMethods(rootCtx context.Context) (MethodUpdateReceiver, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	//defer cancel()

	req := &intf.WatchMethodsRequest{}
	stream, err := self.tubeClient.WatchMethods(ctx, req)
	if err != nil {
		log.Warnf("error on watch methods %+v", err)
		return nil, err
	}
	ch := make(MethodUpdateReceiver, 100)
	go func() {
		defer cancel()
		for {
			update, err := stream.Recv()
			if err == io.EOF {
				log.Infof("watch methods stream closed")
				close(ch)
				return
			} else if err != nil {
				close(ch)
				log.Debugf("error code %d", grpc.Code(err))
				if grpc.Code(err) == codes.Unavailable {
					log.Warnf("server unavailable")
					return
				}
				panic(err)
			}
			log.Debugf("got method update %+v", update)
			ch <- update.MethodInfos
		}
	}()
	return ch, nil
}

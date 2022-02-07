package server

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonz"
)

func ReceiverLoop(ctx context.Context, sender dispatch.ISender, receiver IReceiver, conn *rpcrouter.ConnT, chResult chan dispatch.ResultT) {
	streamDisp := NewStreamDispatcher()
	for {
		//msg, err := msgutil.WSRecv(sender.ws)
		msg, err := receiver.Recv()
		if err != nil {
			log.Warnf("bad request %s", err)
			sender.Done() <- err
			return
		} else if msg == nil {
			sender.Done() <- nil
			return
		}

		if msg.TraceId() == "" {
			msg.SetTraceId(jsonz.NewUuid())
		}

		instRes := streamDisp.HandleMessage(
			ctx, msg,
			conn.Namespace,
			chResult,
			conn, false)

		if instRes != nil {
			err := sender.SendMessage(ctx, instRes)
			if err != nil {
				sender.Done() <- err
				return
			}
			if instRes.IsError() {
				sender.Done() <- nil
				return
			}
		}
	} // end of for
}

func WaitStream(rootCtx context.Context, sender dispatch.ISender, receiver IReceiver, conn *rpcrouter.ConnT) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	chResult := make(chan dispatch.ResultT, misc.DefaultChanSize())
	go dispatch.SenderLoop(ctx, sender, conn, chResult)
	go ReceiverLoop(ctx, sender, receiver, conn, chResult)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err, ok := <-sender.Done():
			if !ok {
				log.Debugf("done received not ok")
				return nil
			} else if err != nil {
				log.Errorf("stream err %+v", err)
				return err
			}
		}
	}
}

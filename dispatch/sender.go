package dispatch

import (
	"context"
	//"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonz"
	"time"
)

func SenderLoop(rootCtx context.Context, sender ISender, conn *rpcrouter.ConnT, chResult chan ResultT) {
	if true {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("recovered ERROR %+v", r)
			}
		}()
	}

	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Debugf("context done")
			return
		case rest, ok := <-chResult:
			{
				if !ok {
					log.Debugf("conn handler channel closed")
					sender.Done() <- nil
					return
				}
				err := sender.SendMessage(ctx, rest.ResMsg)
				if err != nil {
					sender.Done() <- err
					return
				}
			}
		case cmdMsg, ok := <-conn.MsgOutput():
			{
				if !ok {
					log.Debugf("recv channel closed")
					sender.Done() <- nil
					return
				}
				//err := sender.SendMessage(ctx, cmdMsg.Msg)
				err := sender.SendCmdMsg(ctx, cmdMsg)
				if err != nil {
					//panic(err)
					sender.Done() <- err
					return
				}
			}
		case cmdMsg, ok := <-conn.MsgInput():
			{
				if !ok {
					log.Debugf("MsgInput() closed")
					sender.Done() <- nil
					return
				}
				err := conn.HandleRouteMessage(ctx, cmdMsg)
				if err != nil {
					//panic(err)
					sender.Done() <- err
					return
				}
			}
		case state, ok := <-conn.StateChannel():
			{
				if !ok {
					log.Debugf("state channel closed")
					sender.Done() <- nil
					return
				}
				stateJson := make(map[string]interface{})
				err := misc.DecodeStruct(state, &stateJson)
				if err != nil {
					//panic(err)
					sender.Done() <- err
					return
				}
				ntf := jsonz.NewNotifyMessage("_state.changed", []interface{}{stateJson})
				err = sender.SendMessage(ctx, ntf)
				if err != nil {
					//panic(err)
					sender.Done() <- err
					return
				}
			}
		case <-time.After(10 * time.Second):
			conn.ClearPendings()
		}
	} // and for loop
}

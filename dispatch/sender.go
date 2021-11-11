package dispatch

import (
	"context"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"time"
)

func SenderLoop(ctx context.Context, sender ISender, conn *rpcrouter.ConnT, chResult chan ResultT) {
	if true {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("recovered ERROR %+v", r)
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			log.Debugf("context done")
			return
		case rest, ok := <-chResult:
			{
				if !ok {
					log.Debugf("conn handler channel closed")
					return
				}
				err := sender.SendMessage(ctx, rest.ResMsg)
				if err != nil {
					panic(err)
				}
			}
		case cmdMsg, ok := <-conn.MsgOutput():
			{
				if !ok {
					log.Debugf("recv channel closed")
					return
				}
				//err := sender.SendMessage(ctx, cmdMsg.Msg)
				err := sender.SendCmdMsg(ctx, cmdMsg)
				if err != nil {
					panic(err)
				}
			}
		case cmdMsg, ok := <-conn.MsgInput():
			{
				if !ok {
					log.Debugf("MsgInput() closed")
					return
				}
				err := conn.HandleRouteMessage(ctx, cmdMsg)
				if err != nil {
					panic(err)
				}
			}
		case state, ok := <-conn.StateChannel():
			{
				if !ok {
					log.Debugf("state channel closed")
					return
				}
				stateJson := make(map[string]interface{})
				err := misc.DecodeStruct(state, &stateJson)
				if err != nil {
					panic(err)
				}
				ntf := jsonrpc.NewNotifyMessage("_state.changed", []interface{}{stateJson})
				err = sender.SendMessage(ctx, ntf)
				if err != nil {
					panic(err)
				}
			}
		case <-time.After(10 * time.Second):
			conn.ClearPendings()
		}
	} // and for loop
}

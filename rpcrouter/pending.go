package rpcrouter

import (
	"time"
	//"context"
	//log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

// func (self *Router) TryClearPendingRequest(msgId interface{}) {
// 	self.pendingLock.RLock()
// 	defer self.pendingLock.RUnlock()
// 	log.Debugf("try to clear pending request %#v", msgId)

// 	if _, found := self.pendingRequests[msgId]; found {
// 		log.Infof("found pending req %#v", msgId)
// 		go func() {
// 			// sleep for another 1 second
// 			time.Sleep(1 * time.Second)
// 			self.ClearPendingRequest(msgId)
// 		}()
// 	}
// }

// func (self *Router) ClearPendingRequest(msgId interface{}) {
// 	self.lockPending("ClearPendingRequest")
// 	defer self.unlockPending("ClearPendingRequest")

// 	if reqt, found := self.pendingRequests[msgId]; found {
// 		now := time.Now()
// 		if !now.After(reqt.Expire) {
// 			reqt.ReqMsg.Log().Errorf("Expire is not reached even during collecting routine")
// 		}
//		errMsg := jsonrpc.ErrTimeout.ToMessage(reqt.ReqMsg)
// 		msgvec := MsgVec{Msg: errMsg, ToConnId: reqt.FromConnId}
// 		//_ = self.SendTo(reqt.FromConnId, msgvec)
// 		go self.DeliverResultOrError(msgvec)

// 	}
// }

func (self *Router) DeletePending(msgId interface{}) {
	self.lockPending("DeletePending")
	defer self.unlockPending("DeletePending")
	delete(self.pendingRequests, msgId)
}

func (self *Router) getAndDeletePendings(msgId interface{}) (PendingT, bool) {
	self.lockPending("getAndDeletePendings")
	defer self.unlockPending("getAndDeletePendings")

	if reqt, ok := self.pendingRequests[msgId]; ok {
		delete(self.pendingRequests, msgId)
		return reqt, ok
	}
	return PendingT{}, false
}

func (self *Router) addPending(msgId interface{}, pending PendingT) {
	self.lockPending("addPending")
	defer self.unlockPending("addPending")

	self.pendingRequests[msgId] = pending
}

// func (self *Router) StartCollectPendings(rootCtx context.Context) {
// 	ctx, cancel := context.WithCancel(rootCtx)
// 	defer cancel()

// 	select {
// 	case <- ctx.Done():
// 		return
// 	case <- time.After(1 * time.Second):
// 		self.collectPendings()
// 	}
// }

func (self *Router) collectPendings() {
	self.pendingLock.RLock()
	defer self.pendingLock.RUnlock()

	now := time.Now()
	expired := make([]interface{}, 0)
	for msgId, reqt := range self.pendingRequests {
		if now.After(reqt.Expire) {
			expired = append(expired, msgId)
		}
	}
	if len(expired) > 1 {
		go self.removeExpiredPendings(expired)
	}
}

func (self *Router) removeExpiredPendings(expiredMsgIds []interface{}) {
	self.lockPending("removeExpiredPendings")
	defer self.unlockPending("removeExpiredPendings")

	now := time.Now()
	for _, msgId := range expiredMsgIds {
		if reqt, found := self.pendingRequests[msgId]; found {
			if now.After(reqt.Expire) {
				reqt.ReqMsg.Log().Infof("removed from pending due to timeout")
				delete(self.pendingRequests, msgId)
				errMsg := jsonrpc.ErrTimeout.ToMessage(reqt.ReqMsg)
				msgvec := MsgVec{Msg: errMsg, ToConnId: reqt.FromConnId}
				go self.DeliverResultOrError(msgvec)
			} else {

			}
		}
	}
}

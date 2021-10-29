package rpcrouter

import (
	"time"
	//"context"
	//log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

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

func (self *Router) collectPendings() {
	self.rlockPending("collectPendings")
	defer self.runlockPending("collectPendings")

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
				//go self.DeliverResultOrError(msgvec)
				self.ChMsg <- CmdMsg{MsgVec: msgvec}
			} else {

			}
		}
	}
}

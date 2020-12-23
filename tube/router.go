package tube

import (
	//"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewRouter() *Router {
	return new(Router).Init()
}

func RemoveConn(slice []MethodDesc, conn *ConnT) []MethodDesc {
	// for i := range slice {
	// 	if slice[i].Conn == conn {
	// 		//if slice[i].Conn.ConnId == conn.ConnId {
	// 		slice = append(slice[:i], slice[i+1:]...)
	// 	}
	// }
	// return slice
	newarr := make([]MethodDesc, 0, len(slice)-1)
	for _, desc := range slice {
		if desc.Conn != conn {
			newarr = append(newarr, desc)
		}
	}
	return newarr
}

func (self *Router) Init() *Router {
	self.routerLock = new(sync.RWMutex)
	self.MethodConnMap = make(map[string]([]MethodDesc))
	self.ConnMap = make(map[CID](*ConnT))
	self.PendingMap = make(map[PendingKey]PendingValue)
	self.localMethodsSig = ""
	self.setupChannels()
	return self
}

func (self *Router) setupChannels() {
	self.ChMsg = make(chan CmdMsg, 1000)
	//self.ChLeave = make(chan CmdLeave, 100)
	self.ChUpdate = make(chan CmdUpdate, 1000)
}

func (self Router) GetAllMethods() []string {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	methods := []string{}
	for method, _ := range self.MethodConnMap {
		methods = append(methods, method)
	}
	sort.Strings(methods)
	return methods
}

func (self MethodDesc) IsLocal() bool {
	return !self.Info.Delegated
}

func (self MethodInfo) ToMap() MethodInfoMap {
	var schemaIntf interface{}
	if self.Schema != nil {
		schemaIntf = self.Schema.RebuildType()
	}
	return MethodInfoMap{
		"name":      self.Name,
		"help":      self.Help,
		"schema":    schemaIntf,
		"delegated": self.Delegated,
	}
}

func (self Router) GetLocalMethods() []MethodInfo {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	return self.getLocalMethods()
}

func (self Router) getLocalMethods() []MethodInfo {
	minfos := []MethodInfo{}
	for _, descs := range self.MethodConnMap {
		for _, desc := range descs {
			if desc.IsLocal() {
				minfos = append(minfos, desc.Info)
			}
		}
	}
	sort.Slice(minfos, func(i, j int) bool { return minfos[i].Name < minfos[j].Name })
	return minfos
}

func (self Router) getLocalMethodsSig() string {
	var arr []string
	for _, minfo := range self.getLocalMethods() {
		arr = append(arr, minfo.Name)
	}
	return strings.Join(arr, ",")
}

func (self *Router) UpdateMethods(conn *ConnT, methods []MethodInfo) bool {
	self.lock("UpdateMethods")
	defer self.unlock("UpdateMethods")
	return self.updateMethods(conn, methods)
}

func (self *Router) updateMethods(conn *ConnT, methods []MethodInfo) bool {
	connMethods := make(map[string]MethodInfo)
	addingMethods := make([]MethodInfo, 0)
	deletingMethods := make([]string, 0)

	// Find new methods
	for _, minfo := range methods {
		connMethods[minfo.Name] = minfo
		if _, found := conn.Methods[minfo.Name]; !found {
			addingMethods = append(addingMethods, minfo)
		}
	}
	// find old methods to delete
	for method, _ := range conn.Methods {
		if _, found := connMethods[method]; !found {
			deletingMethods = append(deletingMethods, method)
		}
	}

	conn.Methods = connMethods
	maybeChanged := len(addingMethods) > 0 || len(deletingMethods) > 0

	// add methods
	for _, minfo := range addingMethods {
		method := minfo.Name
		methodDesc := MethodDesc{
			Conn: conn,
			Info: minfo,
		}
		// bi direction map
		methodDescArr, methodFound := self.MethodConnMap[method]

		if methodFound {
			methodDescArr = append(methodDescArr, methodDesc)
		} else {
			var tmp []MethodDesc
			methodDescArr = append(tmp, methodDesc)
		}
		self.MethodConnMap[method] = methodDescArr
	}
	// delete methods
	for _, method := range deletingMethods {
		methodDescList, ok := self.MethodConnMap[method]
		if !ok {
			continue
		}
		methodDescList = RemoveConn(methodDescList, conn)
		if len(methodDescList) > 0 {
			self.MethodConnMap[method] = methodDescList
		} else {
			delete(self.MethodConnMap, method)
		}
	}
	if maybeChanged {
		sig := self.getLocalMethodsSig()
		if self.localMethodsSig != sig {
			// notify local methods change by broadcasting notification
			log.Debugf("local methods sig changed %s, %s", self.localMethodsSig, sig)
			self.localMethodsSig = sig
			params := [](interface{}){sig}
			notify := jsonrpc.NewNotifyMessage("localMethods.changed", params)
			self.ChMsg <- CmdMsg{
				MsgVec:    MsgVec{Msg: notify, FromConnId: 0},
				Broadcast: true,
			}
		}
	}

	return maybeChanged

}

func (self *Router) leaveConn(conn *ConnT) {
	for method, _ := range conn.Methods {
		methodDescList, ok := self.MethodConnMap[method]
		if !ok {
			continue
		}
		methodDescList = RemoveConn(methodDescList, conn)
		if len(methodDescList) > 0 {
			self.MethodConnMap[method] = methodDescList
		} else {
			delete(self.MethodConnMap, method)
		}
	}
	conn.Methods = make(map[string]MethodInfo)

	ct, ok := self.ConnMap[conn.ConnId]
	if ok {
		delete(self.ConnMap, conn.ConnId)
		close(ct.RecvChannel)
	}
}

func (self *Router) SelectConn(method string) (*ConnT, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	descs, ok := self.MethodConnMap[method]
	if ok && len(descs) > 0 {
		// or random or round-robin
		return descs[0].Conn, true
	}
	return nil, false
}

func (self *Router) SelectReceiver(method string) (MsgChannel, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	descs, ok := self.MethodConnMap[method]
	if ok && len(descs) > 0 {
		// or random or round-robin
		conn := descs[0].Conn
		return conn.RecvChannel, true
	}
	return nil, false
}

func (self *Router) ClearTimeoutRequests() {
	now := time.Now()
	tmpMap := make(map[PendingKey]PendingValue)

	for pKey, pValue := range self.PendingMap {
		if now.After(pValue.Expire) {
			errMsg := jsonrpc.NewErrorMessage(pKey.MsgId, 408, "request timeout", true)
			msgvec := MsgVec{errMsg, CID(0)}
			_ = self.deliverMessage(pKey.ConnId, msgvec)
		} else {
			tmpMap[pKey] = pValue
		}
	}
	self.PendingMap = tmpMap
}

func (self *Router) ClearPending(connId CID) {
	for pKey, pValue := range self.PendingMap {
		if pKey.ConnId == connId || pValue.ConnId == connId {
			self.deletePending(pKey)
		}
	}
}

func (self *Router) deletePending(pKey PendingKey) {
	delete(self.PendingMap, pKey)
}

func (self *Router) setPending(pKey PendingKey, pValue PendingValue) {
	self.PendingMap[pKey] = pValue
}

func (self *Router) routeMessage(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	fromConnId := cmdMsg.MsgVec.FromConnId
	if msg.IsRequest() {
		toConn, found := self.SelectConn(msg.Method)
		if found {
			pKey := PendingKey{ConnId: fromConnId, MsgId: msg.Id}
			expireTime := time.Now().Add(DefaultRequestTimeout)
			pValue := PendingValue{ConnId: toConn.ConnId, Expire: expireTime}

			self.setPending(pKey, pValue)
			return self.deliverMessage(toConn.ConnId, cmdMsg.MsgVec)
		} else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "method not found", false)
			errMsgVec := MsgVec{errMsg, CID(0)}
			return self.deliverMessage(fromConnId, errMsgVec)
		}
	} else if msg.IsNotify() {
		if cmdMsg.Broadcast {
			self.broadcastMessage(cmdMsg.MsgVec)
		} else {
			toConn, found := self.SelectConn(msg.Method)
			if found {
				return self.deliverMessage(
					toConn.ConnId, cmdMsg.MsgVec)
			}
		}
	} else if msg.IsResultOrError() {
		for pKey, pValue := range self.PendingMap {
			if pKey.MsgId == msg.Id && pValue.ConnId == fromConnId {
				// delete key within a range loop is safe
				// refer to https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-golang-map-within-a-range-loop
				self.deletePending(pKey)
				return self.deliverMessage(
					pKey.ConnId, cmdMsg.MsgVec)
			}
		} // end of for
	}
	return nil
}

func (self *Router) deliverMessage(connId CID, msgvec MsgVec) *ConnT {
	ct, ok := self.ConnMap[connId]
	if ok {
		ct.RecvChannel <- msgvec
		return ct
	}
	return nil
}

func (self *Router) broadcastMessage(msgvec MsgVec) int {
	log.Debugf("broadcast message %+v", msgvec.Msg)
	cnt := 0
	for _, ct := range self.ConnMap {
		if ct.ConnId == msgvec.FromConnId {
			// skip the from addr
			continue
		}
		_, ok := ct.Methods[msgvec.Msg.Method]
		if ok {
			cnt += 1
			ct.RecvChannel <- msgvec
		}
	}
	return cnt
}

func (self *Router) Start(ctx context.Context) {
	//self.setupChannels()
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debugf("Router goroutine done")
				return
				/*case cmd_join := <-self.ChJoin:
				self.Join(cmd_join.ConnId, cmd_join.RecvChannel) */
			// case cmd_leave, ok := <-self.ChLeave:
			// 	if !ok {
			// 		log.Warnf("ChLeave channel not ok")
			// 		return
			// 	}
			// 	conn, found := self.ConnMap[cmd_leave.ConnId]
			// 	if found {
			// 		self.Leave(conn)
			// 	}
			case cmd_update, ok := <-self.ChUpdate:
				{
					if !ok {
						log.Warnf("ChUpdate channel not ok")
						return
					}
					conn, found := self.ConnMap[cmd_update.ConnId]
					if found {
						self.UpdateMethods(conn, cmd_update.Methods)
					} else {
						log.Infof("Conn %d not found for update methods", cmd_update.ConnId)
					}
				}

			case cmd_msg, ok := <-self.ChMsg:
				{
					if !ok {
						log.Warnf("ChMsg channel not ok")
						return
					}

					self.RouteMessage(cmd_msg)
				}
			}
		}
	}()
}

// commands
func (self *Router) RouteMessage(cmdMsg CmdMsg) *ConnT {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	//msg.FromConnId = fromConnId
	//self.ChMsg <- msg
	return self.routeMessage(cmdMsg)
}

func (self *Router) Join() *ConnT {
	conn := NewConn()
	self.JoinConn(conn)
	return conn
}

func (self *Router) JoinConn(conn *ConnT) {
	self.lock("JoinConn")
	defer self.unlock("JoinConn")
	self.ConnMap[conn.ConnId] = conn
}

func (self *Router) lock(wrapper string) {
	//log.Printf("router want lock %s", wrapper)
	self.routerLock.Lock()
	//log.Printf("router locked %s", wrapper)
}
func (self *Router) unlock(wrapper string) {
	//log.Printf("router want unlock %s", wrapper)
	self.routerLock.Unlock()
	//log.Printf("router want unlocked %s", wrapper)
}

func (self *Router) Leave(conn *ConnT) {
	self.lock("Leave")
	defer self.unlock("Leave")

	self.leaveConn(conn)
}

func (self *Router) SingleCall(reqmsg *jsonrpc.RPCMessage, broadcast bool) (*jsonrpc.RPCMessage, error) {
	if reqmsg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		self.ChMsg <- CmdMsg{
			MsgVec: MsgVec{Msg: reqmsg, FromConnId: conn.ConnId}}
		msgvec := <-conn.RecvChannel
		return msgvec.Msg, nil
	} else {
		self.ChMsg <- CmdMsg{
			MsgVec:    MsgVec{Msg: reqmsg, FromConnId: 0},
			Broadcast: broadcast}
		return nil, nil
	}
}

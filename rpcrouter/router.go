package rpcrouter

import (
	//"fmt"
	"context"
	uuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewRouter(name string) *Router {
	return new(Router).Init(name)
}

func RouterFromContext(ctx context.Context) *Router {
	if v := ctx.Value("router"); v != nil {
		if router, ok := v.(*Router); ok {
			return router
		}
		panic("context value router is not a router instance")
	}
	panic("context does not have router")
}

func RemoveConn(slice []MethodDesc, conn *ConnT) []MethodDesc {
	newarr := make([]MethodDesc, 0, len(slice)-1)
	for _, desc := range slice {
		if desc.Conn != conn {
			newarr = append(newarr, desc)
		}
	}
	return newarr
}

func DelegateRemoveConn(slice []MethodDelegation, conn *ConnT) []MethodDelegation {
	newarr := make([]MethodDelegation, 0, len(slice)-1)
	for _, desc := range slice {
		if desc.Conn != conn {
			newarr = append(newarr, desc)
		}
	}
	return newarr
}

func (self *Router) Init(name string) *Router {
	self.name = name
	self.routerLock = new(sync.RWMutex)
	self.methodConnMap = make(map[string]([]MethodDesc))
	self.delegateConnMap = make(map[string][]MethodDelegation)
	self.fallbackConns = make([]*ConnT, 0)
	self.connMap = make(map[CID](*ConnT))
	self.publicConnMap = make(map[string](*ConnT))

	self.pendingRequests = make(map[interface{}]PendingT)
	self.methodsSig = ""
	self.setupChannels()
	return self
}

func (self *Router) setupChannels() {
	self.ChMsg = make(chan CmdMsg, 1000)
	//self.ChLeave = make(chan CmdLeave, 100)
	self.ChServe = make(chan CmdServe, 1000)
	self.ChDelegate = make(chan CmdDelegate, 1000)
}

func (self Router) Name() string {
	return self.name
}

func (self MethodInfo) ToMap() MethodInfoMap {
	var schemaIntf interface{}
	if self.SchemaJson != "" {
		schemaIntf = self.Schema().RebuildType()
	}
	return MethodInfoMap{
		"name":   self.Name,
		"help":   self.Help,
		"schema": schemaIntf,
	}
}

func (self Router) GetDelegates() []string {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	var arr []string
	for name, _ := range self.delegateConnMap {
		arr = append(arr, name)
	}
	return arr
}
func (self Router) GetMethods() []MethodInfo {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	return self.getMethods()
}

func (self Router) GetMethodNames() []string {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	methods := []string{}
	for method, _ := range self.methodConnMap {
		methods = append(methods, method)
	}
	sort.Strings(methods)
	return methods
}

func (self Router) getMethods() []MethodInfo {
	minfos := []MethodInfo{}
	for _, descs := range self.methodConnMap {
		for _, desc := range descs {
			minfos = append(minfos, desc.Info)
		}
	}
	sort.Slice(minfos, func(i, j int) bool { return minfos[i].Name < minfos[j].Name })
	return minfos
}

func (self Router) getMethodsSig() string {
	var arr []string
	var dup map[string]bool = map[string]bool{}
	for _, minfo := range self.getMethods() {
		if _, ok := dup[minfo.Name]; !ok {
			arr = append(arr, minfo.Name)
			dup[minfo.Name] = true
		}
	}
	return strings.Join(arr, ",")
}

func (self *Router) UpdateServeMethods(conn *ConnT, methods []MethodInfo) bool {
	self.lock("CanServeMethods")
	defer self.unlock("CanServeMethods")
	return self.updateServeMethods(conn, methods)
}

func (self *Router) updateServeMethods(conn *ConnT, methods []MethodInfo) bool {
	connMethods := make(map[string]MethodInfo)
	addingMethods := make([]MethodInfo, 0)
	deletingMethods := make([]string, 0)

	// Find new methods
	for _, minfo := range methods {
		connMethods[minfo.Name] = minfo
		if _, found := conn.ServeMethods[minfo.Name]; !found {
			addingMethods = append(addingMethods, minfo)
		}
	}
	// find old methods to delete
	for method, _ := range conn.ServeMethods {
		if _, found := connMethods[method]; !found {
			deletingMethods = append(deletingMethods, method)
		}
	}

	conn.ServeMethods = connMethods
	maybeChanged := len(addingMethods) > 0 || len(deletingMethods) > 0

	// add methods
	for _, minfo := range addingMethods {
		method := minfo.Name
		methodDesc := MethodDesc{
			Conn: conn,
			Info: minfo,
		}
		// bi direction map
		methodDescArr, methodFound := self.methodConnMap[method]

		if methodFound {
			methodDescArr = append(methodDescArr, methodDesc)
		} else {
			var tmp []MethodDesc
			methodDescArr = append(tmp, methodDesc)
		}
		self.methodConnMap[method] = methodDescArr
	}
	// delete methods
	for _, method := range deletingMethods {
		methodDescList, ok := self.methodConnMap[method]
		if !ok {
			continue
		}
		methodDescList = RemoveConn(methodDescList, conn)
		if len(methodDescList) > 0 {
			self.methodConnMap[method] = methodDescList
		} else {
			delete(self.methodConnMap, method)
		}
	}
	if maybeChanged {
		self.probeMethodChange()
	}

	return maybeChanged
}

func (self *Router) UpdateDelegateMethods(conn *ConnT, methodNames []string) bool {
	self.lock("CanDelegateMethods")
	defer self.unlock("CanDelegateMethods")
	return self.updateDelegateMethods(conn, methodNames)
}

func (self *Router) updateDelegateMethods(conn *ConnT, methodNames []string) bool {
	connMethods := make(map[string]bool)
	addingMethods := make([]string, 0)
	deletingMethods := make([]string, 0)

	// Find new methods
	for _, mname := range methodNames {
		connMethods[mname] = true
		if _, found := conn.DelegateMethods[mname]; !found {
			addingMethods = append(addingMethods, mname)
		}
	}
	// find old methods to delete
	for mname, _ := range conn.DelegateMethods {
		if _, found := connMethods[mname]; !found {
			deletingMethods = append(deletingMethods, mname)
		}
	}

	conn.DelegateMethods = connMethods
	maybeChanged := len(addingMethods) > 0 || len(deletingMethods) > 0

	// add methods
	for _, mname := range addingMethods {
		methodDelg := MethodDelegation{
			Conn: conn,
			Name: mname,
		}
		// bi direction map
		methodDelgArr, methodFound := self.delegateConnMap[mname]

		if methodFound {
			methodDelgArr = append(methodDelgArr, methodDelg)
		} else {
			var tmp []MethodDelegation
			methodDelgArr = append(tmp, methodDelg)
		}
		self.delegateConnMap[mname] = methodDelgArr
	}
	// delete methods
	for _, mname := range deletingMethods {
		methodDelgList, ok := self.delegateConnMap[mname]
		if !ok {
			continue
		}
		methodDelgList = DelegateRemoveConn(methodDelgList, conn)
		if len(methodDelgList) > 0 {
			self.delegateConnMap[mname] = methodDelgList
		} else {
			delete(self.delegateConnMap, mname)
		}
	}
	return maybeChanged
}

func (self *Router) probeMethodChange() {
	sig := self.getMethodsSig()
	if self.methodsSig != sig {
		// notify local methods change by broadcasting notification
		log.Debugf("local methods sig changed from %s to %s", self.methodsSig, sig)
		self.methodsSig = sig

		go self.NotifyStateChange()
	}
}

func (self *Router) leaveConn(conn *ConnT) {
	for method, _ := range conn.ServeMethods {
		methodDescList, ok := self.methodConnMap[method]
		if !ok {
			continue
		}
		methodDescList = RemoveConn(methodDescList, conn)
		if len(methodDescList) > 0 {
			self.methodConnMap[method] = methodDescList
		} else {
			delete(self.methodConnMap, method)
		}
	}
	conn.ServeMethods = make(map[string]MethodInfo)

	ct, ok := self.connMap[conn.ConnId]
	if ok {
		delete(self.connMap, conn.ConnId)
		close(ct.RecvChannel)
	}

	if conn.publicId != "" {
		if _, found := self.publicConnMap[conn.publicId]; found {
			delete(self.publicConnMap, conn.publicId)
			conn.publicId = ""
		}
	}

	// remove conn from fallbackConns
	if conn.AsFallback {
		var fbIndex = -1
		for i, c := range self.fallbackConns {
			if c.ConnId == conn.ConnId {
				//conn found in fallback conns
				fbIndex = i
				break
			}
		}
		if fbIndex >= 0 {
			self.fallbackConns = append(
				self.fallbackConns[:fbIndex],
				self.fallbackConns[fbIndex+1:]...)
		}
	}
	self.probeMethodChange()
}

func (self *Router) ListConns(method string, limit int) []CID {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	var arr []CID
	if descs, ok := self.methodConnMap[method]; ok && len(descs) > 0 {
		for _, desc := range descs {
			arr = append(arr, desc.Conn.ConnId)
			if len(arr) >= limit {
				break
			}
		}
	}
	return arr
}

func (self *Router) SelectConn(method string, targetConnId CID) (*ConnT, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	if targetConnId != ZeroCID {
		conn, found := self.connMap[targetConnId]
		return conn, found
	}

	if descs, ok := self.methodConnMap[method]; ok && len(descs) > 0 {
		index := rand.Intn(len(descs))
		return descs[index].Conn, true

	}

	if delgs, ok := self.delegateConnMap[method]; ok && len(delgs) > 0 {
		index := rand.Intn(len(delgs))
		return delgs[index].Conn, true
	}

	// fallback conns
	if len(self.fallbackConns) > 0 {
		index := rand.Intn(len(self.fallbackConns))
		return self.fallbackConns[index], true
	}
	return nil, false
}

func (self *Router) GetConnByPublicId(publicId string) (*ConnT, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	conn, found := self.publicConnMap[publicId]
	return conn, found
}

func (self *Router) SelectReceiver(method string) (MsgChannel, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	descs, ok := self.methodConnMap[method]
	if ok && len(descs) > 0 {
		// or random or round-robin
		conn := descs[0].Conn
		return conn.RecvChannel, true
	}
	return nil, false
}

func (self *Router) ClearTimeoutRequests() {
	now := time.Now()
	tmpMap := make(map[interface{}]PendingT)
	//var tmpMap map[stirng]PendingT

	for msgId, reqt := range self.pendingRequests {
		if now.After(reqt.Expire) {
			errMsg := jsonrpc.RPCErrorMessage(reqt.OrigMsgVec.Msg, 408, "request timeout", true)
			msgvec := MsgVec{Msg: errMsg}
			_ = self.deliverMessage(reqt.OrigMsgVec.FromConnId, msgvec)
		} else {
			tmpMap[msgId] = reqt
		}
	}
	self.pendingRequests = tmpMap
}

func (self *Router) TryClearPendingRequest(msgId interface{}) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	if _, found := self.pendingRequests[msgId]; found {
		go func() {
			// sleep for another 1 second
			time.Sleep(1 * time.Second)
			self.ClearPendingRequest(msgId)
		}()
	}
}

func (self *Router) ClearPendingRequest(msgId interface{}) {
	self.lock("ClearPendingRequest")
	defer self.unlock("ClearPendingRequest")

	if reqt, found := self.pendingRequests[msgId]; found {
		now := time.Now()
		if !now.After(reqt.Expire) {
			reqt.OrigMsgVec.Msg.Log().Errorf("Expire is not reached even during collecting routine")
		}
		errMsg := jsonrpc.RPCErrorMessage(reqt.OrigMsgVec.Msg, 408, "request timeout", true)
		msgvec := MsgVec{Msg: errMsg}
		_ = self.deliverMessage(reqt.OrigMsgVec.FromConnId, msgvec)

		// delete from self.pendingRequests
		delete(self.pendingRequests, msgId)
	}
}

func (self *Router) ClearPending(connId CID) {
	for msgId, reqt := range self.pendingRequests {
		if reqt.ToConnId == connId || reqt.OrigMsgVec.FromConnId == connId {
			self.deletePending(msgId)
		}
	}
}

func (self *Router) deletePending(msgId interface{}) {
	delete(self.pendingRequests, msgId)
}

func (self *Router) routeMessage(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	msg.Log().WithFields(log.Fields{"from": cmdMsg.MsgVec.FromConnId}).Debugf("route message")
	if msg.IsRequest() {
		return self.routeRequest(cmdMsg)
	} else if msg.IsNotify() {
		toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.MsgVec.TargetConnId)
		if found {
			return self.deliverMessage(
				toConn.ConnId, cmdMsg.MsgVec)
		}
	} else if msg.IsResultOrError() {
		return self.routeResultOrError(cmdMsg)
	}
	return nil
}

func (self *Router) routeRequest(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	fromConnId := cmdMsg.MsgVec.FromConnId
	toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.MsgVec.TargetConnId)
	if found {
		timeout := cmdMsg.Timeout
		if timeout <= 0 {
			timeout = DefaultRequestTimeout
		}

		//fmt.Printf("timeout %d\n", timeout)
		expireTime := time.Now().Add(timeout)
		reqMsg, ok := msg.(*jsonrpc.RequestMessage)
		misc.Assert(ok, "bad msg type other than request")

		msgId := reqMsg.Id

		if _, ok := self.pendingRequests[msgId]; ok {
			msgId = uuid.New().String()
			reqMsg.Log().Infof("msg id already exist, change a new one %s", msgId)
			reqMsg = reqMsg.Clone(msgId)
		}
		self.pendingRequests[msgId] = PendingT{
			OrigMsgVec: cmdMsg.MsgVec,
			Expire:     expireTime,
			ToConnId:   toConn.ConnId,
		}
		targetVec := cmdMsg.MsgVec
		targetVec.Msg = reqMsg
		go func() {
			time.Sleep(timeout)
			time.Sleep(1 * time.Second)
			//time.Sleep(int64(timeout.Seconds() + 1) * time.Second)
			self.TryClearPendingRequest(msgId)
		}()
		return self.deliverMessage(toConn.ConnId, targetVec)
	} else {
		errMsg := jsonrpc.RPCErrorMessage(msg, 404, "method not found", false)
		errMsg.SetTraceId(msg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg}
		return self.deliverMessage(fromConnId, errMsgVec)
	}

}
func (self *Router) routeResultOrError(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	//if msgId, ok := msg.MustId().(string); ok {
	msgId := msg.MustId()
	if true {
		if reqt, ok := self.pendingRequests[msgId]; ok {
			self.deletePending(msgId)

			if cmdMsg.MsgVec.FromConnId != reqt.ToConnId {
				msg.Log().Warnf("msg conn %d not from the delivered conn %d", cmdMsg.MsgVec.FromConnId, reqt.ToConnId)
			}
			origReq := reqt.OrigMsgVec.Msg
			if msg.TraceId() != origReq.TraceId() {
				msg.Log().Warnf("result trace is different from request %s", origReq.TraceId())
			}
			if resMsg, ok := msg.(*jsonrpc.ResultMessage); ok {
				newRes := jsonrpc.NewResultMessage(reqt.OrigMsgVec.Msg, resMsg.Result, nil)
				newVec := cmdMsg.MsgVec
				newVec.Msg = newRes
				return self.deliverMessage(reqt.OrigMsgVec.FromConnId, newVec)
			} else if errMsg, ok := msg.(*jsonrpc.ErrorMessage); ok {
				newErr := jsonrpc.NewErrorMessage(reqt.OrigMsgVec.Msg, errMsg.Error, nil)
				newVec := cmdMsg.MsgVec
				newVec.Msg = newErr
				return self.deliverMessage(reqt.OrigMsgVec.FromConnId, newVec)
			} else {
				msg.Log().Fatalf("msg is neither result nor error")
			}
		} else {
			msg.Log().Warn("fail to find request from this result/error")
		}
	} else {
		msg.Log().Warnf("message id is not string, %+v", msg.MustId())
	}
	return nil
}

func (self *Router) deliverMessage(connId CID, msgvec MsgVec) *ConnT {
	ct, ok := self.connMap[connId]
	if ok {
		ct.RecvChannel <- msgvec
		return ct
	}
	return nil
}

func (self *Router) Start(ctx context.Context) {
	//self.setupChannels()
	go self.Loop(ctx)
}

func (self *Router) Loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debugf("Router goroutine done")
			return
		case cmdServe, ok := <-self.ChServe:
			{
				if !ok {
					log.Warnf("ChServe channel not ok")
					return
				}
				conn, found := self.connMap[cmdServe.ConnId]
				if found {
					self.UpdateServeMethods(conn, cmdServe.Methods)
				} else {
					log.Infof("Conn %d not found for update serve methods", cmdServe.ConnId)
				}
			}

		case cmdDelg, ok := <-self.ChDelegate:
			{
				if !ok {
					log.Warnf("ChServe channel not ok")
					return
				}
				conn, found := self.connMap[cmdDelg.ConnId]
				if found {
					self.UpdateDelegateMethods(conn, cmdDelg.MethodNames)
				} else {
					log.Infof("Conn %d not found for update methods", cmdDelg.ConnId)
				}
			}

		case cmdMsg, ok := <-self.ChMsg:
			{
				if !ok {
					log.Warnf("ChMsg channel not ok")
					return
				}

				self.RouteMessage(cmdMsg)
			}
		}
	}
}

// commands
func (self *Router) RouteMessage(cmdMsg CmdMsg) *ConnT {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	return self.routeMessage(cmdMsg)
}

func (self *Router) Join(genUUID bool) *ConnT {
	conn := NewConn()
	self.joinConn(conn, genUUID)
	return conn
}

func (self *Router) JoinFallback() *ConnT {
	conn := NewConn()
	self.joinFallbackConn(conn)
	return conn
}

func (self *Router) joinConn(conn *ConnT, genUUID bool) {
	self.lock("JoinConn")
	defer self.unlock("JoinConn")
	self.connMap[conn.ConnId] = conn
	if genUUID {
		conn.publicId = uuid.New().String()
		self.publicConnMap[conn.publicId] = conn
	}
}

func (self *Router) joinFallbackConn(conn *ConnT) {
	self.lock("JoinConnFallback")
	defer self.unlock("JoinConnFallback")
	self.connMap[conn.ConnId] = conn
	conn.AsFallback = true
	self.fallbackConns = append(self.fallbackConns, conn)
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

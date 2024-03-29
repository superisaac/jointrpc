package rpcrouter

import (
	//"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	//"github.com/superisaac/jointrpc/datadir"
	//"github.com/superisaac/jsonz"
	"github.com/superisaac/jointrpc/misc"
	//misc "github.com/superisaac/jointrpc/misc"
	"math/rand"
	"sort"
	"strings"
	//"sync"
	//"time"
)

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

func NewRouter(factory *RouterFactory, name string) *Router {
	r := &Router{factory: factory, name: name}
	//r.routerLock = new(sync.RWMutex)
	//r.pendingLock = new(sync.RWMutex)
	r.methodConnMap = make(map[string]([]MethodDesc))
	r.delegateConnMap = make(map[string][]MethodDelegation)
	r.connMap = make(map[CID](*ConnT))
	r.methodsSig = ""

	r.setupChannels()
	return r
}

func (self *Router) setupChannels() {
	self.chRouteMsg = make(MsgChannel, 10000)
	self.chSelectConn = make(chan CmdSelectConn, 10000)
	self.ChJoin = make(chan CmdJoin, 10000)
	self.ChLeave = make(chan CmdLeave, 10000)
	self.ChMethods = make(chan CmdMethods, 10000)
	self.ChDelegates = make(chan CmdDelegates, 10000)
}

func (self Router) Name() string {
	return self.name
}

func (self Router) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"namespace": self.name,
	})
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

func (self Router) HasMethod(method string) bool {
	if _, ok := self.methodConnMap[method]; ok {
		return true
	} else if _, ok := self.delegateConnMap[method]; ok {
		return true
	}
	return false
}

func (self Router) GetDelegates() []string {
	//self.rlock("GetDelegates")
	//defer self.runlock("GetDelegates")

	var arr []string
	for name, _ := range self.delegateConnMap {
		arr = append(arr, name)
	}
	return arr
}
func (self Router) GetMethods() []MethodInfo {
	//self.rlock("GetMethods")
	//defer self.runlock("GetMethods")
	return self.getMethods()
}

func (self Router) GetMethodNames() []string {
	//self.rlock("GetMethodNames")
	//defer self.runlock("GetMethods")

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
		if len(descs) > 0 {
			minfos = append(minfos, descs[0].Info)
		}
		// for _, desc := range descs {
		// 	minfos = append(minfos, desc.Info)
		// }
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

func (self *Router) OnCmdMethods(cmdMethods CmdMethods) {
	//self.lock("OnCmdMethods")
	//defer self.unlock("OnCmdMethods")
	conn, found := self.connMap[cmdMethods.ConnId]
	if found {
		self.updateServeMethods(conn, cmdMethods.Methods)
	} else {
		self.Log().Infof("Conn %d not found for update serve methods", cmdMethods.ConnId)
	}
}

func (self *Router) OnCmdDelegates(cmdDelg CmdDelegates) {
	//self.lock("OnCmdDelegates")
	//defer self.unlock("OnCmdDelegates")
	conn, found := self.connMap[cmdDelg.ConnId]
	if found {
		self.updateDelegateMethods(conn, cmdDelg.MethodNames)
	} else {
		self.Log().Infof("Conn %d not found for update methods", cmdDelg.ConnId)
	}
}

func (self *Router) UpdateServeMethods(conn *ConnT, methods []MethodInfo) bool {
	//self.lock("CanServeMethods")
	//defer self.unlock("CanServeMethods")
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
	//self.lock("CanDelegateMethods")
	//defer self.unlock("CanDelegateMethods")
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

	if ct, ok := self.connMap[conn.ConnId]; ok {
		delete(self.connMap, conn.ConnId)
		ct.Destruct()
	}

	self.probeMethodChange()
}

func (self *Router) ListConns(method string, limit int) []CID {
	//self.rlock("ListConns")
	//defer self.runlock("ListConns")

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
	return self.selectConnection(method, targetConnId)
	// conn, found := self.selectConnection(method, targetConnId)
	// if !found && self.Name() != "*" {
	// 	return self.factory.CommonRouter().selectConnection(
	// 		method, targetConnId)
	// }
	// return conn, found
}

func (self *Router) selectConnection(method string, targetConnId CID) (*ConnT, bool) {
	//self.rlock("selectConnection")
	//defer self.runlock("selectConnection")

	if targetConnId != ZeroCID {
		conn, found := self.connMap[targetConnId]
		return conn, found
	}

	// 1st round, select a free connection randomly
	if descs, ok := self.methodConnMap[method]; ok && len(descs) > 0 {
		// choose some free conns
		for i := 0; i < 5; i++ {
			index := rand.Intn(len(descs))
			conn := descs[index].Conn
			if len(conn.MsgOutput()) <= 2 { // skip the choice if there are too many elements in channel buffer
				return conn, true
			}
		}

	}

	// 2nd round, select a random delegate connection
	if delgs, ok := self.delegateConnMap[method]; ok && len(delgs) > 0 {
		index := rand.Intn(len(delgs))
		return delgs[index].Conn, true
	}

	// 3rd round, select a random connection any way
	if descs, ok := self.methodConnMap[method]; ok && len(descs) > 0 {
		// if no free conns just choose the random one anyway
		index := rand.Intn(len(descs))
		return descs[index].Conn, true
	}

	return nil, false
}

func (self *Router) GetConn(connId CID) (*ConnT, bool) {
	//self.rlock("GetConn")
	//defer self.runlock("GetConn")
	conn, found := self.connMap[connId]
	return conn, found
}

func (self *Router) Join() *ConnT {
	conn := NewConn()
	self.joinConn(conn)
	return conn
}

func (self *Router) joinConn(conn *ConnT) {
	misc.Assert(!conn.Joined(), "conn already joined")
	//self.lock("JoinConn")
	//defer self.unlock("JoinConn")
	conn.Namespace = self.name
	conn.router = self
	self.connMap[conn.ConnId] = conn
}

// func (self *Router) lock(wrapper string) {
// 	//log.Printf("router want lock %s", wrapper)
// 	//self.routerLock.Lock()
// 	//log.Printf("router locked %s", wrapper)
// }
// func (self *Router) unlock(wrapper string) {
// 	//log.Printf("router want unlock %s", wrapper)
// 	//self.routerLock.Unlock()
// 	//log.Printf("router want unlocked %s", wrapper)
// }

// func (self *Router) rlock(wrapper string) {
// 	//log.Printf("router want lock %s", wrapper)
// 	//self.routerLock.RLock()
// 	//log.Printf("router locked %s", wrapper)
// }
// func (self *Router) runlock(wrapper string) {
// 	//log.Printf("router want unlock %s", wrapper)
// 	//self.routerLock.RUnlock()
// 	//log.Printf("router want unlocked %s", wrapper)
// }

// func (self *Router) lockPending(wrapper string) {
// 	//log.Printf("router pending want lock %s", wrapper)
// 	//self.pendingLock.Lock()
// 	//log.Printf("router pending locked %s", wrapper)
// }
// func (self *Router) unlockPending(wrapper string) {
// 	//log.Printf("router pending want unlock %s", wrapper)
// 	//self.pendingLock.Unlock()
// 	//log.Printf("router pending want unlocked %s", wrapper)
// }

// func (self *Router) rlockPending(wrapper string) {
// 	//log.Printf("router pending want lock %s", wrapper)
// 	//self.pendingLock.RLock()
// 	//log.Printf("router pending locked %s", wrapper)
// }
// func (self *Router) runlockPending(wrapper string) {
// 	//log.Printf("router pending want unlock %s", wrapper)
// 	//self.pendingLock.RUnlock()
// 	//log.Printf("router pending want unlocked %s", wrapper)
// }

func (self *Router) leave(conn *ConnT) {
	//self.lock("Leave")
	//defer self.unlock("Leave")

	self.leaveConn(conn)
}

func (self Router) Started() bool {
	return self.startCtx != nil
}

func (self *Router) EnsureStart(rootCtx context.Context) {
	if self.startCtx == nil {
		self.startCtx, self.cancelFunc = context.WithCancel(rootCtx)
		go self.loop()
	}
}

func (self *Router) Stop() {
	if self.Started() {
		self.cancelFunc()
		self.startCtx = nil
		self.cancelFunc = nil
	}
}

func (self *Router) loop() {
	for {
		select {
		case <-self.startCtx.Done():
			log.Debugf("router goroutine done")
			return
		case cmdMethods, ok := <-self.ChMethods:
			{
				if !ok {
					log.Warnf("ChMethods channel not ok")
					return
				}

				misc.Assert(cmdMethods.Namespace != "", "bad cmdMethods")
				self.OnCmdMethods(cmdMethods)
			}

		case cmdDelg, ok := <-self.ChDelegates:
			{
				if !ok {
					log.Warnf("ChServe channel not ok")
					return
				}
				misc.Assert(cmdDelg.Namespace != "", "bad cmdDelg namespace")
				log.Infof("cmd delegate %+v", cmdDelg)
				self.OnCmdDelegates(cmdDelg)
			}
		case cmdJoin, ok := <-self.ChJoin:
			{
				if !ok {
					log.Warnf("ChJoin channel not ok")
					return
				}
				log.Infof("join command from conn %d", cmdJoin.Conn.ConnId)
				self.joinConn(cmdJoin.Conn)
				if cmdJoin.ChRet != nil {
					cmdJoin.ChRet <- CmdRet{Ok: true}
				} else {
					log.Infof("conn %d does not take chRet to join", cmdJoin.Conn.ConnId)
				}
			}
		case cmdLeave, ok := <-self.ChLeave:
			{
				if !ok {
					log.Warnf("ChLeave channel not ok")
					return
				}
				log.Infof("leave command from conn %d", cmdLeave.Conn.ConnId)
				self.leave(cmdLeave.Conn)
				if cmdLeave.ChRet != nil {
					cmdLeave.ChRet <- CmdRet{Ok: true}
				} else {
					log.Infof("conn %d does not take chRet to leave", cmdLeave.Conn.ConnId)
				}
			}
		case cmdSelConn, ok := <-self.chSelectConn:
			{
				if !ok {
					log.Warnf("ChSelectConn channel closed")
					return
				}
				toConn, found := self.selectConnection(cmdSelConn.Method, cmdSelConn.ConnId)
				if found {
					cmdSelConn.ChRet <- RetSelectConn{
						Method: cmdSelConn.Method,
						ConnId: cmdSelConn.ConnId,
						Conn:   toConn,
						Found:  true,
					}
				} else {
					cmdSelConn.ChRet <- RetSelectConn{
						Method: cmdSelConn.Method,
						ConnId: cmdSelConn.ConnId,
						Conn:   nil,
						Found:  false,
					}
				}
			}
		case cmdMsg, ok := <-self.chRouteMsg:
			{
				if !ok {
					log.Warnf("ChMsg channel not ok")
					return
				}
				misc.Assert(cmdMsg.Namespace != "", "bad msg namespace")
				cmdMsg.Msg.Log().Debugf("size of ChRouteMsg is %d", len(self.chRouteMsg))
				self.relayMessage(cmdMsg)
			}
		}
	}
}

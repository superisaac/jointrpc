package tube

import (
	"sync"
	"time"
	"github.com/gorilla/websocket"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func NewRouter() *Router {
	return new(Router).Init()
}

func GetConnId(c *websocket.Conn) string {
	return c.UnderlyingConn().RemoteAddr().String()
}

func RemoveElement(slice []jsonrpc.CID, elems jsonrpc.CID) []jsonrpc.CID {
	for i := range slice {
		if slice[i] == elems {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (self *Router) Init() *Router {
	self.routerLock = new(sync.RWMutex)
	self.MethodConnMap = make(map[string]([]jsonrpc.CID))
	self.ConnMethodMap = make(map[jsonrpc.CID]([]string))
	self.ConnMap = make(map[jsonrpc.CID](ConnT))
	self.PendingMap = make(map[PendingKey]PendingValue)
	return self
}

func (self *Router) registerConn(connId jsonrpc.CID, conn ConnT) {
	self.ConnMap[connId] = conn
	// register connId as a service name
}

func (self *Router) RegisterService(connId jsonrpc.CID, method string) error {
	self.routerLock.Lock()
	defer self.routerLock.Unlock()

	// bi direction map
	cidArr, ok := self.MethodConnMap[method]
	if ok {
		cidArr = append(cidArr, connId)
	} else {
		var a []jsonrpc.CID
		cidArr = append(a, connId)
	}
	self.MethodConnMap[method] = cidArr

	snArr, ok := self.ConnMethodMap[connId]
	if ok {
		snArr = append(snArr, method)
	} else {
		var a []string
		snArr = append(a, method)
	}
	self.ConnMethodMap[connId] = snArr

	return nil
}

func (self *Router) UnRegisterService(connId jsonrpc.CID, method string) error {
	self.routerLock.Lock()
	defer self.routerLock.Unlock()

	methods, ok := self.ConnMethodMap[connId]
	if ok {
		var tmpMethods []string

		for _, sname := range methods {
			if sname != method {
				tmpMethods = append(tmpMethods, sname)
			}
		}
		if len(tmpMethods) > 0 {
			self.ConnMethodMap[connId] = tmpMethods
		} else {
			delete(self.ConnMethodMap, connId)
		}
	}

	connIds, ok := self.MethodConnMap[method]
	if ok {
		var tmpConnIds []jsonrpc.CID
		for _, cid := range connIds {
			if cid != connId {
				tmpConnIds = append(tmpConnIds, cid)
			}

			if len(tmpConnIds) > 0 {
				self.MethodConnMap[method] = tmpConnIds
			} else {
				delete(self.MethodConnMap, method)
			}
		}
	}


	ct, ok := self.ConnMap[connId]
	if ok {
		delete(self.ConnMap, connId)
		close(ct.RecvChannel())
	}
	return nil
}

func (self *Router) unregisterConn(connId jsonrpc.CID) {
	self.ClearPending(connId)
	self.routerLock.Lock()
	defer self.routerLock.Unlock()

	methods, ok := self.ConnMethodMap[connId]
	if ok {
		for _, method := range methods {
			connIds, ok := self.MethodConnMap[method]
			if !ok {
				continue
			}
			connIds = RemoveElement(connIds, connId)
			if len(connIds) > 0 {
				self.MethodConnMap[method] = connIds
			} else {
				delete(self.MethodConnMap, method)
			}
		}
		delete(self.ConnMethodMap, connId)
	}

	ct, ok := self.ConnMap[connId]
	if ok {
		delete(self.ConnMap, connId)
		close(ct.RecvChannel())
	}
}

func (self *Router) SelectConn(method string) (jsonrpc.CID, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	connIds, ok := self.MethodConnMap[method]
	if ok && len(connIds) > 0 {
		// or random or round-robin
		return connIds[0], true
	}
	return 0, false
}

func (self *Router) GetMethods(connId jsonrpc.CID) []string {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	return self.ConnMethodMap[connId]
}

func (self *Router) ClearTimeoutRequests() {
	now := time.Now()
	tmpMap := make(map[PendingKey]PendingValue)

	for pKey, pValue := range self.PendingMap {
		if now.After(pValue.Expire) {
			errMsg := jsonrpc.NewErrorMessage(pKey.MsgId, 408, "request timeout")
			_ = self.deliverMessage(pKey.ConnId, errMsg)
		} else {
			tmpMap[pKey] = pValue
		}
	}
	self.PendingMap = tmpMap
}

func (self *Router) ClearPending(connId jsonrpc.CID) {
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

func (self *Router) routeMessage(msg *jsonrpc.RPCMessage) *ConnT {
	fromConnId := msg.FromConnId
	if msg.IsRequest() {
		toConnId, found := self.SelectConn(msg.Method)
		if found {
			pKey := PendingKey{ConnId: fromConnId, MsgId: msg.Id}
			expireTime := time.Now().Add(DefaultRequestTimeout)
			pValue := PendingValue{ConnId: toConnId, Expire: expireTime}

			self.setPending(pKey, pValue)
			return self.deliverMessage(toConnId, msg)
		} /*else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "service not found")
			return self.deliverMessage(fromConnId, errMsg)
		}*/
	} else if msg.IsNotify() {
		toConnId, found := self.SelectConn(msg.Method)
		if found {
			return self.deliverMessage(toConnId, msg)
		}
		/* else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "service not found")
			return self.deliverMessage(fromConnId, errMsg)
		} */
	} else if msg.IsResultOrError() {
		for pKey, pValue := range self.PendingMap {
			if pKey.MsgId == msg.Id && pValue.ConnId == fromConnId {
				// delete key within a range loop is safe
				// refer to https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-golang-map-within-a-range-loop
				self.deletePending(pKey)
				return self.deliverMessage(pKey.ConnId, msg)
			}
		} // end of for
	}
	return nil
}

func (self *Router) broadcastNotify(notify *jsonrpc.RPCMessage) (int, error) {
	if !notify.IsNotify() {
		/*errMsg := jsonrpc.NewErrorMessage(notify.Id, 400, "only notify can be broadcasted")
		self.deliverMessage(notify.FromConnId, errMsg)
		return nil */
		return 0, ErrNotNotify
	}
	cntDeliver := 0
	for connId, conn := range self.ConnMap {
		if conn.CanBroadcast() { // == IntentLocal {
			self.deliverMessage(connId, notify)
			cntDeliver += 1
		}
	}
	return cntDeliver, nil
}

func (self *Router) deliverMessage(connId jsonrpc.CID, msg *jsonrpc.RPCMessage) *ConnT {
	ct, ok := self.ConnMap[connId]
	if ok {
		ct.RecvChannel() <- msg //(*msg)
		return &ct
	}
	return nil
}

//func (self *Router) Start() {
/*	for {
		select {
		case cmdOpen := <-self.ChJoin:
			//self.registerConn(cmdOpen.ConnId, cmdOpen.Channel, cmdOpen.Intent)
		case msg := <-self.ChMsg:
			self.routeMessage(msg)
		case notify := <-self.ChBroadcast:
			self.broadcastNotify(notify)
		case cmdClose := <-self.ChLeave:
			self.unregisterConn(jsonrpc.CID(cmdClose))
		}
	} */
//}

// commands
func (self *Router) RouteMessage(msg *jsonrpc.RPCMessage, fromConnId jsonrpc.CID) *ConnT {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	
	msg.FromConnId = fromConnId
	//self.ChMsg <- msg
	return self.routeMessage(msg)
}

func (self *Router) BroadcastNotify(notify *jsonrpc.RPCMessage, fromConnId jsonrpc.CID) (int, error) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	
	notify.FromConnId = fromConnId
	//self.ChBroadcast <- notify
	return self.broadcastNotify(notify)
}

/*func (self *Router) Join(connId jsonrpc.CID, ch MsgChannel, intent string) {
	conn := &ConnT{RecvChannel: ch, Intent: intent}
	//self.registerConn(cmdOpen.ConnId, conn)
	self.JoinConn(connId, conn)
}*/

func (self *Router) JoinConn(connId jsonrpc.CID, conn ConnT) {
	self.routerLock.Lock()
	defer self.routerLock.Unlock()
	self.registerConn(connId, conn)
}

func (self *Router) Leave(connId jsonrpc.CID) {
	//self.ChLeave <- LeaveCommand(connId)
	self.routerLock.Lock()
	defer self.routerLock.Unlock()
	self.unregisterConn(connId)
}

package tube

import (
	//"fmt"
	"context"
	"errors"
	"sync"
	"time"
	//	"github.com/gorilla/websocket"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func NewRouter() *Router {
	return new(Router).Init()
}

/*func GetConnId(c *websocket.Conn) string {
	return c.UnderlyingConn().RemoteAddr().String()
}
*/

func RemoveElement(slice []CID, elems CID) []CID {
	for i := range slice {
		if slice[i] == elems {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (self *Router) Init() *Router {
	self.routerLock = new(sync.RWMutex)
	self.MethodConnMap = make(map[string]([]CID))
	self.ConnMethodMap = make(map[CID]([]string))
	self.ConnMap = make(map[CID](IConn))
	self.PendingMap = make(map[PendingKey]PendingValue)
	return self
}

func (self Router) GetMethods(connId CID) []string {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()
	return self.ConnMethodMap[connId]
}

func (self Router) GetAllMethods() []string {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	methods := []string{}
	for method, _ := range self.MethodConnMap {
		methods = append(methods, method)
	}
	return methods
}

func (self *Router) registerConn(connId CID, conn IConn) {
	self.ConnMap[connId] = conn
	// register connId as a service name
}

func (self *Router) RegisterMethod(connId CID, method string) error {
	self.routerLock.Lock()
	defer self.routerLock.Unlock()

	// bi direction map
	cidArr, methodFound := self.MethodConnMap[method]
	if methodFound {
		cidArr = append(cidArr, connId)
	} else {
		var a []CID
		cidArr = append(a, connId)
	}
	self.MethodConnMap[method] = cidArr

	snArr, connFound := self.ConnMethodMap[connId]
	if connFound {
		snArr = append(snArr, method)
	} else {
		var a []string
		snArr = append(a, method)
	}
	self.ConnMethodMap[connId] = snArr

	return nil
}

func (self *Router) UnRegisterMethod(connId CID, method string) error {
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
		var tmpConnIds []CID
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

func (self *Router) unregisterConn(connId CID) {
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

func (self *Router) SelectConn(method string) (CID, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	connIds, ok := self.MethodConnMap[method]
	if ok && len(connIds) > 0 {
		// or random or round-robin
		return connIds[0], true
	}
	return 0, false
}

func (self *Router) SelectReceiver(method string) (MsgChannel, bool) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	connIds, ok := self.MethodConnMap[method]
	if ok && len(connIds) > 0 {
		// or random or round-robin
		connId := connIds[0]
		conn, found := self.ConnMap[connId]
		if found {
			return conn.RecvChannel(), found
		}
	}
	return nil, false
}

func (self *Router) ClearTimeoutRequests() {
	now := time.Now()
	tmpMap := make(map[PendingKey]PendingValue)

	for pKey, pValue := range self.PendingMap {
		if now.After(pValue.Expire) {
			errMsg := jsonrpc.NewErrorMessage(pKey.MsgId, 408, "request timeout", true)
			_ = self.deliverMessage(pKey.ConnId, errMsg)
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

func (self *Router) routeMessage(msg *jsonrpc.RPCMessage, fromConnId CID) *IConn {
	if msg.IsRequest() {
		toConnId, found := self.SelectConn(msg.Method)
		if found {
			pKey := PendingKey{ConnId: fromConnId, MsgId: msg.Id}
			expireTime := time.Now().Add(DefaultRequestTimeout)
			pValue := PendingValue{ConnId: toConnId, Expire: expireTime}

			self.setPending(pKey, pValue)
			return self.deliverMessage(toConnId, msg)
		} else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "service not found", false)
			return self.deliverMessage(fromConnId, errMsg)
		}
	} else if msg.IsNotify() {
		toConnId, found := self.SelectConn(msg.Method)
		if found {
			return self.deliverMessage(toConnId, msg)
			/*} else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "service not found", false)
			return self.deliverMessage(fromConnId, errMsg) */
		}
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

func (self *Router) deliverMessage(connId CID, msg *jsonrpc.RPCMessage) *IConn {
	ct, ok := self.ConnMap[connId]
	//fmt.Printf("deliver message %v\n", msg)
	if ok {
		recv_ch := ct.RecvChannel()
		recv_ch <- msg
		return &ct
	}
	return nil
}

func (self *Router) setupChannels() {
	self.ChMsg = make(chan CmdMsg, 100)
	self.ChJoin = make(chan CmdJoin, 100)
	self.ChLeave = make(chan CmdLeave, 100)
	self.ChRegister = make(chan CmdRegister, 100)
}

func (self *Router) Start(ctx context.Context) {
	self.setupChannels()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case cmd_join := <-self.ChJoin:
				self.Join(cmd_join.ConnId, cmd_join.RecvChannel)
			case cmd_leave := <-self.ChLeave:
				self.Leave(cmd_leave.ConnId)
			case cmd_register := <-self.ChRegister:
				self.RegisterMethod(cmd_register.ConnId, cmd_register.Method)
			case cmd_msg := <-self.ChMsg:
				self.RouteMessage(cmd_msg.Msg, cmd_msg.FromConnId)
				//case notify := <-self.ChBroadcast:
				//self.broadcastNotify(notify)
				//case cmdClose := <-self.ChLeave:
				//self.unregisterConn(CID(cmdClose))
			}
		}
	}()
}

// commands
func (self *Router) RouteMessage(msg *jsonrpc.RPCMessage, fromConnId CID) *IConn {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	//msg.FromConnId = fromConnId
	//self.ChMsg <- msg
	return self.routeMessage(msg, fromConnId)
}

func (self *Router) BroadcastNotify(notify *jsonrpc.RPCMessage, fromConnId CID) (int, error) {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	//notify.FromConnId = fromConnId
	//self.ChBroadcast <- notify
	return self.broadcastNotify(notify)
}

// Drop-in implementor of IConn
type SimpleConnT struct {
	recvChannel  MsgChannel
	canBroadcast bool
}

func (self SimpleConnT) RecvChannel() MsgChannel {
	return self.recvChannel
}

func (self SimpleConnT) CanBroadcast() bool {
	return self.canBroadcast
}

func (self *Router) Join(connId CID, ch MsgChannel) {
	conn := &SimpleConnT{recvChannel: ch, canBroadcast: true}
	self.JoinConn(connId, conn)
}

func (self *Router) JoinConn(connId CID, conn IConn) {
	self.routerLock.Lock()
	defer self.routerLock.Unlock()
	self.registerConn(connId, conn)
}

func (self *Router) Leave(connId CID) {
	//self.ChLeave <- LeaveCommand(connId)
	self.routerLock.Lock()
	defer self.routerLock.Unlock()

	self.unregisterConn(connId)
}

func leaveConn(conn_id CID) {
	Tube().Router.ChLeave <- CmdLeave{ConnId: conn_id}
}

func (self *Router) SingleCall(req_msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
	if !req_msg.IsRequest() && !req_msg.IsNotify() {
		return nil, errors.New("only request and notify message accepted")
	}
	if req_msg.IsRequest() {
		conn_id := NextCID()
		recv_ch := make(MsgChannel, 100)
		// router will take care of closing the receive channel
		//defer close(recv_ch)

		self.ChJoin <- CmdJoin{RecvChannel: recv_ch, ConnId: conn_id}
		defer leaveConn(conn_id)

		self.ChMsg <- CmdMsg{Msg: req_msg, FromConnId: conn_id}

		recvmsg := <-recv_ch
		return recvmsg, nil
	} else {
		self.ChMsg <- CmdMsg{Msg: req_msg, FromConnId: 0}
		return nil, nil
	}
}

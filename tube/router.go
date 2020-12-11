package tube

import (
	//"fmt"
	"context"
	"errors"
	"log"
	"sort"
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
	return self
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
	return self.Location == Location_Local
}

func (self Router) GetLocalMethods() []string {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	methods := []string{}
	for method, descs := range self.MethodConnMap {
		for _, desc := range descs {
			if desc.IsLocal() {
				methods = append(methods, method)
			}
		}
	}
	sort.Strings(methods)
	return methods
}

func (self *Router) RegisterLocalMethods(conn *ConnT, methods []string) error {
	return self.RegisterMethods(conn, methods, Location_Local)
}

func (self *Router) RegisterMethods(conn *ConnT, methods []string, location MethodLocation) error {
	self.lock("RegisterMethods")
	defer self.unlock("RegisterMethods")
	for _, method := range methods {
		err := self.registerMethod(conn, method, location)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Router) registerMethod(conn *ConnT, method string, location MethodLocation) error {

	_, found := conn.Methods[method]
	if found {
		// method already attach to this connection
		return nil
	}

	conn.Methods[method] = true

	methodDesc := MethodDesc{Conn: conn, Location: location}
	// bi direction map
	methodDescArr, methodFound := self.MethodConnMap[method]

	if methodFound {
		methodDescArr = append(methodDescArr, methodDesc)
	} else {
		var tmp []MethodDesc
		methodDescArr = append(tmp, methodDesc)
	}
	self.MethodConnMap[method] = methodDescArr
	return nil
}

func (self *Router) UnregisterMethods(conn *ConnT, methods []string) {
	self.lock("UnregisterMethod")
	defer self.unlock("UnregisterMethod")
	for _, method := range methods {
		self.unregisterMethod(conn, method)
	}
}

func (self *Router) unregisterMethod(conn *ConnT, method string) {
	_, found := conn.Methods[method]
	if !found {
		// method is not attached to this connection, just return
		return
	}

	delete(conn.Methods, method)

	methodDescList, ok := self.MethodConnMap[method]
	if ok {
		methodDescList = RemoveConn(methodDescList, conn)
		if len(methodDescList) > 0 {
			self.MethodConnMap[method] = methodDescList
		} else {
			delete(self.MethodConnMap, method)
		}
	}
}

func (self *Router) leaveConn(conn *ConnT) {
	for method, _ := range conn.Methods {
		methodDescList, ok := self.MethodConnMap[method]
		if !ok {
			continue
		}
		log.Printf("method desc pre remove %v", methodDescList)
		methodDescList = RemoveConn(methodDescList, conn)
		log.Printf("method desc post remove %v", methodDescList)
		if len(methodDescList) > 0 {
			self.MethodConnMap[method] = methodDescList
		} else {
			delete(self.MethodConnMap, method)
		}
	}
	conn.Methods = make(map[string]bool)

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

func (self *Router) routeMessage(msg *jsonrpc.RPCMessage, fromConnId CID) *ConnT {
	if msg.IsRequest() {
		toConn, found := self.SelectConn(msg.Method)
		if found {
			pKey := PendingKey{ConnId: fromConnId, MsgId: msg.Id}
			expireTime := time.Now().Add(DefaultRequestTimeout)
			pValue := PendingValue{ConnId: toConn.ConnId, Expire: expireTime}

			self.setPending(pKey, pValue)
			return self.deliverMessage(toConn.ConnId, msg)
		} else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "method not found", false)
			return self.deliverMessage(fromConnId, errMsg)
		}
	} else if msg.IsNotify() {
		toConn, found := self.SelectConn(msg.Method)
		if found {
			return self.deliverMessage(toConn.ConnId, msg)
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

func (self *Router) deliverMessage(connId CID, msg *jsonrpc.RPCMessage) *ConnT {
	ct, ok := self.ConnMap[connId]
	//fmt.Printf("deliver message %v\n", msg)
	if ok {
		recv_ch := ct.RecvChannel
		recv_ch <- msg
		return ct
	}
	return nil
}

func (self *Router) setupChannels() {
	self.ChMsg = make(chan CmdMsg, 100)
	self.ChLeave = make(chan CmdLeave, 100)
	self.ChReg = make(chan CmdReg, 100)
	self.ChUnreg = make(chan CmdUnreg, 100)
}

func (self *Router) Start(ctx context.Context) {
	self.setupChannels()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
				/*case cmd_join := <-self.ChJoin:
				self.Join(cmd_join.ConnId, cmd_join.RecvChannel) */
			case cmd_leave := <-self.ChLeave:
				conn, found := self.ConnMap[cmd_leave.ConnId]
				if found {
					self.Leave(conn)
				}
			case cmd_reg := <-self.ChReg:
				{
					conn, found := self.ConnMap[cmd_reg.ConnId]
					if found {
						self.RegisterMethods(conn, cmd_reg.Methods, cmd_reg.Location)
					}
				}
			case cmd_unreg := <-self.ChUnreg:
				{
					conn, found := self.ConnMap[cmd_unreg.ConnId]
					if found {
						self.UnregisterMethods(conn, cmd_unreg.Methods)
					}
				}
			case cmd_msg := <-self.ChMsg:
				{
					self.RouteMessage(cmd_msg.Msg, cmd_msg.FromConnId)
				}
			}
		}
	}()
}

// commands
func (self *Router) RouteMessage(msg *jsonrpc.RPCMessage, fromConnId CID) *ConnT {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	//msg.FromConnId = fromConnId
	//self.ChMsg <- msg
	return self.routeMessage(msg, fromConnId)
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

func (self *Router) SingleCall(req_msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
	if !req_msg.IsRequest() && !req_msg.IsNotify() {
		return nil, errors.New("only request and notify message accepted")
	}
	if req_msg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		self.ChMsg <- CmdMsg{Msg: req_msg, FromConnId: conn.ConnId}

		recvmsg := <-conn.RecvChannel
		return recvmsg, nil
	} else {
		self.ChMsg <- CmdMsg{Msg: req_msg, FromConnId: 0}
		return nil, nil
	}
}

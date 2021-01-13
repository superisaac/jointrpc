package bridge

import (
	//"fmt"
	"context"
	//"errors"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	"strings"
	//datadir "github.com/superisaac/jointrpc/datadir"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
)

// Bridge
func StartNewBridge(rootCtx context.Context, entries []client.ServerEntry) *Bridge {
	bridge := NewBridge(entries)
	bridge.Start(rootCtx)
	return bridge
}

func NewBridge(entries []client.ServerEntry) *Bridge {
	misc.Assert(len(entries) > 0, "empty entries")
	bridge := new(Bridge)
	bridge.serverEntries = entries
	bridge.edges = make(map[string]*Edge)
	bridge.ChState = make(chan CmdStateChange)
	return bridge
}

func (self *Bridge) Start(rootCtx context.Context) error {
	for _, entry := range self.serverEntries {
		edge := NewEdge(entry)
		self.edges[entry.ServerUrl] = edge
		go edge.Start(rootCtx, self)
	}

	mainCtx, mainCancel := context.WithCancel(rootCtx)
	defer mainCancel()
	for {
		select {
		case <-mainCtx.Done():
			// TODO: log
			return nil
		case stateChange, ok := <-self.ChState:
			if !ok {
				// TODO: log
				return nil
			}
			self.handleStateChange(stateChange)
		}
	}
	return nil
}

func (self *Bridge) requestReceived(msgvec rpcrouter.MsgVec, fromAddress string) (interface{}, error) {
	// stupid methods
	msg := msgvec.Msg
	if msg.IsRequest() {
		for sn, edge := range self.edges {
			if sn == fromAddress {
				continue
			}
			if edge.hasMethod(msg.MustMethod()) {
				resmsg, err := edge.remoteClient.CallMessage(context.Background(), msg)
				if err != nil {
					return nil, err
				}

				if resmsg.MustId() != msg.MustId() {
					log.Fatal("result has not the same id with origial request msg")
				}
				return resmsg.MustResult(), nil
			}
		}
		return nil, &jsonrpc.RPCError{404, "method not found", false}
	} else {
		log.Warnf("unexpected msg received %+v", msg)
		return nil, nil
	}
}

func (self *Bridge) handleStateChange(stateChange CmdStateChange) {
	if fromEdge, ok := self.edges[stateChange.serverUrl]; ok {
		self.exchangeDelegateMethods(fromEdge)
	} else {
		log.Warnf("fail to find edges %s", stateChange.serverUrl)
	}
}

func (self *Bridge) exchangeDelegateMethodsForEdge(aEdge *Edge) {
	uni := misc.NewStringUnifier()
	for _, edge := range self.edges {
		if edge == aEdge {
			continue
		}
		for mname, _ := range edge.methodNames {
			uni.Add(mname)
		}
	}
	aEdge.UpdateDelegateMethods(uni.Result())
}

func (self *Bridge) exchangeDelegateMethods(fromEdge *Edge) {
	for _, edge := range self.edges {
		if edge == fromEdge {
			continue
		}
		self.exchangeDelegateMethodsForEdge(edge)
	}
}

// edge methods
func NewEdge(entry client.ServerEntry) *Edge {
	return &Edge{
		remoteClient: client.NewRPCClient(entry),
		methodNames:  make(misc.StringSet),
	}
}

func (self Edge) hasMethod(methodName string) bool {
	_, ok := self.methodNames[methodName]
	return ok
}

func (self *Edge) onStateChange(state *rpcrouter.TubeState) {
	//fmt.Printf("state change %v\b", state)
	// update edge records
	methodNames := make(misc.StringSet)
	if state != nil {
		for _, minfo := range state.Methods {
			if strings.HasPrefix(minfo.Name, ".") {
				continue
			}
			if _, ok := methodNames[minfo.Name]; ok {
				continue
			}
			methodNames[minfo.Name] = true
		}
	}
	self.methodNames = methodNames
}

func (self *Edge) Start(rootCtx context.Context, bridge *Bridge) error {
	entry := self.remoteClient.ServerEntry()
	self.remoteClient.OnStateChange(func(state *rpcrouter.TubeState) {
		self.onStateChange(state)
		bridge.ChState <- CmdStateChange{
			serverUrl: entry.ServerUrl,
			state:     state,
		}
	})
	self.remoteClient.OnConnected(func() {
		log.Debugf("bridge client connected %s\n", entry.ServerUrl)
		if self.delegateMethods != nil || len(self.delegateMethods) > 0 {
			self.remoteClient.UpdateDelegateMethods(self.delegateMethods)
		}
	})
	self.remoteClient.OnConnectionLost(func() {
		log.Debugf("bridge client on connection lost %s\n", entry.ServerUrl)
		self.onStateChange(nil)
		bridge.ChState <- CmdStateChange{
			serverUrl: entry.ServerUrl,
			state:     nil,
		}
	})

	self.remoteClient.OnDefault(func(req *handler.RPCRequest, method string, params []interface{}) (interface{}, error) {
		return bridge.requestReceived(req.MsgVec, entry.ServerUrl)
	})

	err := self.remoteClient.Connect()
	if err != nil {
		return err
	}
	// TODO: concurrent
	return self.remoteClient.Handle(rootCtx)
}

func (self *Edge) UpdateDelegateMethods(methods []string) {
	//log.Infof("delegate %v", methods)
	self.delegateMethods = methods
	if self.remoteClient.Connected() {
		//log.Infof("client %s update methods %+v", self.remoteClient.ServerEntry().ServerUrl, self.delegateMethods)
		self.remoteClient.UpdateDelegateMethods(self.delegateMethods)
	}
}

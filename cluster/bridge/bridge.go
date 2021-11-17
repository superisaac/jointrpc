package bridge

import (
	"context"
	//"fmt"
	//"errors"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	//"strings"
	//datadir "github.com/superisaac/jointrpc/datadir"
	"github.com/superisaac/jointrpc/dispatch"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	jsonrpc "github.com/superisaac/jsonrpc"
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

func (self *Bridge) requestReceived(cmdMsg rpcrouter.CmdMsg, fromAddress string) (interface{}, error) {
	msg := cmdMsg.Msg
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
	//fmt.Printf("handle state change\n")
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
	aEdge.DeclareDelegateMethods(uni.Result())
}

func (self *Bridge) exchangeDelegateMethods(fromEdge *Edge) {
	log.Infof("exchange delegate methods from edge %s", fromEdge.remoteClient)
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

func (self *Edge) onStateChange(state *rpcrouter.ServerState) {
	// update edge records
	methodNames := make(misc.StringSet)
	if state != nil {
		for _, minfo := range state.Methods {
			if !jsonrpc.IsPublicMethod(minfo.Name) {
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
	disp := dispatch.NewDispatcher()
	stateListener := dispatch.NewStateListener()

	entry := self.remoteClient.ServerEntry()
	self.remoteClient.OnConnected(func() {
		log.Debugf("bridge client connected %s\n", entry.ServerUrl)
		if self.delegateMethods != nil || len(self.delegateMethods) > 0 {
			ctx, cancel := context.WithCancel(rootCtx)
			defer cancel()
			self.remoteClient.DeclareDelegates(ctx, self.delegateMethods)
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

	stateListener.OnStateChange(func(state *rpcrouter.ServerState) {
		//fmt.Printf("on state change %+v\n", state)
		self.onStateChange(state)
		bridge.ChState <- CmdStateChange{
			serverUrl: entry.ServerUrl,
			state:     state,
		}
	})

	disp.OnDefault(func(req *dispatch.RPCRequest, method string, params []interface{}) (interface{}, error) {
		return bridge.requestReceived(req.CmdMsg, entry.ServerUrl)
	})
	client.OnStateChanged(disp, stateListener)

	self.remoteClient.OnAuthorized(func() {
		req := self.remoteClient.NewWatchStateRequest()
		self.remoteClient.LiveCall(rootCtx, req,
			func(res jsonrpc.IMessage) {
				log.Infof("authorized, watch state")
			})
	})

	err := self.remoteClient.Connect()
	if err != nil {
		return err
	}
	// TODO: concurrent
	//go self.remoteClient.SubscribeState(rootCtx, stateListener)

	return self.remoteClient.Live(rootCtx, disp)
}

func (self *Edge) DeclareDelegateMethods(methods []string) {
	log.Infof("delegates %v", methods)
	self.delegateMethods = methods
	if self.remoteClient.Connected() {
		//log.Infof("client %s update methods %+v", self.remoteClient.ServerEntry().ServerUrl, self.delegateMethods)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		self.remoteClient.DeclareDelegates(ctx, self.delegateMethods)
	}
}

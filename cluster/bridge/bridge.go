package bridge

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	"strings"
	//datadir "github.com/superisaac/jointrpc/datadir"
	"github.com/superisaac/jointrpc/joint"
	handler "github.com/superisaac/jointrpc/joint/handler"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
)

// edge methods
func NewEdge() *Edge {
	return &Edge{
		methodNames: make(misc.StringSet),
	}
}

func (self Edge) hasMethod(methodName string) bool {
	_, ok := self.methodNames[methodName]
	return ok
}

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

func (self *Bridge) connectRemote(rootCtx context.Context, entry client.ServerEntry) error {
	if _, ok := self.edges[entry.Address]; ok {
		//log.Warnf("remote client already exist %s", self.remoteClient)
		panic(errors.New("client already exists"))
	}
	c := client.NewRPCClient(entry)

	err := c.Connect()
	if err != nil {
		return err
	}
	edge := NewEdge()
	edge.remoteClient = c
	self.edges[entry.Address] = edge

	c.OnStateChange(func(state *joint.TubeState) {
		self.ChState <- CmdStateChange{
			serverAddress: entry.Address,
			state:         state,
		}
	})
	c.OnDefault(func(req *handler.RPCRequest, method string, params []interface{}) (interface{}, error) {
		return self.messageReceived(req.MsgVec.Msg, entry.Address)
	})
	// TODO: concurrent
	c.Handle(rootCtx)
	return nil
}

func (self *Bridge) Start(rootCtx context.Context) error {
	for _, entry := range self.serverEntries {
		go self.connectRemote(rootCtx, entry)
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

func (self *Bridge) messageReceived(msg jsonrpc.IMessage, fromAddress string) (interface{}, error) {
	// stupid methods
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
	if edge, ok := self.edges[stateChange.serverAddress]; ok {
		// update edge records
		methodNames := make(misc.StringSet)
		for _, minfo := range stateChange.state.Methods {
			if strings.HasPrefix(minfo.Name, ".") {
				continue
			}
			if _, ok := methodNames[minfo.Name]; ok {
				continue
			}
			methodNames[minfo.Name] = true
		}
		edge.methodNames = methodNames

		self.exchangeDelegateMethods(stateChange.serverAddress)
	} else {
		log.Warnf("fail to find edges %s", stateChange.serverAddress)
	}
}

func (self *Bridge) exchangeDelegateMethodsForEdge(aEdge *Edge, fromAddress string) {
	uni := misc.NewStringUnifier()
	for _, edge := range self.edges {
		if edge == aEdge {
			continue
		}
		for mname, _ := range edge.methodNames {
			uni.Add(mname)
		}
	}
	aEdge.remoteClient.UpdateDelegateMethods(uni.Result())
}

func (self *Bridge) exchangeDelegateMethods(fromAddress string) {
	for sname, edge := range self.edges {
		if sname == fromAddress {
			continue
		}
		self.exchangeDelegateMethodsForEdge(edge, fromAddress)
	}
}

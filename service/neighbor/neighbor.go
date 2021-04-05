package neighbor

import (
	"context"
	"errors"
	//"fmt"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//datadir "github.com/superisaac/jointrpc/datadir"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"strings"
)

// edge methods
func NewEdge() *Edge {
	return &Edge{
		methodNames: make(misc.StringSet),
		dlgMethods:  make([]rpcrouter.MethodInfo, 0),
	}
}

func (self Edge) hasMethod(methodName string) bool {
	_, ok := self.methodNames[methodName]
	return ok
}

func NewNeighborService() *NeighborService {
	return new(NeighborService)
}

func (self *NeighborService) Init(rootCtx context.Context) {
	router := rpcrouter.RouterFromContext(rootCtx)
	cfg := router.Config

	var entries []client.ServerEntry
	for _, peer := range cfg.Neighbor.Peers {
		entries = append(entries, client.ServerEntry{
			ServerUrl: peer.ServerUrl,
			CertFile:  peer.CertFile,
		})
	}

	self.router = router
	self.InitDispatcher()
	self.serverEntries = entries
	self.edges = make(map[string]*Edge)
	self.ChState = make(chan CmdStateChange)
}

func (self NeighborService) Name() string {
	return "neighbor"
}

func (self NeighborService) CanRun(rootCtx context.Context) bool {
	router := rpcrouter.RouterFromContext(rootCtx)
	return len(router.Config.Neighbor.Peers) > 0
}

func (self *NeighborService) connectRemote(rootCtx context.Context, entry client.ServerEntry) error {
	if _, ok := self.edges[entry.ServerUrl]; ok {
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
	self.edges[entry.ServerUrl] = edge

	c.OnStateChange(func(state *rpcrouter.ServerState) {
		self.ChState <- CmdStateChange{
			ServerUrl: entry.ServerUrl,
			State:     state,
		}
	})
	c.Handle(rootCtx)
	return nil
}

func (self *NeighborService) Start(rootCtx context.Context) error {
	self.Init(rootCtx)

	for _, entry := range self.serverEntries {
		go self.connectRemote(rootCtx, entry)
	}

	// join connection
	self.conn = self.router.Join(false)
	defer func() {
		self.router.Leave(self.conn)
		self.conn = nil
	}()

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
		case msgvec, ok := <-self.conn.RecvChannel:
			if !ok {
				// TODO: log
				return nil
			}
			err := self.requestReceived(msgvec)
			if err != nil {
				return err
			}
		case resmsg, ok := <-self.ChResult:
			if !ok {
				// TODO: log
				return nil
			}
			self.router.DeliverResultOrError(
				rpcrouter.MsgVec{
					Msg:        resmsg,
					FromConnId: self.conn.ConnId,
				})
		}
	}
	return nil
}

func (self *NeighborService) requestReceived(msgvec rpcrouter.MsgVec) error {
	msg := msgvec.Msg
	// stupid methods
	if msg.IsRequest() {
		for _, edge := range self.edges {
			if edge.hasMethod(msg.MustMethod()) {
				resmsg, err := edge.remoteClient.CallMessage(
					context.Background(),
					msg)
				if err != nil {
					return err
				}
				misc.AssertEqual(resmsg.TraceId(), msgvec.Msg.TraceId(), "")

				if resmsg.MustId() != msg.MustId() {
					log.Fatal("result has not the same id with origial request msg")
				}

				self.ReturnResultMessage(resmsg)
				return nil
			}
		}
	} else {
		log.Warnf("unexpected msg received %+v", msg)
	}
	return nil

}

func (self *NeighborService) handleStateChange(stateChange CmdStateChange) {
	if edge, ok := self.edges[stateChange.ServerUrl]; ok {
		var newMethods []rpcrouter.MethodInfo
		methodNames := make(misc.StringSet)
		for _, minfo := range stateChange.State.Methods {
			if !jsonrpc.IsPublicMethod(minfo.Name) {
				continue
			}
			if _, ok := methodNames[minfo.Name]; ok {
				continue
			}
			newMInfo := minfo
			newMethods = append(newMethods, newMInfo)
			methodNames[minfo.Name] = true
		}

		edge.dlgMethods = newMethods
		edge.methodNames = methodNames
		self.tryUpdateMethods()
	} else {
		log.Warnf("fail to find edges %s", stateChange.ServerUrl)
	}
}

func (self *NeighborService) tryUpdateMethods() {
	uni := misc.NewStringUnifier()
	for _, edge := range self.edges {
		for _, minfo := range edge.dlgMethods {
			uni.Add(minfo.Name)
		}
	}
	// calculate the sig to avoid dup submit
	methodNames := uni.Result()
	sig := strings.Join(methodNames, ",")
	if sig != self.methodSig {
		self.methodSig = sig
		cmdDelegate := rpcrouter.CmdDelegate{
			ConnId:      self.conn.ConnId,
			MethodNames: methodNames,
		}
		self.router.ChDelegate <- cmdDelegate
	}
}

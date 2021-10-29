package neighbor

import (
	"context"
	"github.com/pkg/errors"
	//"fmt"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	datadir "github.com/superisaac/jointrpc/datadir"
	"github.com/superisaac/jointrpc/dispatch"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
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

// NeighborPort
func NewNeighborPort(namespace string, nbrCfg datadir.NeighborConfig) *NeighborPort {
	port := new(NeighborPort)
	var entries []client.ServerEntry
	for _, peer := range nbrCfg.Peers {
		entries = append(
			entries,
			client.ServerEntry{
				ServerUrl: peer.ServerUrl,
				CertFile:  peer.CertFile,
			})
	}
	port.namespace = namespace
	port.serverEntries = entries
	port.edges = make(map[string]*Edge)
	port.dispatcher = dispatch.NewDispatcher()
	port.chResult = make(chan dispatch.ResultT, 1000)
	port.ChState = make(chan CmdStateChange)
	return port
}

func (self *NeighborPort) connectRemote(rootCtx context.Context, entry client.ServerEntry) error {
	if _, ok := self.edges[entry.ServerUrl]; ok {
		//log.Warnf("remote client already exist %s", self.remoteClient)
		panic(errors.New("client already exists"))
	}
	c := client.NewRPCClient(entry)
	stateListener := dispatch.NewStateListener()

	err := c.Connect()
	if err != nil {
		return err
	}
	edge := NewEdge()
	edge.remoteClient = c
	edge.stateListener = stateListener

	self.edges[entry.ServerUrl] = edge

	stateListener.OnStateChange(func(state *rpcrouter.ServerState) {
		self.ChState <- CmdStateChange{
			ServerUrl: entry.ServerUrl,
			State:     state,
		}
	})
	//c.SubscribeState(rootCtx, stateListener)

	disp := dispatch.NewDispatcher()
	client.OnStateChanged(disp, stateListener)
	c.Live(rootCtx, disp)

	return nil
}

func (self *NeighborPort) Start(rootCtx context.Context) {
	factory := rpcrouter.RouterFactoryFromContext(rootCtx)
	router := factory.Get(self.namespace)
	for _, entry := range self.serverEntries {
		go self.connectRemote(rootCtx, entry)
	}

	// join connection
	//self.conn = router.Join()
	self.conn = rpcrouter.NewConn()
	router.ChJoin <- rpcrouter.CmdJoin{Conn: self.conn}

	defer func() {
		//router.Leave(self.conn)
		router.ChLeave <- rpcrouter.CmdLeave{Conn: self.conn}
		self.conn = nil
	}()

	mainCtx, mainCancel := context.WithCancel(rootCtx)
	defer mainCancel()
	for {
		select {
		case <-mainCtx.Done():
			// TODO: log
			return
		case stateChange, ok := <-self.ChState:
			if !ok {
				// TODO: log
				return
			}
			self.handleStateChange(factory, stateChange)
		case msgvec, ok := <-self.conn.RecvChannel:
			if !ok {
				// TODO: log
				return
			}
			err := self.requestReceived(msgvec)
			if err != nil {
				panic(err)
			}
		case result, ok := <-self.chResult:
			if !ok {
				// TODO: log
				return
			}
			//router.DeliverResultOrError(
			router.ChMsg <- rpcrouter.CmdMsg{
				MsgVec: rpcrouter.MsgVec{
					Namespace:  router.Name(),
					Msg:        result.ResMsg,
					FromConnId: self.conn.ConnId,
				},
			}
		}
	}
	return
}

func (self *NeighborPort) requestReceived(msgvec rpcrouter.MsgVec) error {
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

				self.dispatcher.ReturnResultMessage(resmsg, msgvec, self.chResult)
				return nil
			}
		}
	} else {
		log.Warnf("unexpected msg received %+v", msg)
	}
	return nil
}

func (self *NeighborPort) handleStateChange(factory *rpcrouter.RouterFactory, stateChange CmdStateChange) {
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
		self.tryUpdateMethods(factory)
	} else {
		log.Warnf("fail to find edges %s", stateChange.ServerUrl)
	}
}

func (self *NeighborPort) tryUpdateMethods(factory *rpcrouter.RouterFactory) {
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
		cmdDelegates := rpcrouter.CmdDelegates{
			Namespace:   self.namespace,
			ConnId:      self.conn.ConnId,
			MethodNames: methodNames,
		}
		factory.Get(self.namespace).ChDelegates <- cmdDelegates
	}
}

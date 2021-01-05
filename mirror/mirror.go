package mirror

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	datadir "github.com/superisaac/jointrpc/datadir"
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

// Mirror
func StartMirrorsForPeers(rootCtx context.Context) {
	cfg := datadir.ConfigFromContext(rootCtx)
	if len(cfg.Cluster.StaticPeers) > 0 {
		// generate server entry from peers
		var serverEntries []client.ServerEntry
		for _, peer := range cfg.Cluster.StaticPeers {
			serverEntries = append(serverEntries, client.ServerEntry{
				Address:  peer.Address,
				CertFile: peer.CertFile,
			})
		}
		go StartNewMirror(rootCtx, serverEntries)
	}
}

func StartNewMirror(rootCtx context.Context, entries []client.ServerEntry) {
	router := rpcrouter.RouterFromContext(rootCtx)
	mirror := NewMirror(entries, router)
	mirror.Start(rootCtx)
}

func NewMirror(entries []client.ServerEntry, router *rpcrouter.Router) *Mirror {
	mirror := new(Mirror)
	mirror.router = router
	mirror.InitHandlerManager()
	mirror.serverEntries = entries
	mirror.edges = make(map[string]*Edge)
	mirror.ChState = make(chan CmdStateChange)
	return mirror
}

func (self *Mirror) connectRemote(rootCtx context.Context, entry client.ServerEntry) error {
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

	c.OnStateChange(func(state *rpcrouter.TubeState) {
		self.ChState <- CmdStateChange{
			ServerAddress: entry.Address,
			State:         state,
		}
	})
	c.Handle(rootCtx)
	return nil
}

func (self *Mirror) Start(rootCtx context.Context) error {
	for _, entry := range self.serverEntries {
		go self.connectRemote(rootCtx, entry)
	}

	// join connection
	self.conn = self.router.Join()
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
			err := self.messageReceived(msgvec)
			if err != nil {
				return err
			}
		case resmsg, ok := <-self.ChResultMsg:
			if !ok {
				// TODO: log
				return nil
			}
			self.router.ChMsg <- rpcrouter.CmdMsg{
				MsgVec: rpcrouter.MsgVec{
					Msg:        resmsg,
					FromConnId: self.conn.ConnId,
				},
			}
		}
	}
	return nil
}

func (self *Mirror) messageReceived(msgvec rpcrouter.MsgVec) error {
	msg := msgvec.Msg
	// stupid methods
	if msg.IsRequest() {
		for _, edge := range self.edges {
			if edge.hasMethod(msg.MustMethod()) {
				resmsg, err := edge.remoteClient.CallMessage(context.Background(), msg)
				if err != nil {
					return err
				}
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

func (self *Mirror) handleStateChange(stateChange CmdStateChange) {
	if edge, ok := self.edges[stateChange.ServerAddress]; ok {
		var newMethods []rpcrouter.MethodInfo
		methodNames := make(misc.StringSet)
		for _, minfo := range stateChange.State.Methods {
			if strings.HasPrefix(minfo.Name, ".") {
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
		log.Warnf("fail to find edges %s", stateChange.ServerAddress)
	}
}

func (self *Mirror) tryUpdateMethods() {
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

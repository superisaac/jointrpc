package bridge

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/rpctube/client"
	datadir "github.com/superisaac/rpctube/datadir"
	tube "github.com/superisaac/rpctube/tube"
	"sort"
	"strings"
)

// edge methods
func NewEdge() *Edge {
	return &Edge{
		methodNames: make(StringSet),
		dlgMethods:  make([]tube.MethodInfo, 0),
	}
}

func (self Edge) hasMethod(methodName string) bool {
	_, ok := self.methodNames[methodName]
	return ok
}

// Bridge
func StartBridgesForPeers(rootCtx context.Context) {
	cfg := datadir.GetConfig()
	if len(cfg.Cluster.StaticPeers) > 0 {
		// generate server entry from peers
		var serverEntries []client.ServerEntry
		for _, peer := range cfg.Cluster.StaticPeers {
			serverEntries = append(serverEntries, client.ServerEntry{
				Address:  peer.Address,
				CertFile: peer.CertFile,
			})
		}
		go StartNewBridge(rootCtx, serverEntries)
	}
}

func StartNewBridge(rootCtx context.Context, entries []client.ServerEntry) {
	bridge := NewBridge(entries)
	bridge.Start(rootCtx)
}

func NewBridge(entries []client.ServerEntry) *Bridge {
	bridge := new(Bridge)
	bridge.InitHandlerManager()
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

	c.OnStateChange(func(state *tube.TubeState) {
		self.ChState <- CmdStateChange{
			ServerAddress: entry.Address,
			State:         state,
		}
	})
	c.Handle(rootCtx)
	return nil
}

func (self *Bridge) Start(rootCtx context.Context) error {
	for _, entry := range self.serverEntries {
		ctx := context.WithValue(
			rootCtx, "connectTo", entry.Address)
		go self.connectRemote(ctx, entry)
	}

	// join connection
	router := tube.Tube().Router
	self.conn = router.Join()
	defer func() {
		router.Leave(self.conn)
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
			router.ChMsg <- tube.CmdMsg{
				MsgVec: tube.MsgVec{
					Msg:        resmsg,
					FromConnId: self.conn.ConnId,
				},
			}
		}
	}
	return nil
}

func (self *Bridge) messageReceived(msgvec tube.MsgVec) error {
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

func (self *Bridge) handleStateChange(stateChange CmdStateChange) {
	if edge, ok := self.edges[stateChange.ServerAddress]; ok {
		var newMethods []tube.MethodInfo
		methodNames := make(StringSet)
		for _, minfo := range stateChange.State.Methods {
			if strings.HasPrefix(minfo.Name, ".") {
				continue
			}
			if _, ok := methodNames[minfo.Name]; ok {
				continue
			}
			newMInfo := minfo
			newMInfo.Delegated = true
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

func (self *Bridge) tryUpdateMethods() {
	dupChecker := make(StringSet)
	var methods []tube.MethodInfo
	for _, edge := range self.edges {
		for _, minfo := range edge.dlgMethods {
			if _, ok := dupChecker[minfo.Name]; ok {
				continue
			}
			dupChecker[minfo.Name] = true
			methods = append(methods)
		}
	}
	sort.Slice(methods, func(i, j int) bool { return methods[i].Name < methods[j].Name })
	// calculate the sig to avoid dup submit
	var sigArr []string
	for _, minfo := range methods {
		sigArr = append(sigArr, minfo.Name)
	}
	sig := strings.Join(sigArr, ",")
	if sig != self.methodSig {
		self.methodSig = sig
		cmdUpdate := tube.CmdUpdate{
			ConnId:  self.conn.ConnId,
			Methods: methods,
		}
		tube.Tube().Router.ChUpdate <- cmdUpdate
	}
}

package neighbor

import (
	"context"
	//"errors"
	//"fmt"
	//log "github.com/sirupsen/logrus"
	//client "github.com/superisaac/jointrpc/client"
	//"github.com/superisaac/jsonz"
	//datadir "github.com/superisaac/jointrpc/datadir"
	//"github.com/superisaac/jointrpc/dispatch"
	//misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	//"strings"
)

func NewNeighborService() *NeighborService {
	return new(NeighborService)
}

func (self *NeighborService) Init(rootCtx context.Context) {
	factory := rpcrouter.RouterFactoryFromContext(rootCtx)
	cfg := factory.Config

	self.ports = make(map[string]*NeighborPort)

	for namespace, nbrCfg := range cfg.Neighbors {
		port := NewNeighborPort(namespace, nbrCfg)
		self.ports[namespace] = port
	}

	//self.router = router
}

func (self NeighborService) Name() string {
	return "neighbor"
}

func (self NeighborService) CanRun(rootCtx context.Context) bool {
	factory := rpcrouter.RouterFactoryFromContext(rootCtx)
	for _, nbrCfg := range factory.Config.Neighbors {
		if len(nbrCfg.Peers) > 0 {
			return true
		}
	}
	return false
}

func (self *NeighborService) Start(rootCtx context.Context) error {
	self.Init(rootCtx)

	for _, port := range self.ports {
		go port.Start(rootCtx)
	}
	return nil
}

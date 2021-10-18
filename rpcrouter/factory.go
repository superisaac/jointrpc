package rpcrouter

import (
	//"fmt"
	"context"
	//"errors"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/datadir"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"time"
)

func NewRouterFactory(name string) *RouterFactory {

	factory := &RouterFactory{name: name}
	factory.Config = datadir.NewConfig()
	factory.routerMap = make(map[string](*Router))
	factory.setupChannels()
	return factory
}

func RouterFactoryFromContext(ctx context.Context) *RouterFactory {
	if v := ctx.Value("routerfactory"); v != nil {
		if factory, ok := v.(*RouterFactory); ok {
			return factory
		}
		panic("context value router is not a router instance")
	}
	panic("context does not have router")
}

func (self *RouterFactory) Get(name string) *Router {
	misc.Assert(name != "", "factory got empty namespace")
	if r, ok := self.routerMap[name]; ok {
		return r
	} else {
		t := NewRouter(self, name)
		self.routerMap[name] = t
		return t
	}
}

func (self RouterFactory) RouterNames() []string {
	names := make([]string, 0)
	for namespace, _ := range self.routerMap {
		names = append(names, namespace)
	}
	return names
}

func (self RouterFactory) GetOrNil(name string) *Router {
	misc.Assert(name != "", "factory got empty namespace")
	if r, ok := self.routerMap[name]; ok {
		return r
	} else {
		return nil
	}
}

func (self *RouterFactory) CommonRouter() *Router {
	return self.Get("*")
}

func (self *RouterFactory) DefaultRouter() *Router {
	return self.Get("default")
}

func (self *RouterFactory) setupChannels() {
	self.chMsg = make(chan CmdMsg, 1000)
	self.ChMethods = make(chan CmdMethods, 1000)
	self.ChDelegates = make(chan CmdDelegates, 1000)
}

func (self RouterFactory) Name() string {
	return self.name
}

func (self *RouterFactory) Start(ctx context.Context) {
	self.Loop(ctx)
}

func (self *RouterFactory) Loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debugf("RouterFactory goroutine done")
			return
		case cmdMethods, ok := <-self.ChMethods:
			{
				if !ok {
					log.Warnf("ChMethods channel not ok")
					return
				}

				misc.Assert(cmdMethods.Namespace != "", "bad cmdMethods")
				router := self.Get(cmdMethods.Namespace)
				conn, found := router.connMap[cmdMethods.ConnId]
				if found {
					router.UpdateServeMethods(conn, cmdMethods.Methods)
				} else {
					router.Log().Infof("Conn %d not found for update serve methods", cmdMethods.ConnId)
				}
			}

		case cmdDelg, ok := <-self.ChDelegates:
			{
				if !ok {
					log.Warnf("ChServe channel not ok")
					return
				}
				misc.Assert(cmdDelg.Namespace != "", "bad cmdDelg namespace")
				router := self.Get(cmdDelg.Namespace)
				conn, found := router.connMap[cmdDelg.ConnId]
				if found {
					router.UpdateDelegateMethods(conn, cmdDelg.MethodNames)
				} else {
					router.Log().Infof("Conn %d not found for update methods", cmdDelg.ConnId)
				}
			}

		case cmdMsg, ok := <-self.chMsg:
			{
				if !ok {
					log.Warnf("chMsg channel not ok")
					return
				}
				misc.Assert(cmdMsg.MsgVec.Namespace != "", "bad msgvec namespace")
				router := self.Get(cmdMsg.MsgVec.Namespace)
				go router.DeliverMessage(cmdMsg)
			}
		case <-time.After(10 * time.Second):
			for _, router := range self.routerMap {
				router.collectPendings()
			}
		}
	}
}
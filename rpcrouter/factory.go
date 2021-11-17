package rpcrouter

import (
	//"fmt"
	"context"
	//"errors"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/datadir"
	//jsonrpc "github.com/superisaac/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	//"time"
)

func NewRouterFactory(name string) *RouterFactory {

	factory := &RouterFactory{name: name}
	factory.Config = datadir.NewConfig()
	//factory.routerMap = make(map[string](*Router))
	//factory.setupChannels()
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

func (self *RouterFactory) Get(namespace string) *Router {
	misc.Assert(namespace != "", "factory got empty namespace")
	if r, ok := self.routerMap.Load(namespace); ok {
		router, _ := r.(*Router)
		return router
	} else {
		router := NewRouter(self, namespace)
		log.Debugf("router for namespace %s created", namespace)
		//self.routerMap[namespace] = t
		self.routerMap.Store(namespace, router)
		if self.Started() {
			router.EnsureStart(self.startCtx)
		}
		return router
	}
}

func (self RouterFactory) RouterNames() []string {
	names := make([]string, 0)
	self.routerMap.Range(func(key interface{}, value interface{}) bool {
		namespace, _ := key.(string)
		names = append(names, namespace)
		return true
	})
	// for namespace, _ := range self.routerMap {
	// 	names = append(names, namespace)
	// }
	return names
}

func (self RouterFactory) GetOrNil(name string) *Router {
	misc.Assert(name != "", "factory got empty namespace")
	if r, ok := self.routerMap.Load(name); ok {
		router, _ := r.(*Router)
		return router
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

func (self RouterFactory) Name() string {
	return self.name
}

func (self RouterFactory) Started() bool {
	return self.startCtx != nil
}

func (self *RouterFactory) EnsureStart(rootCtx context.Context) {
	if self.startCtx == nil {
		self.startCtx, self.cancelFunc = context.WithCancel(rootCtx)
		//go self.loop(self.startCtx)
		self.routerMap.Range(func(k, v interface{}) bool {
			router, _ := v.(*Router)
			router.EnsureStart(self.startCtx)
			return true
		})
	}
}

func (self *RouterFactory) Stop() {
	if self.Started() {
		self.cancelFunc()
		self.startCtx = nil
		self.cancelFunc = nil
	}
}

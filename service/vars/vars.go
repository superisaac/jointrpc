package vars

import (
	"context"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	datadir "github.com/superisaac/jointrpc/datadir"
	dispatch "github.com/superisaac/jointrpc/dispatch"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type VarsService struct {
	disp     *dispatch.Dispatcher
	chResult chan dispatch.ResultT
	//vars   map[string](map[string](interface{}))
	namedVars map[string]interface{}
	//router *rpcrouter.Router
	conn *rpcrouter.ConnT
}

func NewVarsService() *VarsService {
	srv := new(VarsService)
	return srv
}

func (self VarsService) Name() string {
	return "vars"
}
func (self VarsService) CanRun(rootCtx context.Context) bool {
	varsPath := datadir.Datapath("vars.yml")
	if _, err := os.Stat(varsPath); os.IsNotExist(err) {
		// the vars.yml doesnot exist
		return false
	}
	return true
}

func (self *VarsService) BroadcastVars(factory *rpcrouter.RouterFactory) error {
	for namespace, vars := range self.namedVars {
		router := factory.GetOrNil(namespace)
		if router == nil {
			continue
		}
		notify := jsonrpc.NewNotifyMessage("vars.change", []interface{}{vars})
		notify.SetTraceId(misc.NewUuid())
		notify.Log().Infof("broadcast vars.change")

		_, err := router.CallOrNotify(
			notify, namespace, rpcrouter.WithBroadcast(true))
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *VarsService) ReadVars(varsPath string) error {
	log.Infof("read vars from %s", varsPath)
	data, err := ioutil.ReadFile(varsPath)
	if err != nil {
		return err
	}
	namedVars := make(map[string]interface{})
	err = yaml.Unmarshal(data, namedVars)
	if err != nil {
		return err
	}
	self.namedVars = namedVars
	return nil
}

func (self *VarsService) Start(rootCtx context.Context) error {
	factory := rpcrouter.RouterFactoryFromContext(rootCtx)
	commonRouter := factory.CommonRouter()

	varsPath := datadir.Datapath("vars.yml")
	err := self.ReadVars(varsPath)
	if err != nil {
		return err
	}

	self.disp = dispatch.NewDispatcher()
	self.chResult = make(chan dispatch.ResultT, misc.DefaultChanSize())
	//self.conn = commonRouter.Join()
	self.conn = rpcrouter.NewConn()
	commonRouter.ChJoin <- rpcrouter.CmdJoin{Conn: self.conn}

	ctx, cancel := context.WithCancel(rootCtx)
	defer func() {
		cancel()
		commonRouter.ChLeave <- rpcrouter.CmdLeave{Conn: self.conn}
		//commonRouter.Leave(self.conn)
		self.conn = nil
	}()

	self.disp.On("_vars.list", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		//router := factory.GetOrNil(req.MsgVec.Namespace)
		if vars, ok := self.namedVars[req.MsgVec.Namespace]; ok {
			return vars, nil
		} else {
			return make(map[string]interface{}), nil
		}
	})

	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(varsPath)
	if err != nil {
		return err
	}

	self.declareMethods(factory)
	for {
		select {
		case <-ctx.Done():
			log.Debugf("vars handlers, context done")
			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			log.Debugf("vars watcher event %+v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Infof("modified file %s", event.Name)
				err := self.ReadVars(varsPath)
				if err != nil {
					return err
				}
				err = self.BroadcastVars(factory)
				if err != nil {
					return err
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Warnf("vars watcher error: %+v", err)
		case msgvec, ok := <-self.conn.RecvChannel:
			if !ok {
				log.Debugf("recv channel colosed, leave")
				return nil
			}
			//timeoutCtx, _ := context.WithTimeout(rootCtx, 10 * time.Second)
			self.requestReceived(ctx, msgvec)
		case cmdMsg, ok := <-self.conn.ChRouteMsg:
			if !ok {
				log.Debugf("ChRouteMsg closed")
				return nil
			}
			err := self.conn.HandleRouteMessage(ctx, cmdMsg)
			if err != nil {
				panic(err)
			}
		case result, ok := <-self.chResult:
			if !ok {
				log.Infof("result channel closed, return")
				return nil
			}
			self.conn.ChRouteMsg <- rpcrouter.CmdMsg{
				MsgVec: rpcrouter.MsgVec{
					Msg:        result.ResMsg,
					Namespace:  commonRouter.Name(),
				},
			}

			// commonRouter.DeliverResultOrError(
			// 	rpcrouter.MsgVec{
			// 		Msg:        result.ResMsg,
			// 		Namespace:  commonRouter.Name(),
			// 		FromConnId: self.conn.ConnId,
			// 	})
		}
	}
	return nil
}

func (self *VarsService) declareMethods(factory *rpcrouter.RouterFactory) {
	if self.conn != nil {
		minfos := self.disp.GetMethodInfos()
		cmdMethods := rpcrouter.CmdMethods{
			Namespace: factory.CommonRouter().Name(),
			ConnId:    self.conn.ConnId,
			Methods:   minfos,
		}
		factory.Get(cmdMethods.Namespace).ChMethods <- cmdMethods
	}
}

func (self *VarsService) requestReceived(ctx context.Context, msgvec rpcrouter.MsgVec) {
	msg := msgvec.Msg
	if msg.IsRequest() || msg.IsNotify() {
		self.disp.Feed(ctx, msgvec, self.chResult)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
}

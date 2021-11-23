package vars

import (
	"context"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	datadir "github.com/superisaac/jointrpc/datadir"
	dispatch "github.com/superisaac/jointrpc/dispatch"
	misc "github.com/superisaac/jointrpc/misc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
	jsonrpc "github.com/superisaac/jsonrpc"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	//"time"
)

type VarsService struct {
	disp     *dispatch.Dispatcher
	chResult chan dispatch.ResultT
	//vars   map[string](map[string](interface{}))
	namedVars map[string](map[string]interface{})
	//router *rpcrouter.Router
	conn *rpcrouter.ConnT
	done chan error
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
	namedVars := make(map[string](map[string]interface{}))
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
	self.done = make(chan error, 10)
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

	self.disp.OnTyped("_vars.list", func(req *dispatch.RPCRequest) (map[string]interface{}, error) {
		if vars, ok := self.namedVars[req.CmdMsg.Namespace]; ok {
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

	senderCtx, cancelSender := context.WithCancel(ctx)
	defer cancelSender()
	go dispatch.SenderLoop(senderCtx, self, self.conn, self.chResult)

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
		case err, ok := <-self.Done():
			if !ok {
				return nil
			} else if err != nil {
				return err
			} else {
				return nil
			}
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

func (self VarsService) SendMessage(ctx context.Context, msg jsonrpc.IMessage) error {
	factory := rpcrouter.RouterFactoryFromContext(ctx)
	commonRouter := factory.CommonRouter()
	self.conn.MsgInput() <- rpcrouter.CmdMsg{
		Msg:       msg,
		Namespace: commonRouter.Name(),
	}
	return nil
}

func (self VarsService) SendCmdMsg(ctx context.Context, cmdMsg rpcrouter.CmdMsg) error {
	msg := cmdMsg.Msg
	if msg.IsRequestOrNotify() {
		self.disp.Feed(ctx, cmdMsg, self.chResult)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
	return nil
}

func (self VarsService) Done() chan error {
	return self.done
}

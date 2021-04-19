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
	disp *dispatch.Dispatcher
	//vars   map[string](map[string](interface{}))
	vars   map[string]interface{}
	router *rpcrouter.Router
	conn   *rpcrouter.ConnT
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

func (self *VarsService) BroadcastVars() error {
	notify := jsonrpc.NewNotifyMessage("vars.change", []interface{}{self.vars}, nil)
	notify.SetTraceId(misc.NewUuid())
	notify.Log().Infof("broadcast vars.change")

	_, err := self.router.CallOrNotify(notify,
		rpcrouter.WithBroadcast(true))
	if err != nil {
		return err
	}
	return nil
}

func (self *VarsService) ReadVars(varsPath string) error {
	log.Infof("read vars from %s", varsPath)
	data, err := ioutil.ReadFile(varsPath)
	if err != nil {
		return err
	}
	vars := make(map[string]interface{})
	err = yaml.Unmarshal(data, vars)
	if err != nil {
		return err
	}
	self.vars = vars
	return nil
}

func (self *VarsService) Start(rootCtx context.Context) error {
	varsPath := datadir.Datapath("vars.yml")
	err := self.ReadVars(varsPath)
	if err != nil {
		return err
	}

	self.disp = dispatch.NewDispatcher()
	self.router = rpcrouter.RouterFromContext(rootCtx)
	self.conn = self.router.Join()
	ctx, cancel := context.WithCancel(rootCtx)
	defer func() {
		cancel()
		self.router.Leave(self.conn)
		self.conn = nil
	}()

	self.disp.On("_vars.list", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		return self.vars, nil
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

	self.declareMethods()
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
				err = self.BroadcastVars()
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
			self.requestReceived(msgvec)
		case resmsg, ok := <-self.disp.ChResult:
			if !ok {
				log.Infof("result channel closed, return")
				return nil
			}
			self.router.DeliverResultOrError(rpcrouter.MsgVec{
				Msg:        resmsg,
				FromConnId: self.conn.ConnId})
		}
	}
	return nil
}

func (self *VarsService) declareMethods() {
	if self.conn != nil {
		minfos := make([]rpcrouter.MethodInfo, 0)
		for m, info := range self.disp.MethodHandlers {
			minfo := rpcrouter.MethodInfo{
				Name:       m,
				Help:       info.Help,
				SchemaJson: info.SchemaJson,
			}
			minfos = append(minfos, minfo)
		}
		cmdServe := rpcrouter.CmdServe{ConnId: self.conn.ConnId, Methods: minfos}
		self.router.ChServe <- cmdServe
	}
}

func (self *VarsService) requestReceived(msgvec rpcrouter.MsgVec) {
	msg := msgvec.Msg
	if msg.IsRequest() || msg.IsNotify() {
		self.disp.HandleRequestMessage(msgvec)
	} else {
		log.Warnf("builtin handler, receved none request msg %+v", msg)
	}
}

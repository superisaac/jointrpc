package playbook

import (
	//"fmt"
	"context"
	//"errors"
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/jsonrpc/schema"
	yaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func NewPlaybook() *Playbook {
	return &Playbook{}
}

func (self *PlaybookConfig) ReadConfig(filePath string) error {
	log.Infof("read playbook from %s", filePath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, self)
	if err != nil {
		return err
	}
	err = self.validateValues()
	return err
}

func (self *PlaybookConfig) validateValues() error {
	if self.Version == "" {
		self.Version = "1.0"
	}

	for _, method := range self.Methods {
		if method.SchemaInterface != nil {
			builder := schema.NewSchemaBuilder()
			_, err := builder.BuildYAMLInterface(method.SchemaInterface)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (self MethodT) CanExec() bool {
	return self.Shell != nil && self.Shell.Cmd != ""
}

func (self MethodT) Exec(req *dispatch.RPCRequest, methodName string) (interface{}, error) {
	msg := req.MsgVec.Msg
	// if !msg.IsRequestOrNotify() {
	// }
	var ctx context.Context
	var cancel func()
	if self.Shell.Timeout != nil {
		ctx, cancel = context.WithTimeout(
			context.Background(),
			time.Second*time.Duration(*self.Shell.Timeout))
		defer cancel()
	} else {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, "sh", "-c", self.Shell.Cmd)

	cmd.Env = append(os.Environ(), self.Shell.Env...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	msgJson := msg.MustString()
	io.WriteString(stdin, msgJson)
	stdin.Close()

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if cmd.Process != nil {
		msg.Log().Infof("command for %s received output, pid %#v", methodName, cmd.Process.Pid)
	}
	parsed, err := simplejson.NewJson(out)
	if err != nil {
		return nil, err
	}
	return parsed.Interface(), nil
}

func (self *Playbook) Run(serverEntry client.ServerEntry) error {
	rpcClient := client.NewRPCClient(serverEntry)
	disp := dispatch.NewDispatcher()
	disp.SetSpawnExec(true)

	for name, method := range self.Config.Methods {
		if !method.CanExec() {
			log.Warnf("cannot exec method %s %+v %s\n", name, method, method.Shell.Cmd)
			continue
		}
		opts := make([]func(*dispatch.MethodHandler), 0)
		if method.SchemaInterface != nil {
			n := simplejson.New()
			n.SetPath(nil, method.SchemaInterface)
			sb, err := n.MarshalJSON()
			if err != nil {
				return err
			}
			opts = append(opts, dispatch.WithSchema(string(sb)))
		}
		if method.Description != "" {
			opts = append(opts, dispatch.WithHelp(method.Description))
		}

		disp.On(name, func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			req.MsgVec.Msg.Log().Infof("begin exec %s", name)
			v, err := method.Exec(req, name)
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					req.MsgVec.Msg.Log().Warnf(
						"command exit, code: %d, stderr: %s",
						exitErr.ExitCode(),
						string(exitErr.Stderr)[:100])
					return nil, jsonrpc.ErrWorkerExit
				}

				req.MsgVec.Msg.Log().Warnf("error exec %s, %s", name, err.Error())
			} else {
				req.MsgVec.Msg.Log().Infof("end exec %s", name)
			}
			return v, err
		}, opts...)
	}

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	return rpcClient.Worker(context.Background(), disp)
}

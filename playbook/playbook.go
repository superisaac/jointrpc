package playbook

import (
	//"fmt"
	"context"
	//"errors"
	"github.com/bitly/go-simplejson"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jsonz"
	"github.com/superisaac/jsonz/schema"
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
	return self.ReadConfigBytes(data)
}

func (self *PlaybookConfig) ReadConfigBytes(data []byte) error {
	err := yaml.Unmarshal(data, self)
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
			s, err := builder.BuildYamlInterface(method.SchemaInterface)
			if err != nil {
				return err
			}
			method.innerSchema = s
			method.SchemaInterface = s.RebuildType()
		}
	}
	return nil
}

func (self MethodT) CanExec() bool {
	return self.Shell != nil && self.Shell.Cmd != ""
}

func (self MethodT) Exec(req *dispatch.RPCRequest, methodName string) (interface{}, error) {
	msg := req.CmdMsg.Msg
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

	msgJson := jsonz.MessageString(msg)
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
		log.Infof("playbook register %s", name)
		opts := make([]func(*dispatch.MethodHandler), 0)
		if method.innerSchema != nil {
			schemaJson := schema.SchemaToString(method.innerSchema)
			opts = append(opts, dispatch.WithSchema(schemaJson))
		}
		if method.Description != "" {
			opts = append(opts, dispatch.WithHelp(method.Description))
		}

		disp.On(name, func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			req.CmdMsg.Msg.Log().Infof("begin exec %s", name)
			v, err := method.Exec(req, name)
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					req.CmdMsg.Msg.Log().Warnf(
						"command exit, code: %d, stderr: %s",
						exitErr.ExitCode(),
						string(exitErr.Stderr)[:100])
					return nil, jsonz.ErrLiveExit
				}

				req.CmdMsg.Msg.Log().Warnf("error exec %s, %s", name, err.Error())
			} else {
				req.CmdMsg.Msg.Log().Infof("end exec %s", name)
			}
			return v, err
		}, opts...)
	}

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	return rpcClient.Live(context.Background(), disp)
}

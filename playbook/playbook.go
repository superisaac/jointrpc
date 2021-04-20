package playbook

import (
	//"fmt"
	"context"
	"errors"
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/jsonrpc/schema"
	yaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	//jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

func NewPlaybook() *Playbook {
	return &Playbook{}
}

//
func fixStringMaps(src interface{}) (interface{}, bool) {
	if anyMap, ok := src.(map[interface{}]interface{}); ok {
		strMap := make(map[string]interface{})
		for k, v := range anyMap {
			if sk, ok := k.(string); ok {
				if newV, ok := fixStringMaps(v); ok {
					strMap[sk] = newV
				} else {
					return nil, false
				}
			} else {
				return nil, false
			}
		}
		return strMap, true
	} else if anyList, ok := src.([]interface{}); ok {
		list1 := make([]interface{}, 0)
		for _, elem := range anyList {
			newElem, ok := fixStringMaps(elem)
			if !ok {
				return nil, false
			}
			list1 = append(list1, newElem)
		}
		return list1, true
	} else {
		return src, true
	}
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
			sintf, ok := fixStringMaps(method.SchemaInterface)
			if !ok {
				return errors.New("cannot convert to string map")
			}
			method.SchemaInterface = sintf
			builder := schema.NewSchemaBuilder()
			_, err := builder.Build(method.SchemaInterface)
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

func (self MethodT) Exec(req *dispatch.RPCRequest) (interface{}, error) {
	msg := req.MsgVec.Msg
	// if !msg.IsRequestOrNotify() {
	// }
	cmd := exec.Command("sh", "-c", self.Shell.Cmd)
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
			return method.Exec(req)
		}, opts...)
	}

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	return rpcClient.Worker(context.Background(), disp)
}

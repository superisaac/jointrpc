package datadir

import (
	"os"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
)

var cfg Config
var cfgParsed bool = false

func GetConfig() Config {
	if !cfgParsed {
		err := cfg.ParseConfig()
		if err != nil {
			panic(err)
		}
		cfgParsed = true
	}
	return cfg
}

func (self *Config) ParseConfig() error {
	cfgPath := DataPath("config.yml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		err = self.ValidateValues()
		if err != nil {
			return err
		}
		return nil
	}
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, self)
	if err != nil {
		return err
	}
	return self.ValidateValues()
}

func (self *Config) ValidateValues() error {
	if self.Version == "" {
		self.Version = "1.0"
	}

	if self.Server.Bind == "" {
		self.Server.Bind = "127.0.0.1:38600"
	}

	return nil
}



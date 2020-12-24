package datadir

import (
	"errors"
	log "github.com/sirupsen/logrus"
	logsyslog "github.com/sirupsen/logrus/hooks/syslog"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log/syslog"
	"os"
)

var (
	cfg       *Config
)

func GetConfig() *Config {
	if cfg == nil {
		cfg = new(Config)
		err := cfg.ParseConfig()
		if err != nil {
			panic(err)
		}
		cfg.setupLogger()
	}
	return cfg
}

func (self *Config) ParseConfig() error {
	cfgPath := Datapath("config.yml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		err = self.validateValues()
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
	return self.validateValues()
}

func (self *Config) validateValues() error {
	if self.Version == "" {
		self.Version = "1.0"
	}

	if self.Server.Bind == "" {
		self.Server.Bind = "127.0.0.1:50055"
	}

	if self.Logging.Syslog.URL == "" {
		self.Logging.Syslog.URL = "localhost:514"
	}
	if self.Logging.Syslog.Protocol == "" {
		self.Logging.Syslog.Protocol = "udp"
	}
	if self.Logging.Syslog.Protocol != "udp" &&
		self.Logging.Syslog.Protocol != "tcp" {
		return errors.New("config, invalid syslog protocol")
	}

	if len(self.Cluster.StaticPeers) > 0 && self.Cluster.AdvertisedURL == "" {
		return errors.New("config, advertised url must be specified")
	}

	return nil
}

func (self *Config) setupLogger() {
	log.SetFormatter(&log.JSONFormatter{})

	logOutput := self.Logging.Output
	if logOutput == "" || logOutput == "console" || logOutput == "stdout" {
		log.SetOutput(os.Stdout)
	} else if logOutput == "stderr" {
		log.SetOutput(os.Stderr)
	} else {
		file, err := os.OpenFile(logOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}

	if self.Logging.Syslog.Enabled {
		hook, err := logsyslog.NewSyslogHook(
			self.Logging.Syslog.Protocol,
			self.Logging.Syslog.URL,
			syslog.LOG_INFO, "")
		if err != nil {
			panic(err)
		}
		log.AddHook(hook)
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	switch envLogLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

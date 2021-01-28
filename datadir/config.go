package datadir

import (
	"errors"
	log "github.com/sirupsen/logrus"
	logsyslog "github.com/sirupsen/logrus/hooks/syslog"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log/syslog"
	"net"
	"os"
	"path/filepath"
)

func NewConfig() *Config {
	return new(Config)
}

func (self *Config) ParseDatadir() error {
	cfgPath := Datapath("server.yml")
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

func (self Config) ValidateSchema() bool {
	if self.pValidateSchema == nil {
		return true
	} else {
		return *self.pValidateSchema
	}
}

func (self *Config) validateValues() error {
	if self.Version == "" {
		self.Version = "1.0"
	}

	if self.Server.Bind == "" {
		self.Server.Bind = "127.0.0.1:50055"
	}

	if self.pValidateSchema == nil {
		v := true
		self.pValidateSchema = &v
	}

	// tls
	if self.Server.TLS.CertFile != "" {
		certFile := filepath.Join(Datapath("tls/"), self.Server.TLS.CertFile)
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			return errors.New("config, tls certification file does not exist")
		}
		self.Server.TLS.CertFile = certFile
	}
	if self.Server.TLS.KeyFile != "" {
		keyFile := filepath.Join(Datapath("tls/"), self.Server.TLS.KeyFile)
		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			return errors.New("config, tls key file does not exist")
		}
		self.Server.TLS.KeyFile = keyFile
	}

	// authorizations
	for _, bauth := range self.Authorizations {
		err := bauth.validateValues()
		if err != nil {
			return nil
		}
	}

	// syslog
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

	return nil
}

func (self *Config) SetupLogger() {
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

func (self BasicAuth) Authorize(username string, password string, ipAddr string) bool {
	return self.checkUser(username, password) && self.checkIP(ipAddr)
}

func (self BasicAuth) checkUser(username string, password string) bool {
	return self.Username == username && self.Password == password
}

func (self BasicAuth) checkIP(ipAddr string) bool {
	if len(self.AllowedSources) > 0 {
		ip := net.ParseIP(ipAddr)
		if ip == nil {
			log.Errorf("parse ip failed %s", ipAddr)
		}

		if self.allowedIPNets != nil {
			for _, ipnet := range self.allowedIPNets {
				if ipnet.Contains(ip) {
					return true
				}
			}
		}
		return false
	} else {
		return true
	}
}

func (self *BasicAuth) validateValues() error {
	if self.AllowedSources != nil {
		allowedIPNets := make([]*net.IPNet, 0)
		for _, cidrStr := range self.AllowedSources {
			_, ipnet, err := net.ParseCIDR(cidrStr)
			if err != nil {
				return err
			}
			allowedIPNets = append(allowedIPNets, ipnet)
		}
		self.allowedIPNets = allowedIPNets
	}
	return nil
}

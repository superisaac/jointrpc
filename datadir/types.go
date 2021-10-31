package datadir

import (
	"net"
)

type ServerTLSConfig struct {
	CertFile string `yaml:"cert,omitempty"`
	KeyFile  string `yaml:"key,omitempty"`
}

type ServerConfig struct {
	Bind     string          `yaml:"bind,omitempty"`
	HttpBind string          `yaml:"http_bind,omitempty"`
	TLS      ServerTLSConfig `yaml:"tls,omitempty"`
}

type BasicAuth struct {
	Username       string       `yaml:"username"`
	Password       string       `yaml:"password"`
	Namespace      string       `yaml:"namespace,omitempty"`
	AllowedSources []string     `yaml:"allow,omitempty"`
	allowedIPNets  []*net.IPNet `yaml:"-"`
}

type SyslogConfig struct {
	Enabled  bool   `yaml:"enabled,omitempty"`
	URL      string `yaml:"url,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
}

type LoggingConfig struct {
	Output string       `yaml:"output,omitempty"`
	Level  string       `yaml:"level,omitempty"`
	Syslog SyslogConfig `yaml:"syslog,omitempty"`
}

type PeerConfig struct {
	ServerUrl string `yaml:"url,omitempty"`
	CertFile  string `yaml:"cert,omitempty"`
}

type MetricsConfig struct {
	BearerToken string `yaml:"bearer_token,omitempty"`
}

type NeighborConfig struct {
	Peers []PeerConfig
}

// The root config item
type Config struct {
	Version         string                    `yaml:"version"`
	Logging         LoggingConfig             `yaml:"logging,omitempty"`
	Server          ServerConfig              `yaml:"server"`
	Authorizations  []BasicAuth               `yaml:"auth,omitempty"`
	Metrics         MetricsConfig             `yaml:"metrics,omitempty"`
	Neighbors       map[string]NeighborConfig `yaml:"neighbors,omitempty"`
	pValidateSchema *bool                     `yaml:"validate_schema,omitempty"`
}

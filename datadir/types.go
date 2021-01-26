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
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	AllowedCIDR string `yaml:"allowedCIDR,omitempty"`
	cidrIP      net.IP
	cidrIPNet   *net.IPNet
}

type SyslogConfig struct {
	Enabled  bool   `yaml:"enabled,omitempty"`
	URL      string `yaml:"url,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
}

type LoggingConfig struct {
	Output string       `yaml,"path,omitempty"`
	Level  string       `yaml,"level,omitempty"`
	Syslog SyslogConfig `yaml:"syslog,omitempty"`
}

type PeerConfig struct {
	ServerUrl string `yaml:"url,omitempty"`
	CertFile  string `yaml:"cert,omitempty"`
}

type MetricsConfig struct {
	BearerToken string `yaml:"bearer_token,omitempty"`
}

type ClusterConfig struct {
	AdvertisedURL string       `yaml:"advertised_url"`
	NeighborPeers []PeerConfig `yaml:neighbor_peers,omitempty"`
}

// The root config item
type Config struct {
	Version         string        `yaml:"version"`
	Logging         LoggingConfig `yaml:"log,omitempty"`
	Server          ServerConfig  `yaml:"server"`
	Authorizations  []BasicAuth   `yaml:"auth,omitempty"`
	Metrics         MetricsConfig `yaml:"metrics"`
	Cluster         ClusterConfig `yaml:"cluster,omitempty"`
	pValidateSchema *bool         `yaml:"validate_schema,omitempty"`
}

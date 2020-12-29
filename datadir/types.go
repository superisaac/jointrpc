package datadir

type ServerTLSConfig struct {
	CertFile      string `yaml:"cert,omitempty"`
	KeyFile       string `yaml:"key,omitempty"`
}

type ServerConfig struct {
	Bind      string `yaml:"bind,omitempty"`
	ProxyBind string `yaml:"proxy_bind,omitempty"`
	TLS       ServerTLSConfig `yaml:"tls,omitempty"`
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
	Address string `yaml:"url,omitempty"`
	CertFile string `yaml:"cert,omitempty"`
}

type ClusterConfig struct {
	AdvertisedURL string       `yaml:"advertised_url"`
	StaticPeers   []PeerConfig `yaml:static_peers,omitempty"`
}

// The root config item
type Config struct {
	Version string        `yaml:"version"`
	Logging LoggingConfig `yaml:"log,omitempty"`
	Server  ServerConfig  `yaml:"server"`
	Cluster ClusterConfig `yaml:"cluster,omitempty"`
}

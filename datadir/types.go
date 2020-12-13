package datadir

type ServerConfig struct {
	Bind      string `yaml:"bind,omitempty"`
	ProxyBind string `yaml:"proxy_bind,omitempty"`
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
	URL string `yaml:"url,omitempty"`
	// TODO: certs
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

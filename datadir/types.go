package datadir

type ServerConfig struct {
	Bind string `yaml:"bind,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"`
	ProxyBind string `yaml:"proxy_bind,omitempty"`
}

// type DatastoreConfig struct {
// 	DSN string `yaml:"dsn,omitempty"`
// }

type Config struct {
	Version string `yaml:"version,omitempty"`
	Server ServerConfig `yaml:"server,omitempty"`
	//Datastore DatastoreConfig `yaml:"datastore,omitempty"`
}


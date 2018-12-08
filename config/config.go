package config

type Config struct {
	Server      ServerConfig      `toml:"server"`
	Loader      []LoaderConfig    `toml:"loaders"`
	Development DevelopmentConfig `toml:"dev"`
}

type ServerConfig struct {
	Token string `toml:"token"`
	Port  int    `toml:"port"`
	Debug bool   `toml:"debug"`
}

type LoaderConfig struct {
	Path string `toml:"path"`
}

type DevelopmentConfig struct {
	DebugMode      bool `toml:"debug"`
	IgnoreProtocol bool `toml:"ingore_protocol"`
}

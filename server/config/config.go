package config

type Config struct {
	Server ServerConfig   `toml:"server"`
	Loader []LoaderConfig `toml:"loaders"`
}

type ServerConfig struct {
	Port  int  `toml:"port"`
	Debug bool `toml:"debug"`
}

type LoaderConfig struct {
	Path string `toml:"path"`
}

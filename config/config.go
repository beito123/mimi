package config

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

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
	Name string `toml:"name"`
	Type string `toml:"type"`
	Path string `toml:"path"`
}

type DevelopmentConfig struct {
	DebugMode      bool `toml:"debug"`
	IgnoreProtocol bool `toml:"ingore_protocol"`
}

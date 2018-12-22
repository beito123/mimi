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
	Programs      []ProgramConfig    `toml:"programs"`
	Development DevelopmentConfig `toml:"dev"`
}

type ServerConfig struct {
	Token string `toml:"token"`
	Port  int    `toml:"port"`
	Debug bool   `toml:"debug"`
}

type ProgramConfig struct {
	Name    string            `toml:"name"`
	Path    string            `toml:"path"`
	Loader    string            `toml:"loader"`
	LoaderOptions map[string]string `toml:"loader_options"`
}

type DevelopmentConfig struct {
	DebugMode      bool `toml:"debug"`
	IgnoreProtocol bool `toml:"ingore_protocol"`
}

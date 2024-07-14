package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   Server
	Name     string `yaml:"name"`
	Database Database
	Debug    bool
}

type Server struct {
	Host       string `yaml:"host"`
	CorsHost   string `yaml:"cors_host"`
	WithCors   bool   `yaml:"with_cors"`
	Port       int    `yaml:"port"`
	StaticPath string `yaml:"static_path"`
}

type Database struct {
	Type    string `yaml:"type"`
	Address string `yaml:"address"`
	Cache   string `yaml:"cache"`
	Schema  string `yaml:"schema"`
	MaxConn int    `yaml:"max_conn"`
}

var config *Config

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	config = &conf
	return config, nil
}

func Get() *Config {
	return config
}

package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Driver        string `yaml:"driver"`
	Name          string `yaml:"name"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	TCPConnection string `yaml:"tcp_connection"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

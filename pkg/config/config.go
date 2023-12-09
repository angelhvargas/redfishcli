package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Type     string `yaml:"type"`
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Servers []ServerConfig `yaml:"servers"`
}

type IDRACConfig struct {
	Hostname string
	Username string
	Password string
}

type XClarityConfig struct {
	Hostname string
	Username string
	Password string
	// Any other specific fields
}

// LoadConfig reads a YAML configuration file and unmarshals it into a Config struct.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

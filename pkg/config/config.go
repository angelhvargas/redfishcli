package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type BMCConfig struct {
	Type           string
	IDRACConfig    IDRACConfig
	XClarityConfig XClarityConfig
	Servers        []ServerConfig
}

type ServerConfig struct {
	Type     string `yaml:"type"`
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// type Config struct {
// 	Servers []ServerConfig `yaml:"servers"`
// }

type BMCConnConfig struct {
	Hostname       string
	Username       string
	Password       string
	ControllerType string
}

type IDRACConfig struct {
	BMCConnConfig
}

type XClarityConfig struct {
	BMCConnConfig
}

// LoadConfig reads a YAML configuration file and unmarshals it into a Config struct.
func LoadConfig(path string) (*BMCConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg BMCConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

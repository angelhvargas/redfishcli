package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type BMCConfig struct {
	Servers []ServerConfig `yaml:"servers"`
}

type ServerConfig struct {
	Type     string `yaml:"type"`
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

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

// LoadConfig reads a YAML configuration file and unmarshals it into a BMCConfig struct.
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

// GetDefaultConfigPath returns the default configuration file path.
func GetDefaultConfigPath() string {
	return os.Getenv("HOME") + "/.redfishcli/config.yaml"
}

// LoadConfigFromDefaultPath loads the configuration from the default path.
func LoadConfigFromDefaultPath() (*BMCConfig, error) {
	return LoadConfig(GetDefaultConfigPath())
}

// LoadConfigOrEnv loads the configuration from the specified path, or from environment variables if the path is empty.
func LoadConfigOrEnv(path, bmcType, username, password, hostname string) (*BMCConfig, error) {
	if path != "" {
		return LoadConfig(path)
	}

	server := ServerConfig{
		Type:     bmcType,
		Hostname: hostname,
		Username: username,
		Password: password,
	}

	cfg := &BMCConfig{
		Servers: []ServerConfig{server},
	}

	return cfg, nil
}

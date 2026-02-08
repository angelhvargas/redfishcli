package config

import (
	"os"

	"github.com/angelhvargas/redfishcli/pkg/logger"
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
	logger.Log.Infof("Loading configuration from %s", path)
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
	var cfg *BMCConfig
	var err error

	if path != "" {
		cfg, err = LoadConfig(path)
		if err != nil {
			return nil, err
		}
	} else {
		cfg = &BMCConfig{
			Servers: []ServerConfig{},
		}
	}

	// Override with CLI or environment variables if necessary
	for i, server := range cfg.Servers {
		if server.Username == "" {
			cfg.Servers[i].Username = username
			if envUser := os.Getenv("BMC_USERNAME"); envUser != "" {
				cfg.Servers[i].Username = envUser
			}
		}
		if server.Password == "" {
			cfg.Servers[i].Password = password
			if envPass := os.Getenv("BMC_PASSWORD"); envPass != "" {
				cfg.Servers[i].Password = envPass
			}
		}
	}

	// Add a single server configuration if CLI flags are provided and no servers are defined
	if len(cfg.Servers) == 0 && hostname != "" {
		server := ServerConfig{
			Type:     bmcType,
			Hostname: hostname,
			Username: username,
			Password: password,
		}
		if envUser := os.Getenv("BMC_USERNAME"); envUser != "" {
			server.Username = envUser
		}
		if envPass := os.Getenv("BMC_PASSWORD"); envPass != "" {
			server.Password = envPass
		}
		cfg.Servers = append(cfg.Servers, server)
	}

	return cfg, nil
}

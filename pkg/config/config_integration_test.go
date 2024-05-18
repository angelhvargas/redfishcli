//go:build integration
// +build integration

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempConfigFile(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		return "", err
	}

	_, err = tempFile.Write([]byte(content))
	if err != nil {
		tempFile.Close()
		return "", err
	}

	err = tempFile.Close()
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func TestLoadConfigIntegration(t *testing.T) {
	t.Run("LoadConfig with valid file", func(t *testing.T) {
		content := `
servers:
  - type: "idrac"
    hostname: "192.168.1.1"
    username: "root"
    password: "password1"
  - type: "xclarity"
    hostname: "192.168.1.2"
    username: "admin"
    password: "password2"
`
		configFilePath, err := createTempConfigFile(content)
		assert.NoError(t, err)
		defer os.Remove(configFilePath)

		cfg, err := LoadConfig(configFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Len(t, cfg.Servers, 2)
		assert.Equal(t, "idrac", cfg.Servers[0].Type)
		assert.Equal(t, "192.168.1.1", cfg.Servers[0].Hostname)
		assert.Equal(t, "root", cfg.Servers[0].Username)
		assert.Equal(t, "password1", cfg.Servers[0].Password)
		assert.Equal(t, "xclarity", cfg.Servers[1].Type)
		assert.Equal(t, "192.168.1.2", cfg.Servers[1].Hostname)
		assert.Equal(t, "admin", cfg.Servers[1].Username)
		assert.Equal(t, "password2", cfg.Servers[1].Password)
	})

	t.Run("LoadConfig from default path", func(t *testing.T) {
		// Create a temporary directory to act as the HOME directory
		tempDir, err := os.MkdirTemp("", "home")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Set the HOME environment variable to the temporary directory
		homeDir := os.Getenv("HOME")
		defer os.Setenv("HOME", homeDir)
		os.Setenv("HOME", tempDir)

		// Create the .redfishcli directory within the temporary HOME directory
		configDir := filepath.Join(tempDir, ".redfishcli")
		err = os.Mkdir(configDir, 0755)
		assert.NoError(t, err)

		// Create the config.yaml file within the .redfishcli directory
		configFilePath := filepath.Join(configDir, "config.yaml")
		content := `
servers:
  - type: "idrac"
    hostname: "192.168.1.1"
    username: "root"
    password: "password1"
`
		err = os.WriteFile(configFilePath, []byte(content), 0644)
		assert.NoError(t, err)

		cfg, err := LoadConfigFromDefaultPath()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Len(t, cfg.Servers, 1)
		assert.Equal(t, "idrac", cfg.Servers[0].Type)
		assert.Equal(t, "192.168.1.1", cfg.Servers[0].Hostname)
		assert.Equal(t, "root", cfg.Servers[0].Username)
		assert.Equal(t, "password1", cfg.Servers[0].Password)
	})

	t.Run("Handle missing configuration file", func(t *testing.T) {
		_, err := LoadConfig("/non/existent/path/config.yaml")
		assert.Error(t, err)
	})

	t.Run("LoadConfig with environment variables", func(t *testing.T) {
		cfg, err := LoadConfigOrEnv("", "idrac", "root", "password1", "192.168.1.1")
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Len(t, cfg.Servers, 1)
		assert.Equal(t, "idrac", cfg.Servers[0].Type)
		assert.Equal(t, "192.168.1.1", cfg.Servers[0].Hostname)
		assert.Equal(t, "root", cfg.Servers[0].Username)
		assert.Equal(t, "password1", cfg.Servers[0].Password)
	})
}

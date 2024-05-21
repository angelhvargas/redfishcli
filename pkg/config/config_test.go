package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("LoadConfig with valid file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "config-*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

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
		_, err = tempFile.Write([]byte(content))
		assert.NoError(t, err)
		tempFile.Close()

		cfg, err := LoadConfig(tempFile.Name())
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

	t.Run("LoadConfig with invalid file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "config-*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		content := `invalid yaml content`
		_, err = tempFile.Write([]byte(content))
		assert.NoError(t, err)
		tempFile.Close()

		_, err = LoadConfig(tempFile.Name())
		assert.Error(t, err)
	})
}

func TestLoadConfigOrEnv(t *testing.T) {
	t.Run("LoadConfigOrEnv with path", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "config-*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		content := `
servers:
  - type: "idrac"
    hostname: "192.168.1.1"
    username: "root"
    password: "password1"
`
		_, err = tempFile.Write([]byte(content))
		assert.NoError(t, err)
		tempFile.Close()

		cfg, err := LoadConfigOrEnv(tempFile.Name(), "", "", "", "")
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Len(t, cfg.Servers, 1)
		assert.Equal(t, "idrac", cfg.Servers[0].Type)
		assert.Equal(t, "192.168.1.1", cfg.Servers[0].Hostname)
		assert.Equal(t, "root", cfg.Servers[0].Username)
		assert.Equal(t, "password1", cfg.Servers[0].Password)
	})

	t.Run("LoadConfigOrEnv with environment variables", func(t *testing.T) {
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

func TestLoadConfigFromDefaultPath(t *testing.T) {
	t.Run("LoadConfigFromDefaultPath with valid file", func(t *testing.T) {
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
		configFile := filepath.Join(configDir, "config.yaml")
		err = os.WriteFile(configFile, []byte(`
servers:
  - type: "idrac"
    hostname: "192.168.1.1"
    username: "root"
    password: "password1"
`), 0644)
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
}

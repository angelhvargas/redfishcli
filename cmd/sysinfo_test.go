package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestSysInfoCmd(t *testing.T) {
	// Setup config file
	configFile := "config_test.yaml"
	configContent := `
servers:
  - type: idrac
    hostname: test-server
    username: user
    password: password
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(configFile)

	// Mock Client
	mockClient := new(client.MockServerClient)
	client.ResetRegistry()
	client.Register("idrac", func(cfg config.BMCConnConfig) client.ServerClient {
		return mockClient
	})

	mockInfo := &model.ServerInfo{
		ID:           "test-server-id",
		Manufacturer: "Dell",
		Model:        "R740",
	}
	mockClient.On("GetServerInfo").Return(mockInfo, nil)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command
	rootCmd.SetArgs([]string{"sysinfo", "--config", configFile})
	err = rootCmd.Execute() // This might fail if rootCmd isn't exported or initialized correctly in test context, assuming it works for now

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NoError(t, err)
	assert.Contains(t, output, "Manufacturer: Dell")
	assert.Contains(t, output, "Model: R740")
	mockClient.AssertExpectations(t)
}

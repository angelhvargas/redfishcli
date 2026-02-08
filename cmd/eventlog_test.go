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

func TestEventLogCmd(t *testing.T) {
	// Setup config file
	configFile := "config_test_eventlog.yaml"
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

	mockLogs := []model.EventLogEntry{
		{
			Created:  "2024-01-01T12:00:00Z",
			Severity: "Critical",
			Message:  "System Overheated",
		},
	}
	mockClient.On("GetSystemEventLog").Return(mockLogs, nil).Once()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command
	rootCmd.SetArgs([]string{"eventlog", "--config", configFile})
	err = rootCmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NoError(t, err)
	assert.Contains(t, output, "System Overheated")
	assert.Contains(t, output, "Critical")
	mockClient.AssertExpectations(t)
}

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

func TestBootCmd(t *testing.T) {
	// Setup config file
	configFile := "config_test_boot.yaml"
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

	t.Run("boot status", func(t *testing.T) {
		mockInfo := &model.BootInfo{
			BootOrder: []string{"Hdd", "Pxe"},
		}
		mockClient.On("GetBootInfo").Return(mockInfo, nil).Once()

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd.SetArgs([]string{"boot", "status", "--config", configFile})
		err := rootCmd.Execute()

		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		assert.NoError(t, err)
		assert.Contains(t, output, "Hdd")
		assert.Contains(t, output, "Pxe")
		mockClient.AssertExpectations(t)
	})

	t.Run("boot set", func(t *testing.T) {
		mockClient.On("SetBootOrder", "Pxe").Return(nil).Once()

		rootCmd.SetArgs([]string{"boot", "set", "--device", "Pxe", "--config", configFile})
		err := rootCmd.Execute()

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})
}

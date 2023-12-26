// client/idrac/idrac_integration_test.go

//go:build integration
// +build integration

package idrac

import (
	"os"
	"testing"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIDRACClientIntegration(t *testing.T) {
	// Skipping test if no environment variable is set
	hostname := os.Getenv("IDRAC_HOSTNAME")
	if hostname == "" {
		t.Skip("IDRAC_HOSTNAME environment variable not set")
	}

	client := NewClient(config.IDRACConfig{
		Hostname: hostname,
		Username: os.Getenv("IDRAC_USERNAME"),
		Password: os.Getenv("IDRAC_PASSWORD"),
	})

	t.Run("GetServerInfo", func(t *testing.T) {
		_, err := client.GetServerInfo()
		require.NoError(t, err)
		// Add more assertions here based on expected results
	})

	t.Run("GetStorageInfo", func(t *testing.T) {
		_, err := client.GetStorageInfo()
		require.NoError(t, err)
		// Add more assertions here based on expected results
	})
}

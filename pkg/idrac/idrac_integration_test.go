// client/idrac/idrac_integration_test.go

//go:build integration
// +build integration

package idrac

import (
	"os"
	"testing"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/stretchr/testify/require"
)

func TestIDRACClientIntegration(t *testing.T) {
	// Skipping test if no environment variable is set
	hostname := os.Getenv("IDRAC_HOSTNAME")
	if hostname == "" {
		t.Skip("IDRAC_HOSTNAME environment variable not set")
	}

	client := NewClient(config.IDRACConfig{
		BMCConnConfig: config.BMCConnConfig{
			Hostname: hostname,
			Username: os.Getenv("IDRAC_USERNAME"),
			Password: os.Getenv("IDRAC_PASSWORD"),
		},
	})

	t.Run("GetServerInfo", func(t *testing.T) {
		serverInfo, err := client.GetServerInfo()
		require.NoError(t, err)

		require.Equal(t, "On", serverInfo.PowerState, "Server is not powered on")
		require.Equal(t, "OK", serverInfo.Status.Health, "Server health is not OK")
	})

	// t.Run("GetStorageInfo", func(t *testing.T) {
	// 	_, err := client.GetStorageInfo()
	// 	require.NoError(t, err)
	// 	// Add more assertions here based on expected results
	// })

	// t.Run("GetDrivesInfo", func(t *testing.T) {
	// 	drives, err := client.GetDrivesInfo()
	// 	require.NoError(t, err)

	// 	// Assert that there is at least one drive (or as many as expected)
	// 	require.NotEmpty(t, drives, "No drives found")

	// 	// Iterate through each drive and assert expected conditions
	// 	for _, drive := range drives {
	// 		require.NotEmpty(t, drive.ID, "Drive ID is empty")
	// 		require.Equal(t, "OK", drive.Status.Health, "Drive health is not OK")
	// 		// Add more assertions here based on expected results
	// 	}
	// })
	// t.Run("GetRAIDVirtualVolumes", func(t *testing.T) {
	// 	raid_volumes, err := client.GetRAIDVolumeInfo()
	// 	require.NoErro(t, err)
	// })
}

func TestIDRACClientIntegration_GetRAIDInfo(t *testing.T) {
	hostname := os.Getenv("IDRAC_HOSTNAME")
	if hostname == "" {
		t.Skip("IDRAC_HOSTNAME environment variable not set")
	}

	client := NewClient(config.IDRACConfig{
		BMCConnConfig: config.BMCConnConfig{
			Hostname: hostname,
			Username: os.Getenv("IDRAC_USERNAME"),
			Password: os.Getenv("IDRAC_PASSWORD"),
		},
	})

	client.HTTPClientConfig = httpclient.DefaultConfig()

	controllers, err := client.GetRAIDControllers()
	require.NoError(t, err)
	require.NotEmpty(t, controllers, "No RAID controllers found")

	for _, controller := range controllers {
		for _, volumeRef := range controller.Volumes {
			volume, err := client.GetRAIDVolumeInfo(volumeRef.ID)
			require.NoError(t, err)
			require.Equal(t, "OK", volume.Health, "RAID volume health is not OK")
			// Add more assertions as necessary
		}
	}
}

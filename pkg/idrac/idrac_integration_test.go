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

func newTestIDRACClient(t *testing.T) *Client {
	hostname := os.Getenv("BMC_HOSTNAME")
	if hostname == "" {
		t.Skip("BMC_HOSTNAME environment variable not set")
	}

	client := NewClient(config.IDRACConfig{
		BMCConnConfig: config.BMCConnConfig{
			Hostname: hostname,
			Username: os.Getenv("BMC_USERNAME"),
			Password: os.Getenv("BMC_PASSWORD"),
		},
	})
	client.HTTPClientConfig = httpclient.DefaultConfig()

	return client
}

func TestIDRACClientIntegration(t *testing.T) {
	client := newTestIDRACClient(t)

	t.Run("GetServerInfo", func(t *testing.T) {
		serverInfo, err := client.GetServerInfo()
		require.NoError(t, err)
		require.NotNil(t, serverInfo)

		require.Equal(t, "On", serverInfo.PowerState, "Server is not powered on")
		require.Equal(t, "OK", serverInfo.Status.Health, "Server health is not OK")
	})

	t.Run("GetStorageInfo", func(t *testing.T) {
		storageInfo, err := client.GetStorageInfo()
		require.NoError(t, err)
		require.NotNil(t, storageInfo)
	})

	t.Run("GetDrivesInfo", func(t *testing.T) {
		drives, err := client.GetDrivesInfo()
		require.NoError(t, err)
		require.NotEmpty(t, drives, "No drives found")

		// Iterate through each drive and assert expected conditions
		for _, drive := range drives {
			require.NotEmpty(t, drive.ID, "Drive ID is empty")
			require.Equal(t, "OK", drive.Status.Health, "Drive health is not OK")
		}
	})
}

func TestIDRACClientIntegration_GetRAIDInfo(t *testing.T) {
	client := newTestIDRACClient(t)

	t.Run("GetRAIDControllers", func(t *testing.T) {
		controllers, err := client.GetRAIDControllers()
		require.NoError(t, err)
		require.NotEmpty(t, controllers, "No RAID controllers found")

		for _, controller := range controllers {
			require.NotEmpty(t, controller.ID, "RAID controller ID is empty")

			// Add assertions for other fields in the controller if necessary

			t.Run("GetRAIDControllerInfo", func(t *testing.T) {
				controllerInfo, err := client.GetRAIDControllerInfo(controller.ID)
				require.NoError(t, err)
				require.NotNil(t, controllerInfo)

				t.Run("GetRAIDDriveDetails", func(t *testing.T) {
					for _, driveRef := range controllerInfo.Drives {
						drive, err := client.GetRAIDDriveDetails(driveRef.ID)
						require.NoError(t, err)
						require.NotNil(t, drive)
						require.Equal(t, "OK", drive.Status.Health, "RAID drive health is not OK")
					}
				})
			})

			t.Run("GetRAIDVolumeInfo", func(t *testing.T) {
				for _, volumeRef := range controller.Volumes {
					volume, err := client.GetRAIDVolumeInfo(volumeRef.ID)
					require.NoError(t, err)
					require.NotNil(t, volume)
					require.Equal(t, "OK", volume.Status.Health, "RAID volume health is not OK")
				}
			})
		}
	})
}

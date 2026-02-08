package idrac

import (
	"fmt"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/angelhvargas/redfishcli/pkg/request"
)

// Client represents an iDRAC client.
type Client struct {
	Config           config.IDRACConfig
	HTTPClientConfig httpclient.Config
}

// NewClient initializes a new iDRAC client with default HTTP client configuration.
func NewClient(cfg config.IDRACConfig) *Client {
	return &Client{
		Config:           cfg,
		HTTPClientConfig: httpclient.DefaultConfig(),
	}
}

// GetServerInfo retrieves the server information from iDRAC.
func (c *Client) GetServerInfo() (*model.ServerInfo, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1", c.Config.Hostname)
	var info model.ServerInfo
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetStorageInfo retrieves the storage information from iDRAC.
func (c *Client) GetStorageInfo() (*model.StorageInfo, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	var info model.StorageInfo
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetDrivesInfo retrieves information for all drives from iDRAC.
func (c *Client) GetDrivesInfo() ([]model.Drive, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	var storageCollection model.StorageCollection
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &storageCollection); err != nil {
		return nil, err
	}

	var drives []model.Drive
	for _, member := range storageCollection.Members {
		var storage model.Storage
		if err := request.FetchAndUnmarshal(member.ID, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &storage); err != nil {
			return nil, err
		}

		for _, driveRef := range storage.Drives {
			var drive model.Drive
			if err := request.FetchAndUnmarshal(driveRef.ID, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &drive); err != nil {
				return nil, err
			}
			drives = append(drives, drive)
		}
	}

	return drives, nil
}

// GetStorageControllers retrieves RAID controller information from iDRAC.
func (c *Client) GetStorageControllers(config *model.StorageControllerConfig) ([]model.StorageController, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	var storageResp model.StorageCollection
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &storageResp); err != nil {
		return nil, err
	}

	var StorageControllers []model.StorageController
	for _, member := range storageResp.Members {
		StorageControllers = append(StorageControllers, model.StorageController{ID: member.ID})
	}

	return StorageControllers, nil
}

// GetStorageControllerInfo retrieves detailed information for a specific RAID controller.
func (c *Client) GetStorageControllerInfo(controllerID string) (*model.StorageControllerDetails, error) {
	url := fmt.Sprintf("https://%s%s", c.Config.Hostname, controllerID)
	var raidControllerDetails model.StorageControllerDetails
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &raidControllerDetails); err != nil {
		return nil, err
	}
	return &raidControllerDetails, nil
}

// GetRAIDVolumeInfo retrieves information for a specific RAID volume.
func (c *Client) GetRAIDVolumeInfo(volumeEndpoint string) (*model.RAIDVolume, error) {
	var volume model.RAIDVolume
	if err := request.FetchAndUnmarshal(volumeEndpoint, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &volume); err != nil {
		return nil, err
	}
	return &volume, nil
}

// GetStorageDriveDetails retrieves detailed information for a specific drive.
func (c *Client) GetStorageDriveDetails(driveURL string) (*model.Drive, error) {
	var drive model.Drive
	url := fmt.Sprintf("https://%s%s", c.Config.Hostname, driveURL)
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &drive); err != nil {
		return nil, err
	}
	return &drive, nil
}

// GetPowerState retrieves the current power state of the server.
func (c *Client) GetPowerState() (string, error) {
	info, err := c.GetServerInfo()
	if err != nil {
		return "", err
	}
	return info.PowerState, nil
}

// SetPowerState sets the power state of the server (On, ForceOff, GracefulShutdown).
func (c *Client) SetPowerState(state string) error {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset", c.Config.Hostname)
	payload := map[string]string{
		"ResetType": state,
	}
	return request.Post(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, payload)
}

// Reboot reboots the server (GracefulRestart).
func (c *Client) Reboot() error {
	return c.SetPowerState("GracefulRestart")
}

// GetBootInfo retrieves the boot information.
func (c *Client) GetBootInfo() (*model.BootInfo, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1", c.Config.Hostname)
	var info model.BootInfo
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// SetBootOrder sets the boot order (e.g., PxE, Hdd, Cd).
func (c *Client) SetBootOrder(device string) error {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1", c.Config.Hostname)
	payload := map[string]interface{}{
		"Boot": map[string]string{
			"BootSourceOverrideTarget": device,
		},
	}
	return request.Post(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, payload)
}

// GetSystemEventLog retrieves the system event log.
func (c *Client) GetSystemEventLog() ([]model.EventLogEntry, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Managers/iDRAC.Embedded.1/LogServices/ Sel/Entries", c.Config.Hostname)
	var log model.EventLog
	if err := request.FetchAndUnmarshal(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig, &log); err != nil {
		return nil, err
	}
	return log.Members, nil
}

func init() {
	client.Register("idrac", func(cfg config.BMCConnConfig) client.ServerClient {
		return NewClient(config.IDRACConfig{BMCConnConfig: cfg})
	})
}

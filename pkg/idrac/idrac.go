package idrac

import (
	"encoding/json"
	"fmt"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/angelhvargas/redfishcli/pkg/model"
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

// fetchAndUnmarshal performs a HTTP GET request to the specified Endpoint and unmarshals the response into the given target structure.
func (c *Client) fetchAndUnmarshal(url string, target interface{}) error {
	body, err := httpclient.DoRequest(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		logger.Log.Errorf("Error fetching data: %s", err)
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		logger.Log.Errorf("Error unmarshalling data: %s", err)
		return err
	}

	return nil
}

// GetServerInfo retrieves the server information from iDRAC.
func (c *Client) GetServerInfo() (*model.ServerInfo, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1", c.Config.Hostname)
	var info model.ServerInfo
	if err := c.fetchAndUnmarshal(url, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetStorageInfo retrieves the storage information from iDRAC.
func (c *Client) GetStorageInfo() (*model.StorageInfo, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	var info model.StorageInfo
	if err := c.fetchAndUnmarshal(url, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetDrivesInfo retrieves information for all drives from iDRAC.
func (c *Client) GetDrivesInfo() ([]model.Drive, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	var storageCollection model.StorageCollection
	if err := c.fetchAndUnmarshal(url, &storageCollection); err != nil {
		return nil, err
	}

	var drives []model.Drive
	for _, member := range storageCollection.Members {
		var storage model.Storage
		if err := c.fetchAndUnmarshal(member.ID, &storage); err != nil {
			return nil, err
		}

		for _, driveRef := range storage.Drives {
			var drive model.Drive
			if err := c.fetchAndUnmarshal(driveRef.ID, &drive); err != nil {
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
	if err := c.fetchAndUnmarshal(url, &storageResp); err != nil {
		return nil, err
	}

	var StorageControllers []model.StorageController
	for _, member := range storageResp.Members {
		StorageControllers = append(StorageControllers, model.StorageController{ID: member.ID})
	}

	return StorageControllers, nil
}

// Additional functions (GetStorageControllerInfo, GetRAIDVolumeInfo, GetStorageDriveDetails) follow the same pattern.
// GetStorageControllerInfo retrieves detailed information for a specific RAID controller.
func (c *Client) GetStorageControllerInfo(controllerId string) (*model.StorageControllerDetails, error) {
	url := fmt.Sprintf("https://%s%s", c.Config.Hostname, controllerId)
	var raidControllerDetails model.StorageControllerDetails
	if err := c.fetchAndUnmarshal(url, &raidControllerDetails); err != nil {
		return nil, err
	}
	return &raidControllerDetails, nil
}

// GetRAIDVolumeInfo retrieves information for a specific RAID volume.
func (c *Client) GetRAIDVolumeInfo(volumeEndpoint string) (*model.RAIDVolume, error) {
	var volume model.RAIDVolume
	if err := c.fetchAndUnmarshal(volumeEndpoint, &volume); err != nil {
		return nil, err
	}
	return &volume, nil
}

// GetStorageDriveDetails retrieves detailed information for a specific drive.
func (c *Client) GetStorageDriveDetails(driveUrl string) (*model.Drive, error) {
	var drive model.Drive
	url := fmt.Sprintf("https://%s%s", c.Config.Hostname, driveUrl)
	if err := c.fetchAndUnmarshal(url, &drive); err != nil {
		return nil, err
	}
	return &drive, nil
}

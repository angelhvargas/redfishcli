package idrac

import (
	"encoding/json"
	"fmt"
	"strings"

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

// fetchAndUnmarshal performs a HTTP GET request to the specified URL and unmarshals the response into the given target structure.
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

// GetRAIDControllers retrieves RAID controller information from iDRAC.
func (c *Client) GetRAIDControllers() ([]model.RAIDController, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	var storageResp model.StorageCollection
	if err := c.fetchAndUnmarshal(url, &storageResp); err != nil {
		return nil, err
	}

	var RAIDControllers []model.RAIDController
	for _, member := range storageResp.Members {
		if strings.Contains(member.ID, "RAID") {
			RAIDControllers = append(RAIDControllers, model.RAIDController{ID: member.ID})
		}
	}

	return RAIDControllers, nil
}

// Additional functions (GetRAIDControllerInfo, GetRAIDVolumeInfo, GetRAIDDriveDetails) follow the same pattern.
// GetRAIDControllerInfo retrieves detailed information for a specific RAID controller.
func (c *Client) GetRAIDControllerInfo(controllerId string) (*model.RAIDControllerDetails, error) {
	url := fmt.Sprintf("https://%s%s", c.Config.Hostname, controllerId)
	var raidControllerDetails model.RAIDControllerDetails
	if err := c.fetchAndUnmarshal(url, &raidControllerDetails); err != nil {
		return nil, err
	}
	return &raidControllerDetails, nil
}

// GetRAIDVolumeInfo retrieves information for a specific RAID volume.
func (c *Client) GetRAIDVolumeInfo(volumeURL string) (*model.RAIDVolume, error) {
	var volume model.RAIDVolume
	if err := c.fetchAndUnmarshal(volumeURL, &volume); err != nil {
		return nil, err
	}
	return &volume, nil
}

// GetRAIDDriveDetails retrieves detailed information for a specific drive.
func (c *Client) GetRAIDDriveDetails(driveUrl string) (*model.Drive, error) {
	var drive model.Drive
	url := fmt.Sprintf("https://%s%s", c.Config.Hostname, driveUrl)
	if err := c.fetchAndUnmarshal(url, &drive); err != nil {
		return nil, err
	}
	return &drive, nil
}

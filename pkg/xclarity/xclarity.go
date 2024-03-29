package xclarity

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/angelhvargas/redfishcli/pkg/model"
)

type Client struct {
	Config           config.XClarityConfig
	Debug            bool
	HTTPClientConfig httpclient.Config
}

// NewClient creates a new iDRAC client
func NewClient(cfg config.XClarityConfig) *Client {
	return &Client{
		Config:           cfg,
		Debug:            false,
		HTTPClientConfig: httpclient.DefaultConfig(),
	}
}

// GetServerInfo gets the server information from iDRAC
func (c *Client) GetServerInfo() (*model.ServerInfo, error) {
	// Construct the URL
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1", c.Config.Hostname)

	// Create a new request
	log.Printf("Hostname: %s, Username: %s", c.Config.Hostname, c.Config.Username)

	body, err := httpclient.DoRequest(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		return nil, err
	}

	log.Printf("Raw server info response: %s\n", string(body))
	var info model.ServerInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Client) GetStorageInfo() (*model.StorageInfo, error) {
	// Example URL for storage information
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	log.Printf("Hostname: %s, Username: %s", c.Config.Hostname, c.Config.Username)

	body, err := httpclient.DoRequest(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		return nil, err
	}

	log.Printf("Raw server info response: %s\n", string(body))
	var info model.StorageInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Client) GetDrivesInfo() ([]model.Drive, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)

	storageCollection, err := c.fetchStorageCollection(url)
	if err != nil {
		return nil, err
	}

	var drives []model.Drive
	for _, member := range storageCollection.Members {
		storage, err := c.fetchStorage(member.ID)
		if err != nil {
			return nil, err
		}

		for _, driveRef := range storage.Drives {
			drive, err := c.fetchDrive(driveRef.ID)
			if err != nil {
				return nil, err
			}
			drives = append(drives, *drive)
		}
	}

	return drives, nil
}

func (c *Client) fetchStorageCollection(url string) (*model.StorageCollection, error) {
	log.Printf("Hostname: %s, Username: %s", c.Config.Hostname, c.Config.Username)

	body, err := httpclient.DoRequest(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		return nil, err
	}

	var storageCollection model.StorageCollection
	err = json.Unmarshal(body, &storageCollection)
	if err != nil {
		return nil, err
	}

	return &storageCollection, nil
}

func (c *Client) fetchStorage(url string) (*model.Storage, error) {
	log.Printf("Hostname: %s, Username: %s", c.Config.Hostname, c.Config.Username)
	redfish_url := fmt.Sprintf("https://%s%s", c.Config.Hostname, url)
	log.Printf("Redfish API Url: %s", redfish_url)

	body, err := httpclient.DoRequest(redfish_url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		return nil, err
	}

	log.Printf("Raw server info response: %s\n", string(body))
	var storage model.Storage
	err = json.Unmarshal(body, &storage)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}

func (c *Client) fetchDrive(url string) (*model.Drive, error) {
	log.Printf("Hostname: %s, Username: %s", c.Config.Hostname, c.Config.Username)
	redfish_url := fmt.Sprintf("https://%s%s", c.Config.Hostname, url)
	log.Printf("Redfish API Url: %s", redfish_url)

	body, err := httpclient.DoRequest(redfish_url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		return nil, err
	}
	log.Printf("Raw server info response: %s\n", string(body))

	var drive model.Drive
	err = json.Unmarshal(body, &drive)
	if err != nil {
		return nil, err
	}

	return &drive, nil
}

func (c *Client) GetRAIDControllers() ([]model.RAIDController, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)
	logger.Log.Printf(url)
	body, err := httpclient.DoRequest(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		logger.Log.Errorf(err.Error())
		return nil, err
	}
	log.Printf("Raw server info response: %s\n", string(body))
	var storageResp model.StorageResponse
	if err := json.Unmarshal(body, &storageResp); err != nil {
		return nil, err
	}

	return storageResp.Members, nil
}

func (c *Client) GetRAIDVolumeInfo(volumeURL string) (*model.RAIDVolume, error) {
	body, err := httpclient.DoRequest(volumeURL, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		return nil, err
	}
	log.Printf("Raw server info response: %s\n", string(body))
	var volume model.RAIDVolume
	if err := json.Unmarshal(body, &volume); err != nil {
		return nil, err
	}

	return &volume, nil
}

func (c *Client) GetRAIDControllerInfo(controllerId string) (*model.RAIDControllerDetails, error) {
	url := fmt.Sprintf("https://%s", controllerId)
	logger.Log.Info("calling redfish:", url)
	body, err := httpclient.DoRequest(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	var raidControllerDetails model.RAIDControllerDetails
	if err := json.Unmarshal(body, &raidControllerDetails); err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return &raidControllerDetails, nil

}

func (c *Client) GetRAIDDriveDetails(driveUrl string) (*model.Drive, error) {
	url := fmt.Sprint("http://%%", c.Config.Hostname, driveUrl)
	body, err := httpclient.DoRequest(url, c.Config.Username, c.Config.Password, c.HTTPClientConfig)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	var drive model.Drive
	if err := json.Unmarshal(body, &drive); err != nil {
		logger.Log.Error(err.Error())
	}
	return &drive, nil
}

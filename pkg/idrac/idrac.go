// client/idrac/idrac.go

package idrac

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"log"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/model"
)

// Client represents an iDRAC client
type Client struct {
	Config config.IDRACConfig
	// You can include a http.Client here if you want to customize it (e.g., timeouts)
}

// NewClient creates a new iDRAC client
func NewClient(cfg config.IDRACConfig) *Client {
	return &Client{
		Config: cfg,
		// Initialize http.Client here if needed
	}
}

// GetServerInfo gets the server information from iDRAC
func (c *Client) GetServerInfo() (*model.ServerInfo, error) {
	log.Println("Making HTTP request to get server info")
	// Construct the URL
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1", c.Config.Hostname)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set basic auth
	req.SetBasicAuth(c.Config.Username, c.Config.Password)

	// Perform the request
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info model.ServerInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Client) GetStorageInfo() (*model.StorageInfo, error) {
	log.Println("Making HTTP request to get storage info")
	// Example URL for storage information
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Storage", c.Config.Hostname)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Config.Username, c.Config.Password)

	httpClient := &http.Client{Timeout: time.Second * 30}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Use io.ReadAll instead of ioutil.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info model.StorageInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	// Additional logic might be required here to adapt to different iDRAC versions
	// ...

	return &info, nil
}

// client/idrac/idrac.go

package idrac

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
	body, err := ioutil.ReadAll(resp.Body)
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

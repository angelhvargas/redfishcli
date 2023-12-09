package redfish

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RedfishClient struct {
	BaseURL  string
	Username string
	Password string
	// Add any other configurations like HTTP client
}

type ServerInfo struct {
	SerialNumber string
	PowerStatus  string
	Health       string
}

func (c *RedfishClient) GetServerInfo() (*ServerInfo, error) {
	resp, err := c.makeRequest("/redfish/v1/Systems/System.Embedded.1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", resp.Status)
	}

	var serverInfo ServerInfo
	err = json.NewDecoder(resp.Body).Decode(&serverInfo)
	if err != nil {
		return nil, err
	}

	return &serverInfo, nil
}

func (c *RedfishClient) buildURL(endpoint string) string {
	return c.BaseURL + endpoint
}

func (c *RedfishClient) makeRequest(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.buildURL(endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)

	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	return httpClient.Do(req)
}

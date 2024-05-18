package httpclient

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/angelhvargas/redfishcli/pkg/logger"
)

type Config struct {
	Timeout       time.Duration
	SkipTLSVerify bool
}

// DefaultConfig provides default settings for the HTTP client.
func DefaultConfig() Config {
	return Config{
		Timeout:       60 * time.Second,
		SkipTLSVerify: true,
	}
}

// Function variable for mocking purposes
var DoRequest = doRequest

// doRequest performs an HTTP GET request.
func doRequest(url, username, password string, config Config) ([]byte, error) {
	logger.Log.Printf("Doing http request to: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)

	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipTLSVerify},
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Log.Errorf("Error: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("server returned status code %d", resp.StatusCode)
		logger.Log.Errorf("Error: %s", err)
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

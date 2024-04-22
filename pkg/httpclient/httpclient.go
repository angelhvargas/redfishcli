// pkg/httpclient/httpclient.go

package httpclient

import (
	"crypto/tls"
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
		Timeout:       30 * time.Second,
		SkipTLSVerify: true,
	}
}

func DoRequest(url, username, password string, config Config) ([]byte, error) {
	logger.Log.Printf("Doing http request to: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)

	httpClient := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Log.Errorf("Error: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

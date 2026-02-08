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
		Timeout:       30 * time.Second,
		SkipTLSVerify: true,
	}
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

var (
	ErrAuthentication = &HTTPError{StatusCode: 401, Message: "authentication error"}
	ErrAuthorization  = &HTTPError{StatusCode: 403, Message: "authorization error"}
	ErrNotFound       = &HTTPError{StatusCode: 404, Message: "endpoint not found"}
)

// Function variable for mocking purposes
var DoRequest = doRequest
var Do = do

// doRequest performs an HTTP GET request.
func doRequest(url, username, password string, config Config) ([]byte, error) {
	return do("GET", url, username, password, nil, config)
}

// do performs an HTTP request.
func do(method, url, username, password string, body io.Reader, config Config) ([]byte, error) {
	logger.Log.Printf("API request: %s %s", method, url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

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

	logger.Log.Info(resp.StatusCode)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var httpErr error
		switch resp.StatusCode {
		case 401:
			httpErr = ErrAuthentication
		case 403:
			httpErr = ErrAuthorization
		case 404:
			httpErr = ErrNotFound
		default:
			httpErr = &HTTPError{StatusCode: resp.StatusCode, Message: "unexpected error"}
		}
		logger.Log.Errorf("Error: %s", httpErr)
		return nil, httpErr
	}

	return io.ReadAll(resp.Body)
}

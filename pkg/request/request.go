package request

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/angelhvargas/redfishcli/pkg/logger"
)

// FetchAndUnmarshal performs an HTTP GET request to the specified endpoint and unmarshals the response into the given target structure.
func FetchAndUnmarshal(url, username, password string, config httpclient.Config, target interface{}) error {
	body, err := httpclient.DoRequest(url, username, password, config)
	if err != nil {
		logger.Log.Errorf("Error fetching data: %s", err)
		return HandleHTTPError(err, url)
	}

	if err := json.Unmarshal(body, target); err != nil {
		logger.Log.Errorf("Error unmarshalling data: %s", err)
		return err
	}

	return nil
}

// HandleHTTPError categorizes HTTP errors and returns appropriate error messages
func HandleHTTPError(err error, url string) error {
	if httpErr, ok := err.(*httpclient.HTTPError); ok {
		switch httpErr.StatusCode {
		case 401:
			return fmt.Errorf("url %s: authentication error - %v", url, err)
		case 403:
			return fmt.Errorf("url %s: authorization error - %v", url, err)
		case 404:
			return fmt.Errorf("url %s: endpoint not found - %v", url, err)
		default:
			return fmt.Errorf("url %s: unexpected error - %v", url, err)
		}
	}
	return fmt.Errorf("url %s: unknown error - %v", url, err)
}

// Post performs an HTTP POST request with a JSON payload.
func Post(url, username, password string, config httpclient.Config, payload interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %v", err)
	}

	_, err = httpclient.Do("POST", url, username, password, bytes.NewBuffer(jsonPayload), config)
	if err != nil {
		logger.Log.Errorf("Error posting data: %s", err)
		return HandleHTTPError(err, url)
	}

	return nil
}

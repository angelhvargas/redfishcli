package idrac

import (
	"errors"
	"testing"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/stretchr/testify/assert"
)

func mockDoRequest(mockFunc func(url, username, password string, config httpclient.Config) ([]byte, error)) {
	httpclient.DoRequest = mockFunc
}

func restoreDoRequest() {
	httpclient.DoRequest = httpclient.DoRequest
}

func TestFetchAndUnmarshal(t *testing.T) {
	client := NewClient(config.IDRACConfig{
		BMCConnConfig: config.BMCConnConfig{
			Hostname: "testhost",
			Username: "user",
			Password: "pass",
		},
	})

	defer restoreDoRequest()

	t.Run("success case", func(t *testing.T) {
		mockResponse := []byte(`{"ID": "test-server"}`)
		mockDoRequest(func(url, username, password string, config httpclient.Config) ([]byte, error) {
			return mockResponse, nil
		})

		var result model.ServerInfo
		err := client.fetchAndUnmarshal("https://testhost/redfish/v1/Systems/System.Embedded.1", &result)

		assert.NoError(t, err)
		assert.Equal(t, "test-server", result.ID)
	})

	t.Run("error in request", func(t *testing.T) {
		mockDoRequest(func(url, username, password string, config httpclient.Config) ([]byte, error) {
			return nil, errors.New("request error")
		})

		var result model.ServerInfo
		err := client.fetchAndUnmarshal("https://testhost/redfish/v1/Systems/System.Embedded.1", &result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request error")
	})

	t.Run("error in unmarshal", func(t *testing.T) {
		mockResponse := []byte(`invalid json`)
		mockDoRequest(func(url, username, password string, config httpclient.Config) ([]byte, error) {
			return mockResponse, nil
		})

		var result model.ServerInfo
		err := client.fetchAndUnmarshal("https://testhost/redfish/v1/Systems/System.Embedded.1", &result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character 'i' looking for beginning of value")
	})
}

func TestGetServerInfo(t *testing.T) {
	client := NewClient(config.IDRACConfig{
		BMCConnConfig: config.BMCConnConfig{
			Hostname: "testhost",
			Username: "user",
			Password: "pass",
		},
	})

	defer restoreDoRequest()

	t.Run("success case", func(t *testing.T) {
		mockResponse := []byte(`{"ID": "test-server"}`)
		mockDoRequest(func(url, username, password string, config httpclient.Config) ([]byte, error) {
			return mockResponse, nil
		})

		result, err := client.GetServerInfo()

		assert.NoError(t, err)
		assert.Equal(t, "test-server", result.ID)
	})

	t.Run("error in fetchAndUnmarshal", func(t *testing.T) {
		mockDoRequest(func(url, username, password string, config httpclient.Config) ([]byte, error) {
			return nil, errors.New("request error")
		})

		result, err := client.GetServerInfo()

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

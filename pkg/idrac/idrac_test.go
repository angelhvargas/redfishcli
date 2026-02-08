package idrac

import (
	"errors"
	"testing"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/stretchr/testify/assert"
)

func mockDoRequest(mockFunc func(url, username, password string, config httpclient.Config) ([]byte, error)) {
	httpclient.DoRequest = mockFunc
}

func restoreDoRequest() {
	httpclient.DoRequest = httpclient.DoRequest
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

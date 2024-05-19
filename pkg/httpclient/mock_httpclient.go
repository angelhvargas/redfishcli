package httpclient

import (
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock of HTTPClient interface
type MockHTTPClient struct {
	mock.Mock
}

// DoRequest is a mock implementation of DoRequest
func (m *MockHTTPClient) DoRequest(url, username, password string, config Config) ([]byte, error) {
	args := m.Called(url, username, password, config)
	return args.Get(0).([]byte), args.Error(1)
}

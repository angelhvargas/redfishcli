package client

import (
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/stretchr/testify/mock"
)

// MockServerClient is a mock implementation of ServerClient
type MockServerClient struct {
	mock.Mock
}

func (m *MockServerClient) GetServerInfo() (*model.ServerInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServerInfo), args.Error(1)
}

func (m *MockServerClient) GetStorageInfo() (*model.StorageInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StorageInfo), args.Error(1)
}

func (m *MockServerClient) GetStorageControllers(config *model.StorageControllerConfig) ([]model.StorageController, error) {
	args := m.Called(config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.StorageController), args.Error(1)
}

func (m *MockServerClient) GetRAIDVolumeInfo(volumeEndpoint string) (*model.RAIDVolume, error) {
	args := m.Called(volumeEndpoint)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RAIDVolume), args.Error(1)
}

func (m *MockServerClient) GetStorageControllerInfo(endpoint string) (*model.StorageControllerDetails, error) {
	args := m.Called(endpoint)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StorageControllerDetails), args.Error(1)
}

func (m *MockServerClient) GetStorageDriveDetails(driveEndpoint string) (*model.Drive, error) {
	args := m.Called(driveEndpoint)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Drive), args.Error(1)
}

func (m *MockServerClient) GetPowerState() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockServerClient) SetPowerState(state string) error {
	args := m.Called(state)
	return args.Error(0)
}

func (m *MockServerClient) Reboot() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockServerClient) GetBootInfo() (*model.BootInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BootInfo), args.Error(1)
}

func (m *MockServerClient) SetBootOrder(device string) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockServerClient) GetSystemEventLog() ([]model.EventLogEntry, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.EventLogEntry), args.Error(1)
}

package client

import (
	"fmt"
	"sync"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/model"
)

// ServerClient defines the interface for interacting with BMCs.
type ServerClient interface {
	GetServerInfo() (*model.ServerInfo, error)
	GetStorageInfo() (*model.StorageInfo, error)
	GetStorageControllers(config *model.StorageControllerConfig) ([]model.StorageController, error)
	GetRAIDVolumeInfo(volumeEndpoint string) (*model.RAIDVolume, error)
	GetStorageControllerInfo(endpoint string) (*model.StorageControllerDetails, error)
	GetStorageDriveDetails(driveEndpoint string) (*model.Drive, error)
	// Power Management
	GetPowerState() (string, error)
	SetPowerState(state string) error
	Reboot() error
	// System Information
	GetBootInfo() (*model.BootInfo, error)
	SetBootOrder(device string) error
	GetSystemEventLog() ([]model.EventLogEntry, error)
}

// ClientFactory is a function that creates a new ServerClient.
type ClientFactory func(cfg config.BMCConnConfig) ServerClient

var (
	registryMu sync.RWMutex
	registry   = make(map[string]ClientFactory)
)

// Register registers a new client factory for a specific BMC type.
func Register(bmcType string, factory ClientFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[bmcType] = factory
}

// NewClient creates a new ServerClient based on the BMC type.
func NewClient(bmcType string, cfg config.BMCConnConfig) (ServerClient, error) {
	registryMu.RLock()
	factory, ok := registry[bmcType]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unsupported BMC type: %s", bmcType)
	}

	return factory(cfg), nil
}

// GetSupportedTypes returns a list of supported BMC types.
func GetSupportedTypes() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	var types []string
	for t := range registry {
		types = append(types, t)
	}
	return types
}

// ResetRegistry clears the registry. This is primarily used for testing.
func ResetRegistry() {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = make(map[string]ClientFactory)
}

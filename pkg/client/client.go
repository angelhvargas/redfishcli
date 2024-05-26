package client

import (
	"github.com/angelhvargas/redfishcli/pkg/model"
)

type ServerClient interface {
	GetServerInfo() (*model.ServerInfo, error)
	GetStorageInfo() (*model.StorageInfo, error)
	GetStorageControllers(config *model.StorageControllerConfig) ([]model.StorageController, error)
	GetRAIDVolumeInfo(volumeEndpoint string) (*model.RAIDVolume, error)
	GetStorageControllerInfo(endpoint string) (*model.StorageControllerDetails, error)
	GetStorageDriveDetails(driveEndpoint string) (*model.Drive, error)
	// Other common methods
}

// type BaseClient struct {
// 	Config           config.IDRACConfig
// 	Debug            bool
// 	HTTPClientConfig httpclient.Config
// }

// func (bc *BaseClient) doRequest(url string) ([]byte, error) {
// 	response, err := httpclient.DoRequest(url, bc.Config.Username, bc.Config.Password, bc.HTTPClientConfig)
// 	return response, err
// }

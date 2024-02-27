package client

import (
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/httpclient"
	"github.com/angelhvargas/redfishcli/pkg/model"
)

type ServerClient interface {
	GetServerInfo() (*model.ServerInfo, error)
	GetStorageInfo() (*model.StorageInfo, error)
	GetRAIDControllers() ([]model.RAIDController, error)
	GetRAIDVolumeInfo(volumeURL string) (*model.RAIDVolume, error)
	GetRAIDControllerInfo(RAIDURL string) (*model.RAIDControllerDetails, error)
	GetRAIDDriveDetails(driveURL string) (*model.Drive, error)
	// Other common methods
}

type BaseClient struct {
	Config           config.IDRACConfig
	Debug            bool
	HTTPClientConfig httpclient.Config
}

func (bc *BaseClient) doRequest(url string) ([]byte, error) {
	response, err := httpclient.DoRequest(url, bc.Config.Username, bc.Config.Password, bc.HTTPClientConfig)
	return response, err
}

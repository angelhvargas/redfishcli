package client

import "github.com/angelhvargas/redfishcli/pkg/model"

type ServerClient interface {
	GetServerInfo() (*model.ServerInfo, error)
	GetStorageInfo() (*model.StorageInfo, error)
	// Other common methods
}

type Client struct {
	// Client fields...
}

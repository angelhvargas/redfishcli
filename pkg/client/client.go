package client

import "github.com/angelhvargas/redfishcli/pkg/model"

type ServerClient interface {
	GetServerInfo() (*model.ServerInfo, error)
	// Other common methods
}

package xclarity

import (
	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	// Other imports
)

func NewClient(cfg config.XClarityConfig) client.ServerClient {
	return &Client{Config: cfg}
}

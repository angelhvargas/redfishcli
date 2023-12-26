package idrac

import (
	"testing"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
)

func TestIDRACClientImplementsServerClient(t *testing.T) {
	idracClient := NewClient(config.IDRACConfig{})
	_, ok := interface{}(idracClient).(client.ServerClient)
	if !ok {
		t.Errorf("iDRAC Client does not implement ServerClient")
	}
}

// model/models.go

package model

// ServerInfo represents the basic information of the server
type ServerInfo struct {
	SerialNumber string `json:"SerialNumber"`
	PowerStatus  string `json:"PowerStatus"`
	Health       string `json:"Health"`
	// Add other relevant fields based on the iDRAC API response
}

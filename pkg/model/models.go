// model/models.go

package model

type ServerInfo struct {
	ID           string `json:"ID"`
	SerialNumber string `json:"SerialNumber"`
	PowerStatus  string `json:"PowerStatus"`
	Health       string `json:"Health"`
	// Add additional fields as per the Redfish API response
	Manufacturer string `json:"Manufacturer"`
	Model        string `json:"Model"`
	PowerState   string `json:"PowerState"`
	Status       struct {
		Health string `json:"Health"`
	} `json:"Status"`
}

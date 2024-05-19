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

// RAID health

type RAIDHealthReport struct {
	ID           string  `json:"id" yaml:"id"`
	Name         string  `json:"name" yaml:"name"`
	HealthStatus string  `json:"health_status" yaml:"health_status"`
	State        string  `json:"state" yaml:"state"`
	Drives       []Drive `json:"drives" yaml:"drives"`
	DrivesCount  int8    `json:"drives_count" yaml:"drives_count"`
	Hostname     string  `json:"hostname" yaml:"hostname"`
}

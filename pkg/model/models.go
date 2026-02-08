// model/models.go

package model

type ServerInfo struct {
	ID           string `json:"Id"`
	SerialNumber string `json:"SerialNumber"`
	PowerState   string `json:"PowerState"`
	Status       struct {
		Health string `json:"Health"`
		State  string `json:"State"`
	} `json:"Status"`
	Manufacturer string `json:"Manufacturer"`
	Model        string `json:"Model"`
	SKU          string `json:"SKU"`
	BiosVersion  string `json:"BiosVersion"`
}

type BootInfo struct {
	BootSourceOverrideTarget  string   `json:"BootSourceOverrideTarget"`
	BootSourceOverrideEnabled string   `json:"BootSourceOverrideEnabled"`
	BootOrder                 []string `json:"BootOrder"`
}

type EventLogEntry struct {
	ID        string `json:"Id"`
	Name      string `json:"Name"`
	Created   string `json:"Created"`
	Message   string `json:"Message"`
	Severity  string `json:"Severity"`
	EntryType string `json:"EntryType"`
	Action    string `json:"SensorType"`
}

type EventLog struct {
	Members []EventLogEntry `json:"Members"`
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

type ControllersReport struct {
	Controllers []StorageController
	Hostname    string `json:"hostname" yaml:"hostname"`
}

// model/models.go

package model

type ServerInfo struct {
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

type StorageResponse struct {
	Members []RAIDController `json:"Members"`
}

type StorageInfo struct{}

type StorageCollection struct {
	Members []struct {
		ID string `json:"@odata.id"`
	} `json:"Members"`
}

type Storage struct {
	Drives []struct {
		ID string `json:"@odata.id"`
	} `json:"Drives"`
	// Other fields as necessary
}

type Drive struct {
	ID            string `json:"Id"`
	Description   string `json:"Description"`
	CapacityBytes int64  `json:"CapacityBytes"`
	Status        struct {
		Health string `json:"Health"`
	} `json:"Status"`
	// Other fields as necessary
}

// RAIDController represents a RAID controller
type RAIDController struct {
	ID      string `json:"Id"`
	Volumes []struct {
		ID string `json:"@odata.id"`
	} `json:"Volumes"`
}

// RAIDVolume represents a RAID volume
type RAIDVolume struct {
	ID     string `json:"Id"`
	Name   string `json:"Name"`
	Health string `json:"Health"`
	State  string `json:"State"`
	// Include other relevant fields
}

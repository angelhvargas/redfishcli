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

type StorageInfo struct {
	Id      string `json:"@odata.id"`
	Context string `json:"@odata.context"`
	Members []struct {
		ID string `json:"@odata.id"`
	} `json:"Members"`
}

type StorageCollection struct {
	Id      string `json:"@odata.id"`
	Context string `json:"@odata.context"`
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
	ID      string        `json:"Id"`
	Volumes []OdataObject `json:"Volumes"`
}

type RAIDControllerDetails struct {
	ID                      string               `json:"Id"`
	Name                    string               `json:"Name"`
	Description             string               `json:"Description"`
	Drives                  []OdataObject        `json:"Drives"`
	DrivesCount             int                  `json:"Drives@odata.count"`
	Links                   RAIDControllerLinks  `json:"Links"`
	Status                  RAIDControllerStatus `json:"Status"`
	StorageControllers      []StorageController  `json:"StorageControllers"`
	StorageControllersCount int                  `json:"StorageControllers@odata.count"`
	Volumes                 OdataObject          `json:"Volumes"`
}

type OdataObject struct {
	ID string `json:"@odata.id"`
}

type RAIDControllerLinks struct {
	Enclosures      []OdataObject `json:"Enclosures"`
	EnclosuresCount int           `json:"Enclosures@odata.count"`
}

type RAIDControllerStatus struct {
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
	State        string `json:"State"`
}

type StorageController struct {
	OdataID                      string                 `json:"@odata.id"`
	Assembly                     OdataObject            `json:"Assembly"`
	FirmwareVersion              string                 `json:"FirmwareVersion"`
	Identifiers                  []ControllerIdentifier `json:"Identifiers"`
	Manufacturer                 string                 `json:"Manufacturer"`
	MemberId                     string                 `json:"MemberId"`
	Model                        string                 `json:"Model"`
	Name                         string                 `json:"Name"`
	SpeedGbps                    int                    `json:"SpeedGbps"`
	Status                       RAIDControllerStatus   `json:"Status"`
	SupportedControllerProtocols []string               `json:"SupportedControllerProtocols"`
	SupportedDeviceProtocols     []string               `json:"SupportedDeviceProtocols"`
}

type ControllerIdentifier struct {
	DurableName       string `json:"DurableName"`
	DurableNameFormat string `json:"DurableNameFormat"`
}

type RAIDVolume struct {
	ID                 string        `json:"Id"`
	Name               string        `json:"Name"`
	BlockSizeBytes     int           `json:"BlockSizeBytes"`
	CapacityBytes      int64         `json:"CapacityBytes"`
	Description        string        `json:"Description"`
	Encrypted          bool          `json:"Encrypted"`
	EncryptionTypes    []string      `json:"EncryptionTypes"`
	Identifiers        []interface{} `json:"Identifiers"`
	OptimumIOSizeBytes int           `json:"OptimumIOSizeBytes"`
	VolumeType         string        `json:"VolumeType"`
	Status             struct {
		Health       string `json:"Health"`
		HealthRollup string `json:"HealthRollup"`
		State        string `json:"State"`
	} `json:"Status"`
	Links struct {
		Drives []struct {
			OdataID string `json:"@odata.id"`
		} `json:"Drives"`
		DrivesCount int `json:"Drives@odata.count"`
	} `json:"Links"`
}

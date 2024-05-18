package model

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

// RAIDController represents a RAID controller
type RAIDController struct {
	ID      string        `json:"Id"`
	Volumes []OdataObject `json:"Volumes"`
}

type RAIDControllerDetails struct {
	ID          string        `json:"Id"`
	Name        string        `json:"Name"`
	Description string        `json:"Description"`
	Drives      []OdataObject `json:"Drives"`
	DrivesCount int           `json:"Drives@odata.count"`
	// we need the links to get the drive info.
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

type Drive struct {
	OdataContext                  string        `json:"@odata.context"`
	OdataID                       string        `json:"@odata.id"`
	OdataType                     string        `json:"@odata.type"`
	Actions                       DriveActions  `json:"Actions"`
	Assembly                      OdataObject   `json:"Assembly"`
	BlockSizeBytes                int           `json:"BlockSizeBytes"`
	CapableSpeedGbs               int           `json:"CapableSpeedGbs"`
	CapacityBytes                 int64         `json:"CapacityBytes"`
	Description                   string        `json:"Description"`
	EncryptionAbility             string        `json:"EncryptionAbility"`
	EncryptionStatus              string        `json:"EncryptionStatus"`
	FailurePredicted              bool          `json:"FailurePredicted"`
	HotspareType                  string        `json:"HotspareType"`
	ID                            string        `json:"Id"`
	Identifiers                   []Identifier  `json:"Identifiers"`
	Links                         DriveLinks    `json:"Links"`
	Manufacturer                  string        `json:"Manufacturer"`
	MediaType                     string        `json:"MediaType"`
	Model                         string        `json:"Model"`
	Name                          string        `json:"Name"`
	NegotiatedSpeedGbs            int           `json:"NegotiatedSpeedGbs"`
	Operations                    []interface{} `json:"Operations"` // Adjust type based on expected content
	PartNumber                    string        `json:"PartNumber"`
	PredictedMediaLifeLeftPercent interface{}   `json:"PredictedMediaLifeLeftPercent"` // Adjust type based on expected content
	Protocol                      string        `json:"Protocol"`
	Revision                      string        `json:"Revision"`
	RotationSpeedRPM              int           `json:"RotationSpeedRPM"`
	SerialNumber                  string        `json:"SerialNumber"`
	Status                        DriveStatus   `json:"Status"`
}

type DriveActions struct {
	SecureErase struct {
		Target string `json:"target"`
	} `json:"#Drive.SecureErase"`
}

type Identifier struct {
	DurableName       string `json:"DurableName"`
	DurableNameFormat string `json:"DurableNameFormat"`
}

type DriveLinks struct {
	Chassis      OdataObject   `json:"Chassis"`
	Volumes      []OdataObject `json:"Volumes"`
	VolumesCount int           `json:"Volumes@odata.count"`
}

type DriveStatus struct {
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
	State        string `json:"State"`
}

// RAID health

type RAIDHealthReport struct {
	ID           string  `json:"id"`
	Name         string  `json:"Name"`
	Drives       []Drive `json:"drives"`
	DrivesCount  int8    `json:"DrivesCount"`
	HealthStatus string  `json:"healthstatus"`
	State        string  `json:"enable"`
}

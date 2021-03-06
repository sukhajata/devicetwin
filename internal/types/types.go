package types

// DesiredConfig represents a desired config request
type DesiredConfig struct {
	DeviceEUI  string `json:"deviceEUI"`
	FieldName  string `json:"name"`
	FieldValue string `json:"value"`
}

// UntypedConfigDoc represents a config doc
type UntypedConfigDoc struct {
	Desired  map[string]interface{} `json:"desired"`
	Reported map[string]interface{} `json:"reported"`
}

// ConfigFieldDetails represents config field details
type ConfigFieldDetails struct {
	Index       int32       `json:"i"`
	Name        string      `json:"n"`
	Type        interface{} `json:"t"`
	Default     interface{} `json:"d"`
	Description string      `json:"a"`
	Min         interface{} `json:"b"`
	Max         interface{} `json:"c"`
}

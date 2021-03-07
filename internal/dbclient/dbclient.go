package dbclient

import (
	"github.com/sukhajata/devicetwin/internal/types"
	pb "github.com/sukhajata/ppconfig"
)

// Client represents a database client
type Client interface {
	GetConfigByName(firmware string, fieldDetails types.ConfigFieldDetails, req *pb.GetConfigByNameRequest) (*pb.ConfigField, error)
	GetConfigByIndex(req *pb.GetConfigByIndexRequest) (*pb.ConfigField, error)
	GetFieldDetailsByIndex(index int32, firmwareVersion string, docType string) (types.ConfigFieldDetails, error)
	GetFieldDetailsByName(fieldName string, firmwareVersion string, docType string) (types.ConfigFieldDetails, error)
	UpdateDbDesired(req *pb.SetDesiredRequest, fieldDetails types.ConfigFieldDetails) error
	GetS11ConfigKey(identifier string, slot int32) (string, error)
	UpdateDbReported(req *pb.UpdateReportedRequest, fieldDetails types.ConfigFieldDetails) error
	GetDeviceConfig(req *pb.Identifier) (*pb.ConfigFields, error)
	GetNextRadioOffset() (int, error)
	GetLatestFirmware(docType string) (string, error)
	GetFieldDetails(firmware string, docType string) (map[string]types.ConfigFieldDetails, error)
	UpdateFirmwareAllDevices() error
	UpdateConfigToNewFirmware(identifier string, slot int, configFields map[string]types.ConfigFieldDetails)
	GetDLResmin(identifier string) (string, error)
	GetInconsistentDevices() ([]string, error)
	DeleteConfig(identifier string, slot int) error
}

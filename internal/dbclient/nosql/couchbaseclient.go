package nosql

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/sukhajata/devicetwin.git/internal/types"
	"github.com/sukhajata/devicetwin.git/internal/utility"
	"github.com/sukhajata/devicetwin.git/pkg/db"
	"github.com/sukhajata/devicetwin.git/pkg/loggerhelper"
	pb "github.com/sukhajata/ppconfig"
	pbLogger "github.com/sukhajata/pplogger"
)

const (
	docTypeConnection        = "connection"
	docTypePendingConnection = "pending-connection"
	//docTypeS11Config         = "s11config"

	// DocTypeConfigSchema normal config schema
	DocTypeConfigSchema = "configschema"

	// DocTypeS11ConfigSchema s11 config schema
	DocTypeS11ConfigSchema = "s11configschema"
)

//CouchbaseClient a client for accessing couchbase
type CouchbaseClient struct {
	dbEngine         db.NoSQLEngine
	bucketName       string
	bucketNameShared string
	loggerHelper     loggerhelper.Helper
}

// NewCouchbaseClient factory method for creating couchbase client
func NewCouchbaseClient(dbEngine db.NoSQLEngine, bucketName string, bucketNameShared string, loggerHelper loggerhelper.Helper) *CouchbaseClient {
	return &CouchbaseClient{
		dbEngine:         dbEngine,
		bucketName:       bucketName,
		bucketNameShared: bucketNameShared,
		loggerHelper:     loggerHelper,
	}
}

// GetFieldDetailsByIndex get the field details for a given config index
func (c *CouchbaseClient) GetFieldDetailsByIndex(index int32, firmwareVersion string, docType string) (types.ConfigFieldDetails, error) {
	queryString := fmt.Sprintf("SELECT f.i, f.n, f.t, f.d, f.a, f.b, f.c FROM %s "+
		" s UNNEST ppschema f WHERE s.type=$1"+
		" AND f.i = $2"+
		" AND s.ppver = $3", c.bucketNameShared)
	var fieldDetails types.ConfigFieldDetails
	results, err := c.dbEngine.Query(c.bucketNameShared, queryString, []interface{}{docType, index, firmwareVersion})
	if err != nil {
		return fieldDetails, err
	}

	if len(results) >= 1 {
		fmap, ok := results[0].(map[string]interface{})
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to map[string]interface{}, type is %v", results[0], reflect.TypeOf(results[0]))
		}

		index, ok := fmap["i"].(float64)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to float64, type is %v", fmap["i"], reflect.TypeOf(fmap["i"]))
		}

		fieldDetails = types.ConfigFieldDetails{
			Index:       int32(index),
			Name:        fmt.Sprintf("%v", fmap["n"]),
			Type:        fmap["t"],
			Default:     fmap["d"],
			Description: fmt.Sprintf("%v", fmap["a"]),
			Min:         fmap["b"],
			Max:         fmap["c"],
		}

		return fieldDetails, nil
	}

	return fieldDetails, fmt.Errorf("config details not found for Index %v", index)

}

// GetFieldDetailsByName get the field details for a given config name
func (c *CouchbaseClient) GetFieldDetailsByName(fieldName string, firmwareVersion string, docType string) (types.ConfigFieldDetails, error) {
	queryString := fmt.Sprintf("SELECT f.i, f.n, f.t, f.d, f.a, f.b, f.c FROM %s "+
		" s UNNEST ppschema f WHERE s.type=$1"+
		" AND f.n = $2"+
		" AND s.ppver = $3", c.bucketNameShared)
	var fieldDetails types.ConfigFieldDetails
	results, err := c.dbEngine.Query(c.bucketNameShared, queryString, []interface{}{docType, fieldName, firmwareVersion})
	if err != nil {
		return fieldDetails, err
	}

	if len(results) >= 1 {
		fmap, ok := results[0].(map[string]interface{})
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v map[string]interface{}", results[0])
		}

		index, ok := fmap["i"].(float64)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to float, type is %v", fmap["i"], reflect.TypeOf(fmap["i"]))
		}

		fieldDetails = types.ConfigFieldDetails{
			Index:       int32(index),
			Name:        fmt.Sprintf("%v", fmap["n"]),
			Type:        fmap["t"],
			Default:     fmap["d"],
			Description: fmt.Sprintf("%v", fmap["a"]),
			Min:         fmap["b"],
			Max:         fmap["c"],
		}

		return fieldDetails, nil
	}

	return fieldDetails, fmt.Errorf("field name %v not found for firmware %v using: %v", fieldName, firmwareVersion, queryString)
}

// UpdateDbDesired update the dbclient desired config with given details
func (c *CouchbaseClient) UpdateDbDesired(req *pb.SetDesiredRequest, fieldDetails types.ConfigFieldDetails) error {
	value, err := utility.StringToInterface(fieldDetails, req.GetFieldValue())
	if err != nil {
		return err
	}

	key := req.Identifier

	if req.Slot > 0 {
		key, err = c.GetS11ConfigKey(req.Identifier, req.Slot)
		if err != nil {
			return err
		}
	}

	fieldPath := "config.desired." + fieldDetails.Name
	err = c.dbEngine.Update(c.bucketName, key, fieldPath, value)
	if err != nil {
		loggerhelper.WriteToLog("Error updating desired config in dbclient: " + fieldDetails.Name)
		return err
	}

	return nil
}

// GetS11ConfigKey get the key for s11 doc
func (c *CouchbaseClient) GetS11ConfigKey(identifier string, slot int32) (string, error) {
	//find s11 config doc key
	slotName := fmt.Sprintf("s%v", slot)
	queryString := fmt.Sprintf("SELECT connection.slots.%s FROM %s WHERE meta().id = $1", slotName, c.bucketName)
	results, err := c.dbEngine.Query(c.bucketName, queryString, []interface{}{identifier})
	if err != nil {
		return "", err
	}

	if len(results) >= 1 {
		slot, ok := results[0].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("could not convert %v to map[string]interface{}", results[0])
		}
		key := fmt.Sprintf("%v", slot[slotName])
		return key, nil
	}

	return "", errors.New("S11 doc key not found")

}

// UpdateDbReported update reported config value
func (c *CouchbaseClient) UpdateDbReported(req *pb.UpdateReportedRequest, fieldDetails types.ConfigFieldDetails) error {
	value, err := utility.DecodeFieldValue(fieldDetails, req.GetFieldValue())
	if err != nil {
		return err
	}

	fieldPath := "config.reported." + fieldDetails.Name
	loggerhelper.WriteToLog(fmt.Sprintf("Updating %s for %s slot %v", fieldPath, req.DeviceEUI, req.Slot))

	key := req.DeviceEUI

	if req.Slot > 0 {
		key, err = c.GetS11ConfigKey(req.DeviceEUI, req.Slot)
		if err != nil {
			return err
		}
	}

	err = c.dbEngine.Update(c.bucketName, key, fieldPath, value)
	return err
}

//GetDeviceConfig get the config for a device
func (c *CouchbaseClient) GetDeviceConfig(req *pb.Identifier) (*pb.ConfigFields, error) {
	//var configMap map[string]interface{}
	var configDoc *types.UntypedConfigDoc //pb.ConfigDoc
	results := &pb.ConfigFields{}
	var fields []*pb.ConfigField

	key := req.Identifier
	docType := DocTypeConfigSchema

	var err error
	if req.Slot > 0 {
		key, err = c.GetS11ConfigKey(req.Identifier, req.Slot)
		if err != nil {
			return results, err
		}
		docType = DocTypeS11ConfigSchema
	}

	err = c.dbEngine.Lookup(c.bucketName, key, "config", &configDoc)
	if err != nil {
		loggerhelper.WriteToLog("Error getting device config")
		return results, err
	}

	firmware, err := c.GetLatestFirmware(docType)
	if err != nil {
		loggerhelper.WriteToLog("Error getting firmware")
		return results, err
	}
	configData, err := c.GetFieldDetails(firmware, docType)
	if err != nil {
		loggerhelper.WriteToLog("Error getting config data")
		return results, err
	}

	for k := range configDoc.Desired {
		field := &pb.ConfigField{
			Name:        k,
			Index:       configData[k].Index,
			Desired:     utility.GetFormattedValue(configDoc.Desired[k]),
			Reported:    utility.GetFormattedValue(configDoc.Reported[k]),
			FieldType:   utility.GetFormattedValue(configData[k].Type),
			Default:     utility.GetFormattedValue(configData[k].Default),
			Description: configData[k].Description,
			Min:         utility.GetFormattedValue(configData[k].Min),
			Max:         utility.GetFormattedValue(configData[k].Max),
		}
		fields = append(fields, field)
	}
	results.Fields = fields

	return results, nil
}

// GetNextRadioOffset gets an incremental radio offset
func (c *CouchbaseClient) GetNextRadioOffset() (int, error) {
	queryString := fmt.Sprintf("SELECT MAX(config.desired.roffset) AS maxroffset FROM %s WHERE type=$1 and config.desired.roffset != ''", c.bucketName)
	results, err := c.dbEngine.Query(c.bucketName, queryString, []interface{}{docTypeConnection})
	if err != nil {
		return -1, err
	}

	result, ok := results[0].(map[string]interface{})
	if !ok {
		return -1, fmt.Errorf("could not convert %v to map[string]interface{}", results[0])
	}

	if result["maxroffset"] != nil {
		maxRoffsetString := fmt.Sprintf("%v", result["maxroffset"])
		maxRoffset, err := strconv.Atoi(maxRoffsetString)
		if err != nil {
			return -1, err
		}

		roffset := maxRoffset + 10
		if roffset >= 2800 {
			//too big, just get random value
			roffset = rand.Intn(2800)
		}
		return roffset, nil
	}

	//just return random value
	return rand.Intn(2800), nil

}

//GetLatestFirmware get the most recent firmware version
func (c *CouchbaseClient) GetLatestFirmware(docType string) (string, error) {
	queryString := fmt.Sprintf("SELECT ppver FROM %s WHERE type = $1 ORDER BY pporder DESC LIMIT 1", c.bucketNameShared)
	results, err := c.dbEngine.Query(c.bucketNameShared, queryString, []interface{}{docType})
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", errors.New("no firmware version found")
	}

	row, ok := results[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not convert %v to map[string]interface{}", results[0])
	}

	firmware := fmt.Sprintf("%v", row["ppver"])

	return firmware, nil
}

// GetFieldDetails - get all the config field details for a firmware
func (c *CouchbaseClient) GetFieldDetails(firmware string, docType string) (map[string]types.ConfigFieldDetails, error) {
	configData := make(map[string]types.ConfigFieldDetails)

	queryString := fmt.Sprintf("SELECT f.i, f.n, f.t, f.d, f.a, f.b, f.c FROM %s s UNNEST s.ppschema AS f WHERE s.type=$1 "+
		"AND s.ppver=$2 ORDER BY f.n", c.bucketNameShared)
	results, err := c.dbEngine.Query(c.bucketNameShared, queryString, []interface{}{docType, firmware})
	if err != nil {
		return configData, err
	}

	if len(results) == 0 {
		return configData, fmt.Errorf("no details found for %s", firmware)
	}

	for _, v := range results {
		fmap, ok := v.(map[string]interface{})
		if !ok {
			return configData, fmt.Errorf("could not convert %v map[string]interface{}, type is %v", results[0], reflect.TypeOf(results[0]))
		}

		index, ok := fmap["i"].(float64)
		if !ok {
			return configData, fmt.Errorf("could not convert %v to float64, type is %v", fmap["i"], reflect.TypeOf(fmap["i"]))
		}

		fieldDetails := types.ConfigFieldDetails{
			Index:       int32(index),
			Name:        fmt.Sprintf("%v", fmap["n"]),
			Type:        fmap["t"],
			Default:     fmap["d"],
			Description: fmt.Sprintf("%v", fmap["a"]),
			Min:         fmap["b"],
			Max:         fmap["c"],
		}

		configData[fieldDetails.Name] = fieldDetails
	}

	return configData, nil
}

// UpdateFirmwareAllDevices update all devices to new firmware
func (c *CouchbaseClient) UpdateFirmwareAllDevices() error {
	//get a list of connections
	queryString := fmt.Sprintf("SELECT meta(c).id FROM %s c WHERE c.type = $1 OR c.type = $2", c.bucketName)
	results, err := c.dbEngine.Query(c.bucketName, queryString, []interface{}{docTypeConnection, docTypePendingConnection})
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return errors.New("no devices found to update")
	}

	firmware, err := c.GetLatestFirmware(DocTypeConfigSchema)
	if err != nil {
		return err
	}
	configFields, err := c.GetFieldDetails(firmware, DocTypeConfigSchema)
	if err != nil {
		return err
	}

	loggerhelper.WriteToLog("Updating firmware")

	for _, v := range results {
		row, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("could not convert %v to map[string]interface{}", v)
		}
		identifier := fmt.Sprintf("%s", row["id"])
		c.UpdateConfigToNewFirmware(identifier, 0, configFields)
		time.Sleep(time.Millisecond * 200) //stagger
	}

	return nil

}

// UpdateConfigToNewFirmware - given a connection id and a new config schema, update the connection, copying across existing values
func (c *CouchbaseClient) UpdateConfigToNewFirmware(identifier string, slot int, configFields map[string]types.ConfigFieldDetails) {
	var configDoc *types.UntypedConfigDoc
	err := c.dbEngine.Lookup(c.bucketName, identifier, "config", &configDoc)
	if err != nil {
		c.loggerHelper.LogError("updateConfigToNewFirmware", "error getting config"+err.Error(), pbLogger.ErrorMessage_SEVERE)
		return
	}

	//new config schema
	newConfigDoc := &pb.ConfigDoc{
		Desired:  make(map[string]string),
		Reported: make(map[string]string),
	}

	//copy existing values to new doc
	for _, v := range configFields {
		if val, ok := configDoc.Desired[v.Name]; ok {
			if val != "" {
				newConfigDoc.Desired[v.Name] = fmt.Sprintf("%v", val)
			} else {
				newConfigDoc.Desired[v.Name] = ""
			}

			if configDoc.Reported[v.Name] != "" {
				newConfigDoc.Reported[v.Name] = fmt.Sprintf("%v", configDoc.Reported[v.Name])
			} else {
				newConfigDoc.Reported[v.Name] = ""
			}
		} else {
			newConfigDoc.Desired[v.Name] = ""
			newConfigDoc.Reported[v.Name] = ""
		}

	}

	loggerhelper.WriteToLog(fmt.Sprintf("Updating config for %s", identifier))

	//replace config document
	err = c.dbEngine.Update(c.bucketName, identifier, "config", newConfigDoc)
	if err != nil {
		c.loggerHelper.LogError("updateConfigToNewFirmware", "error updating config doc: "+err.Error(), pbLogger.ErrorMessage_SEVERE)
		return
	}

}

//GetConfigByName get config details by name
func (c *CouchbaseClient) GetConfigByName(firmware string, fieldDetails types.ConfigFieldDetails, req *pb.GetConfigByNameRequest) (*pb.ConfigField, error) {
	key := req.Identifier

	var err error
	if req.Slot > 0 {
		key, err = c.GetS11ConfigKey(req.Identifier, req.Slot)
		if err != nil {
			return nil, err
		}
	}
	//n1ql allows placeholders in the where clause only. For the rest we can use Sprintf
	queryString := fmt.Sprintf("SELECT config.desired.%s AS desired, config.reported.%s AS reported FROM %s c WHERE meta(c).id = $1", req.GetFieldName(), req.GetFieldName(), c.bucketName)
	results, err := c.dbEngine.Query(c.bucketName, queryString, []interface{}{key})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		loggerhelper.WriteToLog(fmt.Sprintf("Field %s not found", req.FieldName))
		return nil, errors.New("field not found")
	}

	row, ok := results[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not convert %v to map[string]interface{}", results[0])
	}

	return &pb.ConfigField{
		Name:        req.GetFieldName(),
		Index:       fieldDetails.Index,
		Desired:     utility.GetFormattedValue(row["desired"]),
		Reported:    utility.GetFormattedValue(row["reported"]),
		FieldType:   fmt.Sprintf("%v", fieldDetails.Type),
		Description: fieldDetails.Description,
		Default:     utility.GetFormattedValue(fieldDetails.Default),
		Min:         utility.GetFormattedValue(fieldDetails.Min),
		Max:         utility.GetFormattedValue(fieldDetails.Max),
	}, nil

}

//GetConfigByIndex get config by index
func (c *CouchbaseClient) GetConfigByIndex(req *pb.GetConfigByIndexRequest) (*pb.ConfigField, error) {
	docType := DocTypeConfigSchema
	if req.Slot != 0 {
		docType = DocTypeS11ConfigSchema
	}
	firmware, err := c.GetLatestFirmware(docType)
	if err != nil {
		c.loggerHelper.LogError("GetConfigByIndex1", fmt.Sprintf("error getting latest firmware version: %v", err.Error()), pbLogger.ErrorMessage_SEVERE)
		return nil, err
	}

	fieldDetails, err := c.GetFieldDetailsByIndex(req.GetIndex(), firmware, docType)
	if err != nil {
		c.loggerHelper.LogError("GetConfigByIndex2", fmt.Sprintf("error getting %v field details for %v, firmware=%v: %v", req.GetIndex(), req.GetIdentifier(), firmware, err.Error()), pbLogger.ErrorMessage_SEVERE)
		return nil, err
	}

	queryString := fmt.Sprintf("SELECT config.desired.%s AS desired, config.reported.%s AS reported FROM %s c WHERE meta(c).id = $1", fieldDetails.Name, fieldDetails.Name, c.bucketName)
	results, err := c.dbEngine.Query(c.bucketName, queryString, []interface{}{req.GetIdentifier()})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		loggerhelper.WriteToLog(fmt.Sprintf("Field %s not found", fieldDetails.Name))
		return nil, errors.New("field not found")
	}

	row, ok := results[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not convert %v to map[string]interface{}", results[0])
	}

	return &pb.ConfigField{
		Name:        fieldDetails.Name,
		Index:       fieldDetails.Index,
		Desired:     utility.GetFormattedValue(row["desired"]),
		Reported:    utility.GetFormattedValue(row["reported"]),
		FieldType:   fmt.Sprintf("%v", fieldDetails.Type),
		Description: fieldDetails.Description,
		Default:     utility.GetFormattedValue(fieldDetails.Default),
		Min:         utility.GetFormattedValue(fieldDetails.Min),
		Max:         utility.GetFormattedValue(fieldDetails.Max),
	}, nil
}

// GetDLResmin get the downlink reserved minutes for a given device
func (c *CouchbaseClient) GetDLResmin(identifier string) (string, error) {
	queryString := fmt.Sprintf("SELECT config.reported.dlresmin AS dlresmin FROM %s c WHERE meta(c).id = $1", c.bucketName)
	results, err := c.dbEngine.Query(c.bucketName, queryString, []interface{}{identifier})
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", errors.New("no value found")
	}

	value, ok := results[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not convert %v to map[string]interface{}", results[0])
	}

	reservedMinutes := fmt.Sprintf("%v", value["dlresmin"]) //eg."6,8"
	return reservedMinutes, nil
}

// GetInconsistentDevices returns an array of devices with inconsistent config
func (c *CouchbaseClient) GetInconsistentDevices() ([]string, error) {
	queryString := fmt.Sprintf("SELECT meta(b).id FROM %s b WHERE b.type = $1 AND b.config.desired != b.config.reported", c.bucketName)
	rows, err := c.dbEngine.Query(c.bucketName, queryString, []interface{}{docTypeConnection})
	if err != nil {
		return nil, err
	}

	inconsistent := make([]string, 0)

	for _, v := range rows {
		row, ok := v.(map[string]interface{})
		if !ok {
			return inconsistent, fmt.Errorf("could not convert %v to map[string]interface{}", v)
		}
		id := row["id"].(string)
		inconsistent = append(inconsistent, id)
	}

	return inconsistent, nil
}

func (c *CouchbaseClient) DeleteConfig(identifier string, slot int) error {
	return nil
}

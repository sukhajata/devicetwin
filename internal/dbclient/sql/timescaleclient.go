package sql

import (
	"errors"
	"fmt"
	"github.com/sukhajata/devicetwin/internal/dbclient/nosql"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/sukhajata/devicetwin/internal/types"
	"github.com/sukhajata/devicetwin/internal/utility"
	"github.com/sukhajata/devicetwin/pkg/db"
	"github.com/sukhajata/devicetwin/pkg/loggerhelper"
	pb "github.com/sukhajata/ppconfig"
	pbLogger "github.com/sukhajata/pplogger"
)

// TimescaleClient implements dbclient.Client
type TimescaleClient struct {
	dbEngine  db.SQLEngine
	errorChan chan *pbLogger.ErrorMessage
}

//NewTimescaleClient - factory method for generating timescale client
func NewTimescaleClient(dbEngine db.SQLEngine, errorChan chan *pbLogger.ErrorMessage) *TimescaleClient {
	return &TimescaleClient{
		dbEngine:  dbEngine,
		errorChan: errorChan,
	}
}

// GetConfigByName get a config field by name
func (t *TimescaleClient) GetConfigByName(firmware string, fieldDetails types.ConfigFieldDetails, req *pb.GetConfigByNameRequest) (*pb.ConfigField, error) {
	key := req.Identifier

	queryString := `SELECT "DESIRED", "REPORTED" FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	results, err := t.dbEngine.Query(queryString, key, req.Slot, req.FieldName)
	if err != nil {
		return nil, err
	}

	desired := ""
	reported := ""

	if len(results) > 0 {
		row, ok := results[0].([]interface{})
		if !ok {
			msg := fmt.Sprintf("could not convert %v to []interface{}", results[0])
			errMsg := &pbLogger.ErrorMessage{
				Service:  "config-service",
				Function: "GetConfigByName",
				Severity: pbLogger.ErrorMessage_FATAL,
				Message:  msg,
			}
			t.errorChan <- errMsg

			return nil, errors.New(msg)
		}

		desired = utility.GetFormattedValue(row[0])
		reported = utility.GetFormattedValue(row[1])
	}

	return &pb.ConfigField{
		Name:        req.GetFieldName(),
		Index:       fieldDetails.Index,
		Desired:     desired,
		Reported:    reported,
		FieldType:   fmt.Sprintf("%v", fieldDetails.Type),
		Description: fieldDetails.Description,
		Default:     utility.GetFormattedValue(fieldDetails.Default),
		Min:         utility.GetFormattedValue(fieldDetails.Min),
		Max:         utility.GetFormattedValue(fieldDetails.Max),
	}, nil
}

// GetS11ConfigKey get the key for s11 config doc for the given slot
func (t *TimescaleClient) GetS11ConfigKey(identifier string, slot int32) (string, error) {
	return "", errors.New("not supported")
	/*find s11 config doc key
	slotName := fmt.Sprintf("s%v", slot)
	queryString := fmt.Sprintf(`SELECT "DATA" -> 'slots' ->> %s FROM "CONNECTIONS_JSON" WHERE "ID" = $1`, slotName)
	results, err := t.dbEngine.Query(queryString, identifier)
	if err != nil {
		return "", err
	}

	if len(results) >= 1 {
		row, ok := results[0].([]interface{})
		if !ok {
			return "", fmt.Errorf("could not convert %v to []interface{}", results[0])
		}

		key := fmt.Sprintf("%v", row[0])
		return key, nil
	}

	return "", errors.New("S11 doc key not found")*/
}

// GetConfigByIndex - get a config field by its index
func (t *TimescaleClient) GetConfigByIndex(req *pb.GetConfigByIndexRequest) (*pb.ConfigField, error) {
	docType := nosql.DocTypeConfigSchema
	if req.Slot != 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}
	firmware, err := t.GetLatestFirmware(docType)
	if err != nil {
		return nil, err
	}

	fieldDetails, err := t.GetFieldDetailsByIndex(req.GetIndex(), firmware, docType)
	if err != nil {
		return nil, err
	}

	queryString := `SELECT "DESIRED", "REPORTED" FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	results, err := t.dbEngine.Query(queryString, req.Identifier, req.Slot, fieldDetails.Name)
	if err != nil {
		errMessage := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetLatestFirmware",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  err.Error(),
		}
		t.errorChan <- errMessage

		return nil, err
	}

	desired := ""
	reported := ""

	if len(results) > 0 {
		row, ok := results[0].([]interface{})
		if !ok {
			msg := fmt.Sprintf("could not convert %v to []interface{}, type is %v", results[0], reflect.TypeOf(results[0]))
			errMessage := &pbLogger.ErrorMessage{
				Service:  "config-service",
				Function: "GetLatestFirmware",
				Severity: pbLogger.ErrorMessage_FATAL,
				Message:  msg,
			}
			t.errorChan <- errMessage

			return nil, errors.New(msg)
		}
		desired = utility.GetFormattedValue(row[0])
		reported = utility.GetFormattedValue(row[1])
	}

	return &pb.ConfigField{
		Name:        fieldDetails.Name,
		Index:       fieldDetails.Index,
		Desired:     desired,
		Reported:    reported,
		FieldType:   fmt.Sprintf("%v", fieldDetails.Type),
		Description: fieldDetails.Description,
		Default:     utility.GetFormattedValue(fieldDetails.Default),
		Min:         utility.GetFormattedValue(fieldDetails.Min),
		Max:         utility.GetFormattedValue(fieldDetails.Max),
	}, nil
}

//GetLatestFirmware - get the latest firmware version
func (t *TimescaleClient) GetLatestFirmware(docType string) (string, error) {
	ppdev := "meter"
	if docType == nosql.DocTypeS11ConfigSchema {
		ppdev = "controller"
	}
	queryString := `SELECT "PPVER" FROM "CONFIG_SCHEMA" WHERE "PPDEV" = $1 ORDER BY "PPORDER" DESC LIMIT 1`
	results, err := t.dbEngine.Query(queryString, ppdev)
	if err != nil {
		errMessage := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetLatestFirmware",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  err.Error(),
		}
		t.errorChan <- errMessage

		return "", err
	}

	if len(results) == 0 {
		return "", errors.New("no firmware version found")
	}

	row, ok := results[0].([]interface{})
	if !ok {
		msg := fmt.Sprintf("could not convert %v to []interface{}, type is %v", results[0], reflect.TypeOf(results[0]))
		errMessage := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetLatestFirmware",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  msg,
		}
		t.errorChan <- errMessage

		return "", errors.New(msg)
	}

	firmware := fmt.Sprintf("%v", row[0])

	return firmware, nil
}

// GetFieldDetailsByIndex - get field details by config index
func (t *TimescaleClient) GetFieldDetailsByIndex(index int32, firmwareVersion string, docType string) (types.ConfigFieldDetails, error) {
	ppdev := "meter"
	if docType == nosql.DocTypeS11ConfigSchema {
		ppdev = "controller"
	}
	queryString := `SELECT "INDEX", "NAME", "TYPE", "DEFAULT", "DESCRIPTION", "MIN", "MAX"
		FROM "CONFIG_SCHEMA"
		WHERE "PPDEV" = $1
		AND "INDEX" = $2
		AND "PPVER" = $3`
	var fieldDetails types.ConfigFieldDetails
	results, err := t.dbEngine.Query(queryString, ppdev, index, firmwareVersion)
	if err != nil {
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetFieldDetailsByIndex",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  err.Error(),
		}
		t.errorChan <- errMsg

		return fieldDetails, err
	}

	if len(results) >= 1 {
		var ok bool
		var row []interface{}
		row, ok = results[0].([]interface{})
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to []interface{}, type is %v", results[0], reflect.TypeOf(results[0]))
		}

		index, ok := row[0].(int32)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert index %v to int32, type is %v", row[0], reflect.TypeOf(row[0]))
		}

		name, ok := row[1].(string)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert name %v to string, type is %v", row[1], reflect.TypeOf(row[1]))
		}

		description, ok := row[4].(string)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert description %v to string, type is %v", row[5], reflect.TypeOf(row[5]))
		}

		fieldDetails := types.ConfigFieldDetails{
			Index:       index,
			Name:        name,
			Type:        row[2],
			Default:     row[3],
			Description: description,
			Min:         row[5],
			Max:         row[6],
		}
		return fieldDetails, nil
	}

	return fieldDetails, fmt.Errorf("config details not found for Index %v", index)
}

// GetFieldDetailsByName - get field details by name
func (t *TimescaleClient) GetFieldDetailsByName(fieldName string, firmwareVersion string, docType string) (types.ConfigFieldDetails, error) {
	ppdev := "meter"
	if docType == nosql.DocTypeS11ConfigSchema {
		ppdev = "controller"
	}
	queryString := `SELECT "INDEX", "NAME", "TYPE", "DEFAULT", "DESCRIPTION", "MIN", "MAX"
		FROM "CONFIG_SCHEMA"
		WHERE "PPDEV" = $1
		AND "NAME" = $2
		AND "PPVER" = $3`
	var fieldDetails types.ConfigFieldDetails
	results, err := t.dbEngine.Query(queryString, ppdev, fieldName, firmwareVersion)
	if err != nil {
		return fieldDetails, err
	}

	if len(results) >= 1 {
		var ok bool
		var row []interface{}
		row, ok = results[0].([]interface{})
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to []interface{}", results[0])
		}

		index, ok := row[0].(int32)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to int32, type is %v", row[0], reflect.TypeOf(row[0]))
		}

		name, ok := row[1].(string)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to string", row[1])
		}

		description, ok := row[4].(string)
		if !ok {
			return fieldDetails, fmt.Errorf("could not convert %v to string", row[5])
		}

		fieldDetails := types.ConfigFieldDetails{
			Index:       int32(index),
			Name:        name,
			Type:        row[2],
			Default:     row[4],
			Description: description,
			Min:         row[5],
			Max:         row[6],
		}
		return fieldDetails, nil
	}

	return fieldDetails, fmt.Errorf("field name %v not found for firmware %v using: %v", fieldName, firmwareVersion, queryString)
}

// UpdateDbDesired - update the desired value for a config field
func (t *TimescaleClient) UpdateDbDesired(req *pb.SetDesiredRequest, fieldDetails types.ConfigFieldDetails) error {
	value, err := utility.StringToInterface(fieldDetails, req.GetFieldValue())
	if err != nil {
		return err
	}

	key := req.Identifier

	queryString := `SELECT 1 FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	rows, err := t.dbEngine.Query(queryString, key, req.Slot, fieldDetails.Name)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		queryString = `INSERT INTO "CONFIG" ("CONNECTIONID", "SLOT", "NAME", "DESIRED", "REPORTED") VALUES($1, $2, $3, $4, $5)`
		err = t.dbEngine.Exec(queryString, key, req.Slot, fieldDetails.Name, fmt.Sprintf("%v", value), "")
	} else {
		queryString := `UPDATE "CONFIG" SET "DESIRED" = $1 WHERE "CONNECTIONID" = $2 AND "SLOT" = $3 AND "NAME" = $4`
		err = t.dbEngine.Exec(queryString, fmt.Sprintf("%v", value), key, req.Slot, fieldDetails.Name)
	}

	if err != nil {
		msg := fmt.Sprintf("Error updating desired %s for %s: %v", fieldDetails.Name, key, err)
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "UpdateDbDesired",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  msg,
		}
		t.errorChan <- errMsg

		return err
	}

	return nil
}

// UpdateDbReported - update the reported value for a config field
func (t *TimescaleClient) UpdateDbReported(req *pb.UpdateReportedRequest, fieldDetails types.ConfigFieldDetails) error {
	value, err := utility.DecodeFieldValue(fieldDetails, req.GetFieldValue())
	if err != nil {
		return err
	}

	fieldPath := "config.reported." + fieldDetails.Name
	loggerhelper.WriteToLog(fmt.Sprintf("Updating %s for %s slot %v", fieldPath, req.DeviceEUI, req.Slot))

	key := req.DeviceEUI

	queryString := `SELECT 1 FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	rows, err := t.dbEngine.Query(queryString, key, req.Slot, fieldDetails.Name)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		queryString = `INSERT INTO "CONFIG" ("CONNECTIONID", "SLOT", "NAME", "DESIRED", "REPORTED") VALUES($1, $2, $3, $4, $5)`
		err = t.dbEngine.Exec(queryString, key, req.Slot, fieldDetails.Name, "", fmt.Sprintf("%v", value))
	} else {
		queryString = `UPDATE "CONFIG" SET "REPORTED" = $1 WHERE "CONNECTIONID" = $2 AND "SLOT" = $3 AND "NAME" = $4`
		err = t.dbEngine.Exec(queryString, fmt.Sprintf("%v", value), key, req.Slot, fieldDetails.Name)
	}

	if err != nil {
		msg := fmt.Sprintf("Error updating reported %s for %s: %v", fieldDetails.Name, key, err)
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "UpdateDbDesired",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  msg,
		}
		t.errorChan <- errMsg

		return err
	}

	return nil
}

// GetDeviceConfig - get all config fields for a device
func (t *TimescaleClient) GetDeviceConfig(req *pb.Identifier) (*pb.ConfigFields, error) {
	configFields := &pb.ConfigFields{}
	var fields []*pb.ConfigField

	key := req.Identifier
	docType := nosql.DocTypeConfigSchema
	ppdev := "meter"

	var err error
	if req.Slot > 0 {
		ppdev = "controller"
		docType = nosql.DocTypeS11ConfigSchema
	}

	firmware, err := t.GetLatestFirmware(docType)
	if err != nil {
		return configFields, err
	}
	queryString := `SELECT s."NAME", s."INDEX", COALESCE(c."DESIRED", ''), COALESCE(c."REPORTED", ''), s."TYPE", s."DEFAULT", s."DESCRIPTION", s."MIN", s."MAX"
		FROM "CONFIG_SCHEMA" s 
		LEFT JOIN "CONFIG" c
		ON c."NAME" = s."NAME"
		WHERE s."PPVER" = $1
		AND s."PPDEV" = $2
		AND c."CONNECTIONID" = $3
		ORDER BY s."INDEX"`
	results, err := t.dbEngine.Query(queryString, firmware, ppdev, key)
	if err != nil {
		return configFields, err
	}

	for _, v := range results {
		row, ok := v.([]interface{})
		if !ok {
			return configFields, fmt.Errorf("failed to convert %v to []interface{}", v)
		}

		name, ok := row[0].(string)
		if !ok {
			return configFields, fmt.Errorf("failed to convert %v to string", row[0])
		}

		index, ok := row[1].(int32)
		if !ok {
			return configFields, fmt.Errorf("failed to convert %v to int32, type is %v", row[1], reflect.TypeOf(row[1]))
		}

		desired, ok := row[2].(string)
		if !ok {
			return configFields, fmt.Errorf("failed to convert %v to string", row[2])
		}

		reported, ok := row[3].(string)
		if !ok {
			return configFields, fmt.Errorf("failed to convert %v to string", row[3])
		}

		description, ok := row[6].(string)
		if !ok {
			return configFields, fmt.Errorf("failed to convert %v to string", row[6])
		}

		field := &pb.ConfigField{
			Name:        name,
			Index:       index,
			Desired:     utility.GetFormattedValue(desired),
			Reported:    utility.GetFormattedValue(reported),
			FieldType:   utility.GetFormattedValue(row[4]),
			Default:     utility.GetFormattedValue(row[5]),
			Description: description,
			Min:         utility.GetFormattedValue(row[7]),
			Max:         utility.GetFormattedValue(row[8]),
		}
		fields = append(fields, field)
	}
	configFields.Fields = fields

	return configFields, nil
}

func (t *TimescaleClient) GetNextRadioOffset() (int, error) {
	queryString := `SELECT MAX(CAST("DESIRED" AS int)) AS maxroffset FROM "CONFIG" WHERE "NAME" = 'roffset' and "DESIRED" != ''`
	results, err := t.dbEngine.Query(queryString)
	if err != nil {
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetNextRadioOffset",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  err.Error(),
		}
		t.errorChan <- errMsg

		// return random value
		return rand.Intn(2800), nil
	}

	if len(results) == 0 {
		// return random value
		return rand.Intn(2800), nil
	}

	result, ok := results[0].([]interface{})
	if !ok {
		msg := fmt.Sprintf("could not convert %v to []interface{}", results[0])
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetNextRadioOffset",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  msg,
		}
		t.errorChan <- errMsg

		return rand.Intn(2800), nil
	}

	if result[0] != nil {
		maxRoffsetString := fmt.Sprintf("%v", result[0])
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

// GetFieldDetails - get the config schema
func (t *TimescaleClient) GetFieldDetails(firmware string, docType string) (map[string]types.ConfigFieldDetails, error) {
	configData := make(map[string]types.ConfigFieldDetails)
	ppdev := "meter"
	if docType == nosql.DocTypeS11ConfigSchema {
		ppdev = "controller"
	}

	queryString := `SELECT "NAME", "INDEX", "TYPE", "DEFAULT", "DESCRIPTION", "MIN", "MAX"
		FROM "CONFIG_SCHEMA"
		WHERE "PPVER" = $1
		AND "PPDEV" = $2`
	results, err := t.dbEngine.Query(queryString, firmware, ppdev)
	if err != nil {
		return configData, err
	}

	if len(results) == 0 {
		return configData, fmt.Errorf("no details found for %s", firmware)
	}

	for _, v := range results {
		row, ok := v.([]interface{})
		if !ok {
			return configData, fmt.Errorf("could not convert %v to []interface{}", v)
		}
		name, ok := row[0].(string)
		if !ok {
			return configData, fmt.Errorf("failed to convert %v to string", row[0])
		}

		index, ok := row[1].(int32)
		if !ok {
			return configData, fmt.Errorf("failed to convert %v to int32, type is %v", row[1], reflect.TypeOf(row[1]))
		}

		description, ok := row[4].(string)
		if !ok {
			return configData, fmt.Errorf("failed to convert %v to string", row[6])
		}

		field := types.ConfigFieldDetails{
			Index:       index,
			Name:        name,
			Type:        row[2],
			Default:     row[3],
			Description: description,
			Min:         row[5],
			Max:         row[6],
		}
		configData[name] = field
	}

	return configData, nil
}

// UpdateFirmwareAllDevices - set all devices to the latest firmware
func (t *TimescaleClient) UpdateFirmwareAllDevices() error {
	//get a list of connections
	queryString := `SELECT "ID" FROM "CONNECTIONS_JSON"`
	results, err := t.dbEngine.Query(queryString)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return errors.New("no devices found to update")
	}

	//get latest firmware
	firmware, err := t.GetLatestFirmware(nosql.DocTypeConfigSchema)
	if err != nil {
		return err
	}
	//get all config fields for this firmware
	configFields, err := t.GetFieldDetails(firmware, nosql.DocTypeConfigSchema)
	if err != nil {
		return err
	}

	loggerhelper.WriteToLog("Updating firmware")

	for _, v := range results {
		row, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("could not convert %v to []interface{}", v)
		}
		identifier := fmt.Sprintf("%s", row[0])
		t.UpdateConfigToNewFirmware(identifier, 0, configFields)
		time.Sleep(time.Millisecond * 200) //stagger
	}

	return nil
}

// UpdateConfigToNewFirmware - ensure device has an entry for each field in config schema
func (t *TimescaleClient) UpdateConfigToNewFirmware(identifier string, slot int, configFields map[string]types.ConfigFieldDetails) {
	// check if device has an entry for each field in new firmware
	for k := range configFields {
		queryString := `SELECT 1 FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "NAME" = $2 AND "SLOT" = $3`
		results, err := t.dbEngine.Query(queryString, identifier, k, slot)
		if err != nil {
			errMsg := &pbLogger.ErrorMessage{
				Service:  "config-service",
				Function: "UpdateConfigToNewFirmware1",
				Severity: pbLogger.ErrorMessage_FATAL,
				Message:  err.Error(),
			}
			t.errorChan <- errMsg

			return
		}

		if len(results) == 0 {
			queryString = `INSERT INTO "CONFIG" ("CONNECTIONID", "NAME", "SLOT", "DESIRED", "REPORTED") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
			err = t.dbEngine.Exec(queryString, identifier, k, slot, "", "")
			if err != nil {
				errMsg := &pbLogger.ErrorMessage{
					Service:  "config-service",
					Function: "UpdateConfigToNewFirmware2",
					Severity: pbLogger.ErrorMessage_FATAL,
					Message:  err.Error(),
				}
				t.errorChan <- errMsg

				return
			}
		}
	}

	loggerhelper.WriteToLog(fmt.Sprintf("Updating config for %s", identifier))
}

// GetDLResmin - get the downlink reserved minutes for a device
func (t *TimescaleClient) GetDLResmin(identifier string) (string, error) {
	queryString := `SELECT "REPORTED" FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = 0 AND "NAME" = 'dlresmin'`
	results, err := t.dbEngine.Query(queryString, identifier)
	if err != nil {
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetDLResmin",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  err.Error(),
		}
		t.errorChan <- errMsg

		// return default
		return "6,8", nil
	}

	if len(results) == 0 {
		return "6,8", nil
	}

	row, ok := results[0].([]interface{})
	if !ok {
		msg := fmt.Sprintf("could not convert %v to []interface{}", results[0])
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "GetDLResmin",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  msg,
		}
		t.errorChan <- errMsg

		return "6,8", nil
	}

	reservedMinutes := fmt.Sprintf("%v", row[0]) //eg."6,8"
	return reservedMinutes, nil
}

// GetInconsistentDevices returns an array of devices with inconsistent config
func (t *TimescaleClient) GetInconsistentDevices() ([]string, error) {
	queryString := `SELECT DISTINCT "CONNECTIONID" FROM "CONFIG" WHERE "DESIRED" != "REPORTED"`
	rows, err := t.dbEngine.Query(queryString)
	if err != nil {
		return nil, err
	}

	inconsistent := make([]string, 0)

	for _, v := range rows {
		row, ok := v.([]interface{})
		if !ok {
			return inconsistent, fmt.Errorf("could not convert %v to []interface{}", v)
		}
		id := row[0].(string)
		inconsistent = append(inconsistent, id)
	}

	return inconsistent, nil
}

// DeleteConfig - for a device
func (t *TimescaleClient) DeleteConfig(identifier string, slot int) error {
	queryString := `DELETE FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2`
	err := t.dbEngine.Exec(queryString, identifier, slot)
	if err != nil {
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "DeleteConfig",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  err.Error(),
		}
		t.errorChan <- errMsg
	}

	return err
}

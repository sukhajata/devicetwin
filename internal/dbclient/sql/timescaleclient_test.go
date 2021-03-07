package sql

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sukhajata/devicetwin/internal/dbclient/nosql"
	"github.com/sukhajata/devicetwin/internal/types"
	"github.com/sukhajata/devicetwin/internal/utility"
	"github.com/sukhajata/devicetwin/mocks"
	pb "github.com/sukhajata/ppconfig"
	pbLogger "github.com/sukhajata/pplogger"
	"testing"
)

func setupTimescaleTest(mockCtrl *gomock.Controller) (*TimescaleClient, *mocks.MockSQLEngine) {
	mockSQLEngine := mocks.NewMockSQLEngine(mockCtrl)
	errorChan := make(chan *pbLogger.ErrorMessage, 2)
	client := NewTimescaleClient(mockSQLEngine, errorChan)
	return client, mockSQLEngine
}

func TestTimescaleClient_GetFieldDetailsByIndex(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row := []interface{}{int32(3), "roffset", "i", int32(0), "The radio offset", int32(0), int32(3000)}
	results := []interface{}{row}

	ppdev := "meter"
	index := int32(3)
	firmwareVersion := "1.2.0"
	//id := "123"

	queryString := `SELECT "INDEX", "NAME", "TYPE", "DEFAULT", "DESCRIPTION", "MIN", "MAX"
		FROM "CONFIG_SCHEMA"
		WHERE "PPDEV" = $1
		AND "INDEX" = $2
		AND "PPVER" = $3`

	mockDBEngine.EXPECT().Query(queryString, ppdev, index, firmwareVersion).Return(results, nil).Times(1)

	/*row = []interface{}{int32(3), "roffset", "i", int32(0), "The radio offset", int32(0), int32(3000)}
	results := []interface{}{row}

	queryString = `SELECT "DESIRED", "REPORTED" FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	mockDBEngine.EXPECT().Query(queryString, id, 0, fieldDetails.Name).Return()*/

	field, err := client.GetFieldDetailsByIndex(index, firmwareVersion, nosql.DocTypeConfigSchema)
	require.Nil(t, err)
	require.Equal(t, "roffset", field.Name)
	require.Equal(t, int32(3), field.Index)
}

func TestTimescaleClient_GetFieldDetailsByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row := []interface{}{int32(3), "roffset", "i", int32(0), "The radio offset", int32(0), int32(3000)}

	results := []interface{}{row}

	queryString := `SELECT "INDEX", "NAME", "TYPE", "DEFAULT", "DESCRIPTION", "MIN", "MAX"
		FROM "CONFIG_SCHEMA"
		WHERE "PPDEV" = $1
		AND "NAME" = $2
		AND "PPVER" = $3`
	ppdev := "meter"
	fieldName := "roffset"
	firmwareVersion := "1.2.0"

	mockDBEngine.EXPECT().Query(queryString, ppdev, fieldName, firmwareVersion).Return(results, nil).Times(1)

	field, err := client.GetFieldDetailsByName(fieldName, firmwareVersion, nosql.DocTypeConfigSchema)
	require.Nil(t, err)
	require.Equal(t, "roffset", field.Name)
	require.Equal(t, int32(3), field.Index)
}

func TestTimescaleClient_UpdateDbDesired(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	details := types.ConfigFieldDetails{
		Index: 3.0,
		Name:  "roffset",
		Type:  "i",
	}
	req := pb.SetDesiredRequest{
		Identifier: "123",
		FieldName:  "roffset",
		FieldValue: "200",
		Slot:       0,
	}
	value, err := utility.StringToInterface(details, req.GetFieldValue())

	queryString := `SELECT 1 FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	mockDBEngine.EXPECT().Query(queryString, req.Identifier, req.Slot, details.Name).Return(make([]interface{}, 1), nil).Times(1)

	queryString = `UPDATE "CONFIG" SET "DESIRED" = $1 WHERE "CONNECTIONID" = $2 AND "SLOT" = $3 AND "NAME" = $4`
	mockDBEngine.EXPECT().Exec(queryString, fmt.Sprintf("%v", value), req.Identifier, req.Slot, details.Name).Return(nil).Times(1)

	err = client.UpdateDbDesired(&req, details)
	require.Nil(t, err)
}

/*
func TestTimescaleClient_GetS11ConfigKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	id := "123"
	slot := int32(100)
	row := []interface{}{"456"}
	results := []interface{}{row}
	slotName := fmt.Sprintf("s%v", slot)
	queryString := fmt.Sprintf(`SELECT "DATA" -> 'slots' ->> %s FROM "CONNECTIONS_JSON" WHERE "ID" = $1`, slotName)

	mockDBEngine.EXPECT().Query(queryString, id).Return(results, nil).Times(1)

	key, err := client.GetS11ConfigKey(id, slot)
	require.Nil(t, err)
	require.Equal(t, "456", key)
}*/

func TestTimescaleClient_UpdateDbReported(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	details := types.ConfigFieldDetails{
		Index: 3.0,
		Name:  "roffset",
		Type:  "i",
	}
	req := pb.UpdateReportedRequest{
		DeviceEUI:  "123",
		FieldIndex: 3,
		FieldValue: []byte{0x00, 0x00, 0x12, 0x12},
		Slot:       0,
	}
	key := req.DeviceEUI

	value, err := utility.DecodeFieldValue(details, req.GetFieldValue())
	require.NoError(t, err)

	queryString := `SELECT 1 FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	mockDBEngine.EXPECT().Query(queryString, key, req.Slot, details.Name).Return(make([]interface{}, 1), nil).Times(1)

	queryString = `UPDATE "CONFIG" SET "REPORTED" = $1 WHERE "CONNECTIONID" = $2 AND "SLOT" = $3 AND "NAME" = $4`
	mockDBEngine.EXPECT().Exec(queryString, fmt.Sprintf("%v", value), key, req.Slot, details.Name).Return(nil).Times(1)

	err = client.UpdateDbReported(&req, details)
	require.Nil(t, err)
}

func TestTimescaleClient_GetNextRadioOffset(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row := []interface{}{1000}
	results := []interface{}{row}

	queryString := `SELECT MAX(CAST("DESIRED" AS int)) AS maxroffset FROM "CONFIG" WHERE "NAME" = 'roffset' and "DESIRED" != ''`
	mockDBEngine.EXPECT().Query(queryString).Return(results, nil).Times(1)

	roffset, err := client.GetNextRadioOffset()
	require.Nil(t, err)
	require.Equal(t, 1010, roffset)
}

func TestTimescaleClient_GetLatestFirmware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row := []interface{}{"1.1.4"}
	results := []interface{}{row}
	ppdev := "meter"

	queryString := `SELECT "PPVER" FROM "CONFIG_SCHEMA" WHERE "PPDEV" = $1 ORDER BY "PPORDER" DESC LIMIT 1`
	mockDBEngine.EXPECT().Query(queryString, ppdev).Return(results, nil).Times(1)

	firmware, err := client.GetLatestFirmware("configschema")
	require.Nil(t, err)
	require.Equal(t, "1.1.4", firmware)
}

func TestTimescaleClient_GetFieldDetails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row1 := []interface{}{"roffset", int32(3), "i", int32(0), "The radio offset", int32(0), int32(3000)}
	row2 := []interface{}{"dlresmin", int32(4), 10, "6,8", "Downlink reserved minutes", "5,7", "6,8"}
	//"NAME", "INDEX", "TYPE", "DEFAULT", "DESCRIPTION", "MIN", "MAX"
	results := []interface{}{row1, row2}

	ppdev := "meter"
	firmware := "1.2.0"
	queryString := `SELECT "NAME", "INDEX", "TYPE", "DEFAULT", "DESCRIPTION", "MIN", "MAX"
		FROM "CONFIG_SCHEMA"
		WHERE "PPVER" = $1
		AND "PPDEV" = $2`
	mockDBEngine.EXPECT().Query(queryString, firmware, ppdev).Return(results, nil).Times(1)

	fieldDetails, err := client.GetFieldDetails(firmware, nosql.DocTypeConfigSchema)
	require.Nil(t, err)
	require.Equal(t, int32(3), fieldDetails["roffset"].Index)
	require.Equal(t, int32(4), fieldDetails["dlresmin"].Index)
	require.Equal(t, "dlresmin", fieldDetails["dlresmin"].Name)
}

func TestTimescaleClient_UpdateConfigToNewFirmware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	rows := make(map[string]types.ConfigFieldDetails)
	rows["roffset"] = types.ConfigFieldDetails{
		Index: 3.0,
		Name:  "roffset",
		Type:  "i",
	}
	id := "123"

	for k := range rows {
		queryString := `SELECT 1 FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "NAME" = $2 AND "SLOT" = $3`

		//return empty slice
		mockDBEngine.EXPECT().Query(queryString, id, k, 0).Return([]interface{}{}, nil).Times(1)

		queryString = `INSERT INTO "CONFIG" ("CONNECTIONID", "NAME", "SLOT", "DESIRED", "REPORTED") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
		mockDBEngine.EXPECT().Exec(queryString, id, k, 0, "", "").Return(nil).Times(1)
	}

	client.UpdateConfigToNewFirmware("123", 0, rows)
}

func TestTimescaleClient_GetConfigByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row := []interface{}{2000, 1500}
	results := []interface{}{row}

	details := types.ConfigFieldDetails{
		Index: 3.0,
		Name:  "roffset",
		Type:  "i",
	}
	req := pb.GetConfigByNameRequest{
		Identifier: "123",
		FieldName:  "roffset",
		Slot:       0,
	}

	queryString := `SELECT "DESIRED", "REPORTED" FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = $2 AND "NAME" = $3`
	mockDBEngine.EXPECT().Query(queryString, req.Identifier, req.Slot, req.FieldName).Return(results, nil).Times(1)

	field, err := client.GetConfigByName("1.1.4", details, &req)
	require.Nil(t, err)
	require.Equal(t, "roffset", field.Name)
	require.Equal(t, int32(3), field.Index)
	require.Equal(t, "2000", field.Desired)
	require.Equal(t, "1500", field.Reported)
}

func TestTimescaleClient_GetDLResmin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row := []interface{}{"6,8"}
	results := []interface{}{row}
	id := "123"

	queryString := `SELECT "REPORTED" FROM "CONFIG" WHERE "CONNECTIONID" = $1 AND "SLOT" = 0 AND "NAME" = 'dlresmin'`
	mockDBEngine.EXPECT().Query(queryString, id).Return(results, nil).Times(1)

	dlresmin, err := client.GetDLResmin("123")
	require.Nil(t, err)
	require.Equal(t, "6,8", dlresmin)
}

func TestTimescaleClient_GetInconsistentDevices(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setupTimescaleTest(mockCtrl)

	row1 := []interface{}{"123"}
	row2 := []interface{}{"456"}

	results := []interface{}{row1, row2}

	queryString := `SELECT DISTINCT "CONNECTIONID" FROM "CONFIG" WHERE "DESIRED" != "REPORTED"`
	mockDBEngine.EXPECT().Query(queryString).Return(results, nil).Times(1)

	inconsistent, err := client.GetInconsistentDevices()
	require.Nil(t, err)
	require.Equal(t, 2, len(inconsistent))
	require.Equal(t, "123", inconsistent[0])
}

package nosql

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/sukhajata/devicetwin.git/internal/types"
	"github.com/sukhajata/devicetwin.git/internal/utility"
	pb "github.com/sukhajata/ppconfig"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sukhajata/devicetwin.git/mocks"
)

var (
	bucketNameShared = "shared"
	bucketName       = "test"
)

func setup(mockCtrl *gomock.Controller) (*CouchbaseClient, *mocks.MockNoSQLEngine) {
	mockNoSQLEngine := mocks.NewMockNoSQLEngine(mockCtrl)
	mockHelper := mocks.NewMockHelper(mockCtrl)
	client := NewCouchbaseClient(mockNoSQLEngine, bucketName, bucketNameShared, mockHelper)
	return client, mockNoSQLEngine
}

func TestCouchbaseClient_GetFieldDetailsByIndex(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row := map[string]interface{}{
		"i": float64(3),
		"n": "roffset",
		"t": "i",
	}
	results := []interface{}{row}
	queryString := fmt.Sprintf("SELECT f.i, f.n, f.t, f.d, f.a, f.b, f.c FROM %s "+
		" s UNNEST ppschema f WHERE s.type=$1"+
		" AND f.i = $2"+
		" AND s.ppver = $3", bucketNameShared)
	index := int32(3)
	firmware := "1.2.0"

	mockDBEngine.EXPECT().Query(bucketNameShared, queryString, []interface{}{DocTypeConfigSchema, index, firmware}).Return(results, nil).Times(1)

	field, err := client.GetFieldDetailsByIndex(int32(3), firmware, DocTypeConfigSchema)
	require.Nil(t, err)
	require.Equal(t, "roffset", field.Name)
	require.Equal(t, int32(3), field.Index)
}

func TestCouchbaseClient_GetFieldDetailsByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row := map[string]interface{}{
		"i": 3.0,
		"n": "roffset",
		"t": "i",
	}

	results := []interface{}{row}
	firmware := "1.2.0"
	queryString := fmt.Sprintf("SELECT f.i, f.n, f.t, f.d, f.a, f.b, f.c FROM %s "+
		" s UNNEST ppschema f WHERE s.type=$1"+
		" AND f.n = $2"+
		" AND s.ppver = $3", bucketNameShared)
	name := "roffset"

	mockDBEngine.EXPECT().Query(bucketNameShared, queryString, []interface{}{DocTypeConfigSchema, name, firmware}).Return(results, nil).Times(1)

	field, err := client.GetFieldDetailsByName("roffset", firmware, DocTypeConfigSchema)
	require.Nil(t, err)
	require.Equal(t, "roffset", field.Name)
	require.Equal(t, int32(3), field.Index)
}

func TestCouchbaseClient_UpdateDbDesired(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

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
	fieldPath := "config.desired.roffset"
	value, err := utility.StringToInterface(details, req.GetFieldValue())
	require.NoError(t, err)

	mockDBEngine.EXPECT().Update(bucketName, req.GetIdentifier(), fieldPath, value).Return(nil).Times(1)

	err = client.UpdateDbDesired(&req, details)
	require.Nil(t, err)
}

func TestCouchbaseClient_GetS11ConfigKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row := make(map[string]interface{})
	row["s100"] = "456"
	results := []interface{}{row}

	id := "123"
	slotName := "s100"
	queryString := fmt.Sprintf("SELECT connection.slots.%s FROM %s WHERE meta().id = $1", slotName, bucketName)

	mockDBEngine.EXPECT().Query(bucketName, queryString, []interface{}{id}).Return(results, nil).Times(1)

	key, err := client.GetS11ConfigKey(id, 100)
	require.Nil(t, err)
	require.Equal(t, "456", key)
}

func TestCouchbaseClient_UpdateDbReported(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

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

	fieldPath := "config.reported.roffset"
	value, err := utility.DecodeFieldValue(details, req.GetFieldValue())
	require.NoError(t, err)

	mockDBEngine.EXPECT().Update(bucketName, req.DeviceEUI, fieldPath, value).Return(nil).Times(1)

	err = client.UpdateDbReported(&req, details)
	require.NoError(t, err)
}

func TestCouchbaseClient_GetNextRadioOffset(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row := make(map[string]interface{})
	row["maxroffset"] = 1000
	results := []interface{}{row}
	queryString := fmt.Sprintf("SELECT MAX(config.desired.roffset) AS maxroffset FROM %s WHERE type=$1 and config.desired.roffset != ''", bucketName)

	mockDBEngine.EXPECT().Query(bucketName, queryString, []interface{}{docTypeConnection}).Return(results, nil).Times(1)

	roffset, err := client.GetNextRadioOffset()
	require.Nil(t, err)
	require.Equal(t, 1010, roffset)
}

func Test_GetLatestFirmware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row := make(map[string]interface{})
	row["ppver"] = "1.1.4"

	results := []interface{}{row}
	queryString := fmt.Sprintf("SELECT ppver FROM %s WHERE type = $1 ORDER BY pporder DESC LIMIT 1", bucketNameShared)

	mockDBEngine.EXPECT().Query(bucketNameShared, queryString, []interface{}{DocTypeConfigSchema}).Return(results, nil).Times(1)

	firmware, err := client.GetLatestFirmware("configschema")
	require.Nil(t, err)
	require.Equal(t, "1.1.4", firmware)
}

func Test_GetFieldDetails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row1 := map[string]interface{}{
		"i": 3.0,
		"n": "roffset",
		"t": "i",
	}

	row2 := map[string]interface{}{
		"i": 4.0,
		"n": "dlresmin",
		"t": 10,
	}

	results := []interface{}{row1, row2}
	firmware := "1.2.0"

	queryString := fmt.Sprintf("SELECT f.i, f.n, f.t, f.d, f.a, f.b, f.c FROM %s s UNNEST s.ppschema AS f WHERE s.type=$1 "+
		"AND s.ppver=$2 ORDER BY f.n", bucketNameShared)
	mockDBEngine.EXPECT().Query(bucketNameShared, queryString, []interface{}{DocTypeConfigSchema, firmware}).Return(results, nil).Times(1)

	fieldDetails, err := client.GetFieldDetails(firmware, DocTypeConfigSchema)
	require.Nil(t, err)
	require.Equal(t, int32(3), fieldDetails["roffset"].Index)
	require.Equal(t, int32(4), fieldDetails["dlresmin"].Index)
	require.Equal(t, "dlresmin", fieldDetails["dlresmin"].Name)
}

func Test_UpdateConfigToNewFirmware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	desired := make(map[string]interface{})
	desired["roffset"] = 2000
	desired["dlresmin"] = "6,8"

	reported := make(map[string]interface{})
	reported["roffset"] = 5
	reported["dlresmin"] = "6,7"

	configDoc := &types.UntypedConfigDoc{
		Desired:  desired,
		Reported: reported,
	}
	rows := make(map[string]types.ConfigFieldDetails)
	rows["roffset"] = types.ConfigFieldDetails{
		Index: 3.0,
		Name:  "roffset",
		Type:  "i",
	}
	rows["dlresmin"] = types.ConfigFieldDetails{
		Index: 4.0,
		Name:  "dlresmin",
		Type:  10,
	}
	id := "123"
	slot := 0

	//get the old doc
	mockDBEngine.EXPECT().Lookup(bucketName, id, "config", gomock.Any()).Return(nil).SetArg(3, configDoc).Times(1)

	//insert with new config
	mockDBEngine.EXPECT().Update(bucketName, id, "config", gomock.Any()).Return(nil).Times(1)

	client.UpdateConfigToNewFirmware(id, slot, rows)
}

func Test_GetConfigByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row := make(map[string]interface{})
	row["desired"] = 2000
	row["reported"] = 1500

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

	queryString := fmt.Sprintf("SELECT config.desired.%s AS desired, config.reported.%s AS reported FROM %s c WHERE meta(c).id = $1", req.GetFieldName(), req.GetFieldName(), bucketName)
	mockDBEngine.EXPECT().Query(bucketName, queryString, []interface{}{req.Identifier}).Return(results, nil).Times(1)

	field, err := client.GetConfigByName("1.1.4", details, &req)
	require.Nil(t, err)
	require.Equal(t, "roffset", field.Name)
	require.Equal(t, int32(3), field.Index)
	require.Equal(t, "2000", field.Desired)
	require.Equal(t, "1500", field.Reported)
}

func Test_GetDLResmin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row := make(map[string]interface{})
	row["dlresmin"] = "6,8"

	results := []interface{}{row}
	id := "123"
	queryString := fmt.Sprintf("SELECT config.reported.dlresmin AS dlresmin FROM %s c WHERE meta(c).id = $1", bucketName)

	mockDBEngine.EXPECT().Query(bucketName, queryString, []interface{}{id}).Return(results, nil).Times(1)

	dlresmin, err := client.GetDLResmin(id)
	require.Nil(t, err)
	require.Equal(t, "6,8", dlresmin)
}

func Test_GetInconsistentDevices(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client, mockDBEngine := setup(mockCtrl)

	row1 := make(map[string]interface{})
	row1["id"] = "123"
	row2 := make(map[string]interface{})
	row2["id"] = "456"

	results := []interface{}{row1, row2}
	queryString := fmt.Sprintf("SELECT meta(b).id FROM %s b WHERE b.type = $1 AND b.config.desired != b.config.reported", bucketName)

	mockDBEngine.EXPECT().Query(bucketName, queryString, []interface{}{docTypeConnection}).Return(results, nil).Times(1)

	inconsistent, err := client.GetInconsistentDevices()
	require.Nil(t, err)
	fmt.Println(inconsistent)
	require.Equal(t, 2, len(inconsistent))
	require.Equal(t, "123", inconsistent[0])
}

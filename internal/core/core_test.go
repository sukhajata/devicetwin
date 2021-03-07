package core

import (
	"github.com/sukhajata/devicetwin/internal/dbclient/nosql"
	"github.com/sukhajata/devicetwin/internal/types"
	pbLogger "github.com/sukhajata/pplogger"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sukhajata/devicetwin/mocks"
	pb "github.com/sukhajata/ppconfig"
	"github.com/sukhajata/ppmessage/ppdownlink"
	"github.com/sukhajata/ppmessage/ppuplink"
)

func setup(mockCtrl *gomock.Controller) (*Service, *mocks.MockClient, *mocks.MockConnectionServiceClient) {
	mockHelper := mocks.NewMockHelper(mockCtrl)
	mockDBClient := mocks.NewMockClient(mockCtrl)
	mockConnectionClient := mocks.NewMockConnectionServiceClient(mockCtrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(mockCtrl)
	transmitChan := make(chan *ppdownlink.ConfigDownlinkMessage, 2)
	deviceEventChan := make(chan *pbLogger.DeviceLogMessage, 2)
	errorChan := make(chan *pbLogger.ErrorMessage, 2)
	consistencyService := mocks.NewMockConsistencyChecker(mockCtrl)

	service := NewService(
		mockDBClient,
		mockConnectionClient,
		mockAuthClient,
		consistencyService,
		"servicekey",
		mockHelper,
		transmitChan,
		errorChan,
		deviceEventChan,
		"powerpilot-admin",
		"poewrpilot-installer",
		"powerpilot-superuser",
	)

	return service, mockDBClient, mockConnectionClient
}

func Test_HandleConfigUplink(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	service, mockDBClient, _ := setup(mockCtrl)

	/*row := map[string]interface{}{
		"i": 3.0,
		"n": "roffset",
		"t": "i",
	}*/

	//results := []interface{}{row}
	firmware := "1.2.0"
	req := &pb.UpdateReportedRequest{
		DeviceEUI:  "ABC",
		FieldIndex: int32(3),
		FieldValue: []byte{0x00, 0x00, 0x12, 0x12},
		Slot:       int32(0),
	}
	details := types.ConfigFieldDetails{
		Index: 3.0,
		Name:  "roffset",
		Type:  "i",
	}

	mockDBClient.EXPECT().GetLatestFirmware(nosql.DocTypeConfigSchema).Return(firmware, nil).Times(1)

	mockDBClient.EXPECT().GetFieldDetailsByIndex(req.FieldIndex, firmware, nosql.DocTypeConfigSchema).Return(details, nil).Times(1)

	mockDBClient.EXPECT().UpdateDbReported(req, details).Return(nil).Times(1)

	//mockConnectionClient.EXPECT().GetConnection(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

	//mockConnectionClient.EXPECT().UpdateConnection(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

	uplink := &ppuplink.ConfigUplinkMessage{
		Deviceeui: "ABC",
		Index:     3,
		Value:     []byte{0x00, 0x00, 0x12, 0x12},
	}
	service.HandleConfigUplink(uplink)

}

/*
func Test_SetDesired(t *testing.T) {
	row := map[string]interface{}{
		"i": 3.0,
		"n": "roffset",
		"t": "i",
	}

	results := []interface{}{row}
	mockDBEngine := mocks.MockNoSQLEngine{
		Results: results,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	service := setup(mockCtrl, mockDBEngine)

	req := &pb.SetDesiredRequest{
		Identifier: "123",
		FieldName:  "roffset",
		FieldValue: "2000",
	}
	response, err := service.SetDesired("test", req)
	require.NoError(t, err)
	fmt.Println(response)

}

*/

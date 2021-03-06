package messageprocessor

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sukhajata/devicetwin.git/internal/dbclient/nosql"
	"github.com/sukhajata/devicetwin.git/internal/types"
	"github.com/sukhajata/devicetwin.git/mocks"
	"github.com/sukhajata/devicetwin.git/pkg/ppmqtt"
	pbLogger "github.com/sukhajata/pplogger"
	"github.com/sukhajata/ppmessage/ppuplink"
	"google.golang.org/protobuf/proto"
	"testing"
)

func setup(mockCtrl *gomock.Controller) (*MessageProcessor, *mocks.MockClient, *mocks.MockConfigHandler, <-chan *pbLogger.ErrorMessage) {
	coreService := mocks.NewMockConfigHandler(mockCtrl)
	errorChan := make(chan *pbLogger.ErrorMessage, 2)
	dbClient := mocks.NewMockClient(mockCtrl)

	return &MessageProcessor{
		coreService: coreService,
		dbClient:    dbClient,
		errorChan:   errorChan,
	}, dbClient, coreService, errorChan
}

func Test_ProcessUplinkMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	processor, _, coreService, errorChan := setup(mockCtrl)

	// fail on any error
	go func(errorChan <-chan *pbLogger.ErrorMessage) {
		for msg := range errorChan {
			t.Fatal(msg.Message)
		}
	}(errorChan)

	// setup
	uplink := &ppuplink.ConfigUplinkMessage{
		Deviceeui: "123",
		Slot:      0,
		Firmware:  "1.2.0",
		Value:     []byte{0x01, 0x02},
	}
	msgBytes, err := proto.Marshal(uplink)
	require.NoError(t, err)

	msg := ppmqtt.Message{
		Topic:   "application/powerpilot/uplink/config/123",
		Payload: msgBytes,
	}

	// expect
	coreService.EXPECT().HandleConfigUplink(uplink)

	// call
	processor.ProcessMessage(msg)
}

func Test_ProcessConnectionsMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	processor, dbClient, _, errorChan := setup(mockCtrl)

	// fail on any error
	go func(errorChan <-chan *pbLogger.ErrorMessage) {
		for msg := range errorChan {
			t.Fatal(msg.Message)
		}
	}(errorChan)

	// setup
	update := &pbLogger.DeviceLogMessage{
		User:      "bob",
		DeviceEUI: "17e6c639-670e-4659-84ce-ac781fd09597",
		Message:   "Created pending connection: 17e6c639-670e-4659-84ce-ac781fd09597",
	}
	msgBytes, err := proto.Marshal(update)
	require.NoError(t, err)
	msg := ppmqtt.Message{
		Topic:   "application/powerpilot/connections",
		Payload: msgBytes,
	}
	firmware := "1.2.0"
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

	// expect
	dbClient.EXPECT().GetLatestFirmware(nosql.DocTypeConfigSchema).Return(firmware, nil)
	dbClient.EXPECT().GetFieldDetails(firmware, nosql.DocTypeConfigSchema).Return(rows, nil)
	dbClient.EXPECT().UpdateConfigToNewFirmware(update.DeviceEUI, 0, rows)

	// call
	processor.ProcessMessage(msg)
}

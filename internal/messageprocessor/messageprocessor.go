package messageprocessor

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sukhajata/devicetwin.git/internal/core"
	"github.com/sukhajata/devicetwin.git/internal/dbclient"
	"github.com/sukhajata/devicetwin.git/internal/dbclient/nosql"
	"github.com/sukhajata/devicetwin.git/pkg/loggerhelper"
	"github.com/sukhajata/devicetwin.git/pkg/ppmqtt"
	pbLogger "github.com/sukhajata/pplogger"
	"github.com/sukhajata/ppmessage/ppuplink"
	"strconv"
	"strings"
)

type MessageProcessor struct {
	coreService core.ConfigHandler
	dbClient    dbclient.Client
	errorChan   chan *pbLogger.ErrorMessage
}

func NewMessageProcessor(coreService core.ConfigHandler, dbClient dbclient.Client, errorChan chan *pbLogger.ErrorMessage) *MessageProcessor {
	return &MessageProcessor{
		coreService: coreService,
		dbClient:    dbClient,
		errorChan:   errorChan,
	}
}

func (p *MessageProcessor) ProcessMessage(msg ppmqtt.Message) {
	loggerhelper.WriteToLog(fmt.Sprintf("TOPIC: %s\n", msg.Topic))
	if strings.Contains(msg.Topic, "uplink/config") {
		var configMessage ppuplink.ConfigUplinkMessage
		err := proto.Unmarshal(msg.Payload, &configMessage)
		if err != nil {
			errMsg := &pbLogger.ErrorMessage{
				Service:  "config-service",
				Function: "ProcessMessage",
				Message:  err.Error(),
				Severity: pbLogger.ErrorMessage_FATAL,
			}
			p.errorChan <- errMsg
			return
		}

		go p.coreService.HandleConfigUplink(&configMessage)

	} else if strings.Contains(msg.Topic, "connections") {
		var logMessage pbLogger.DeviceLogMessage
		err := proto.Unmarshal(msg.Payload, &logMessage)
		if err != nil {
			errMsg := &pbLogger.ErrorMessage{
				Service:  "config-service",
				Function: "ProcessMessage",
				Message:  err.Error(),
				Severity: pbLogger.ErrorMessage_FATAL,
			}
			p.errorChan <- errMsg
			return
		}
		loggerhelper.WriteToLog(logMessage.Message)

		if strings.Contains(logMessage.Message, "Created pending") || strings.Contains(logMessage.Message, "Created connection") {
			// create blank config fields
			firmware, err := p.dbClient.GetLatestFirmware(nosql.DocTypeConfigSchema)
			if err != nil {
				return
			}
			fieldDetails, err := p.dbClient.GetFieldDetails(firmware, nosql.DocTypeConfigSchema)
			if err != nil {
				return
			}

			p.dbClient.UpdateConfigToNewFirmware(logMessage.DeviceEUI, 0, fieldDetails)

		} else if strings.Contains(logMessage.Message, "Deleted connection") {
			err := p.dbClient.DeleteConfig(logMessage.DeviceEUI, 0)
			if err != nil {
				loggerhelper.WriteToLog(err)
			}
		} else if strings.Contains(logMessage.Message, "Added slot") {
			// get slot
			ww := strings.Split(logMessage.Message, " ") // Added slot 100
			ss := ww[2]
			slot, err := strconv.Atoi(ss)
			if err != nil {
				errMsg := &pbLogger.ErrorMessage{
					Service:  "config-service",
					Function: "ProcessMessage",
					Message:  err.Error(),
					Severity: pbLogger.ErrorMessage_FATAL,
				}
				p.errorChan <- errMsg
				return
			}
			// create s11 config fields
			firmware, err := p.dbClient.GetLatestFirmware(nosql.DocTypeS11ConfigSchema)
			if err != nil {
				errMsg := &pbLogger.ErrorMessage{
					Service:  "config-service",
					Function: "ProcessMessage",
					Message:  err.Error(),
					Severity: pbLogger.ErrorMessage_FATAL,
				}
				p.errorChan <- errMsg
				return
			}
			fieldDetails, err := p.dbClient.GetFieldDetails(firmware, nosql.DocTypeS11ConfigSchema)
			if err != nil {
				errMsg := &pbLogger.ErrorMessage{
					Service:  "config-service",
					Function: "ProcessMessage",
					Message:  err.Error(),
					Severity: pbLogger.ErrorMessage_FATAL,
				}
				p.errorChan <- errMsg
				return
			}
			p.dbClient.UpdateConfigToNewFirmware(logMessage.DeviceEUI, slot, fieldDetails)

		}
	}
}

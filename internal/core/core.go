package core

import (
	"errors"
	"fmt"
	"github.com/sukhajata/devicetwin/internal/consistency"
	"github.com/sukhajata/devicetwin/internal/dbclient"
	"github.com/sukhajata/devicetwin/internal/dbclient/nosql"
	"github.com/sukhajata/devicetwin/internal/utility"
	"github.com/sukhajata/devicetwin/pkg/authhelper"
	"github.com/sukhajata/devicetwin/pkg/loggerhelper"
	pbAuth "github.com/sukhajata/ppauth"
	pb "github.com/sukhajata/ppconfig"
	pbConnection "github.com/sukhajata/ppconnection"
	pbLogger "github.com/sukhajata/pplogger"
	"github.com/sukhajata/ppmessage/ppdownlink"
	"github.com/sukhajata/ppmessage/ppuplink"
)

type ConfigHandler interface {
	HandleConfigUplink(msg *ppuplink.ConfigUplinkMessage)
	AssignRadioOffset(token string, identifier *pb.Identifier) (*pb.Response, error)
	SetDesired(token string, req *pb.SetDesiredRequest) (*pb.Response, error)
	SendConsistencyCheckRequest(downlink *ppdownlink.ConfigDownlinkMessage)
	UpdateReported(req *pb.UpdateReportedRequest) (*pb.Response, error)
	GetConfigByName(token string, req *pb.GetConfigByNameRequest) (*pb.ConfigField, error)
	GetNewConfigDoc(token string, req *pb.Identifier) (*pb.ConfigDoc, error)
	GetConfigByIndex(token string, req *pb.GetConfigByIndexRequest) (*pb.ConfigField, error)
	GetDeviceConfig(token string, req *pb.Identifier) (*pb.ConfigFields, error)
	UpdateFirmwareAllDevices(token string) error
}

// Service provides core services
type Service struct {
	dbClient             dbclient.Client
	grpcConnectionClient pbConnection.ConnectionServiceClient
	grpcAuthClient       pbAuth.AuthServiceClient
	consistencyService   consistency.ConsistencyChecker
	serviceKey           string
	transmitChan         chan<- *ppdownlink.ConfigDownlinkMessage
	loggerHelper         loggerhelper.Helper
	errorChan            chan<- *pbLogger.ErrorMessage
	deviceEventChan      chan<- *pbLogger.DeviceLogMessage
	adminRole            string
	installerRole        string
	superuserRole        string
}

// NewService factory method
func NewService(
	dbClient dbclient.Client,
	grpcConnectionClient pbConnection.ConnectionServiceClient,
	grpcAuthClient pbAuth.AuthServiceClient,
	consistencyService consistency.ConsistencyChecker,
	serviceKey string,
	loggerHelper loggerhelper.Helper,
	transmitChan chan<- *ppdownlink.ConfigDownlinkMessage,
	errorChan chan<- *pbLogger.ErrorMessage,
	deviceEventChan chan<- *pbLogger.DeviceLogMessage,
	adminRole string,
	installerRole string,
	superuserRole string,
) *Service {

	cs := &Service{
		dbClient:             dbClient,
		grpcConnectionClient: grpcConnectionClient,
		grpcAuthClient:       grpcAuthClient,
		consistencyService:   consistencyService,
		serviceKey:           serviceKey,
		loggerHelper:         loggerHelper,
		transmitChan:         transmitChan,
		deviceEventChan:      deviceEventChan,
		errorChan:            errorChan,
		adminRole:            adminRole,
		installerRole:        installerRole,
		superuserRole:        superuserRole,
	}

	return cs
}

// HandleConfigUplink handle reported config messages
func (c *Service) HandleConfigUplink(configMessage *ppuplink.ConfigUplinkMessage) {
	loggerhelper.WriteToLog(fmt.Sprintf("Received uplink index %v value %v", configMessage.Index, configMessage.Value))

	updateRequest := &pb.UpdateReportedRequest{
		DeviceEUI:  configMessage.Deviceeui,
		FieldIndex: int32(configMessage.Index),
		FieldValue: configMessage.Value,
		Slot:       int32(configMessage.Slot),
	}

	_, err := c.UpdateReported(updateRequest)
	loggerhelper.WriteToLog(err)
}

// AssignRadioOffset assign incremental value
func (c *Service) AssignRadioOffset(token string, identifier *pb.Identifier) (*pb.Response, error) {
	allowedRoles := []string{c.adminRole, c.installerRole, c.superuserRole}
	_, err := authhelper.CheckToken(c.grpcAuthClient, token, allowedRoles)
	if err != nil {
		return &pb.Response{
			Reply: "NOT AUTHORIZED",
		}, err
	}

	roffset, err := c.dbClient.GetNextRadioOffset()
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	request := &pb.SetDesiredRequest{
		Identifier: identifier.Identifier,
		FieldName:  "roffset",
		FieldValue: fmt.Sprintf("%v", roffset),
	}

	return c.SetDesired(token, request)
}

// SetDesired - publish to mqtt, log the change
func (c *Service) SetDesired(token string, req *pb.SetDesiredRequest) (*pb.Response, error) {
	allowedRoles := []string{c.adminRole, c.installerRole, c.superuserRole}
	username, err := authhelper.CheckToken(c.grpcAuthClient, token, allowedRoles)
	if err != nil {
		return &pb.Response{
			Reply: "NOT AUTHORIZED",
		}, err
	}

	loggerhelper.WriteToLog(fmt.Sprintf("Setting config %v value %v", req.FieldName, req.FieldValue))
	if req.GetIdentifier() == "" {
		return &pb.Response{
			Reply: "NOT OK",
		}, errors.New("missing identifier")
	}

	if req.GetFieldName() == "" {
		return &pb.Response{
			Reply: "NOT OK",
		}, errors.New("missing field name")
	}

	docType := nosql.DocTypeConfigSchema

	if req.Slot > 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}

	firmware, err := c.dbClient.GetLatestFirmware(docType)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}
	fieldDetails, err := c.dbClient.GetFieldDetailsByName(req.FieldName, firmware, docType)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	// build downlink message, validate value
	downlink, err := utility.BuildDownlinkMessage(req.Identifier, fieldDetails, req.FieldValue, firmware, 0, uint32(req.Slot))
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	// log change
	// get old value
	configField, err := c.dbClient.GetConfigByName(firmware, fieldDetails, &pb.GetConfigByNameRequest{
		Identifier: req.Identifier,
		FieldName:  req.FieldName,
		Slot:       req.Slot,
	})
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	logMessage := &pbLogger.DeviceLogMessage{
		User:      username,
		DeviceEUI: req.Identifier,
		Message:   fmt.Sprintf("Changed %s from %s to %s slot %v", req.FieldName, configField.Desired, req.FieldValue, req.Slot),
	}
	c.deviceEventChan <- logMessage

	// update dbclient
	err = c.dbClient.UpdateDbDesired(req, fieldDetails)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}
	loggerhelper.WriteToLog("Updated dbclient")

	// check if this is installed before sending downlink
	connectionRequest := pbConnection.Identifier{
		Identifier: req.Identifier,
	}
	ctx, cancel := authhelper.GetContextWithAuth(c.serviceKey)
	defer cancel()
	conn, err := c.grpcConnectionClient.GetConnection(ctx, &connectionRequest)
	if err != nil {
		errMsg := &pbLogger.ErrorMessage{
			Service:  "config-service",
			Function: "SetDesired",
			Severity: pbLogger.ErrorMessage_FATAL,
			Message:  fmt.Sprintf("error getting connection - %v", err.Error()),
		}
		c.errorChan <- errMsg

		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}
	if conn.Device != nil && conn.Device.DeviceEUI != "" {
		loggerhelper.WriteToLog(fmt.Sprintf("Sending command: %v", conn.Device.DeviceEUI))

		// send
		c.transmitChan <- downlink

		// schedule consistency check
		go c.SendConsistencyCheckRequest(downlink)
	} else {
		loggerhelper.WriteToLog("Not sending command")
	}

	// send response
	return &pb.Response{
		Reply: "OK",
	}, nil
}

// SendConsistencyCheckRequest - schedule consistency check
func (c *Service) SendConsistencyCheckRequest(downlink *ppdownlink.ConfigDownlinkMessage) {
	checkConsistencyRequest := &pb.CheckConsistencyRequest{
		DeviceEUI:  downlink.Deviceeui, //identifier,
		Slot:       int32(downlink.Slot),
		FieldIndex: int32(downlink.Index),
		NumRetries: int32(downlink.Numretries),
	}

	_, err := c.consistencyService.ProcessCheckConsistencyRequest(checkConsistencyRequest)
	loggerhelper.WriteToLog(err)
}

// UpdateReported update reported config
func (c *Service) UpdateReported(req *pb.UpdateReportedRequest) (*pb.Response, error) {
	docType := nosql.DocTypeConfigSchema

	if req.Slot > 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}

	firmware, err := c.dbClient.GetLatestFirmware(docType)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}
	fieldDetails, err := c.dbClient.GetFieldDetailsByIndex(req.FieldIndex, firmware, docType)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	// update dbclient
	err = c.dbClient.UpdateDbReported(req, fieldDetails)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	/*
		if fieldDetails.Name == "roffset" || fieldDetails.Name == "firmware" {
			// trigger timescale update
			ctx, cancel := helpers.GetContextWithAuth(c.serviceKey)
			defer cancel()
			request := &pbConnection.Identifier{
				Identifier: req.DeviceEUI,
			}
			conn, err := c.grpcConnectionClient.GetConnection(ctx, request)
			if err != nil {
				c.authLoggerHelper.HandleGrpcError(err)

				errMsg := &pbLogger.ErrorMessage{
					Service:  "config-service",
					Function: "UpdateReported",
					Severity: pbLogger.ErrorMessage_FATAL,
					Message:  err.Error(),
				}
				c.errorChan <- errMsg
			}

			updateRequest := &pbConnection.UpdateConnectionRequest{
				Identifier: req.DeviceEUI,
				Connection: conn,
			}
			_, err = c.grpcConnectionClient.UpdateConnection(ctx, updateRequest)
			if err != nil {
				c.authLoggerHelper.HandleGrpcError(err)

				errMsg := &pbLogger.ErrorMessage{
					Service:  "config-service",
					Function: "UpdateReported",
					Severity: pbLogger.ErrorMessage_FATAL,
					Message:  err.Error(),
				}
				c.errorChan <- errMsg

			}
		}*/

	// send response
	return &pb.Response{
		Reply: "OK",
	}, nil

}

// GetConfigByName get config details by name
func (c *Service) GetConfigByName(token string, req *pb.GetConfigByNameRequest) (*pb.ConfigField, error) {
	allowedRoles := []string{c.adminRole, c.installerRole, c.superuserRole}
	_, err := authhelper.CheckToken(c.grpcAuthClient, token, allowedRoles)
	if err != nil {
		return nil, err
	}

	docType := nosql.DocTypeConfigSchema
	if req.Slot > 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}
	firmware, err := c.dbClient.GetLatestFirmware(docType)
	if err != nil {
		return nil, err
	}

	fieldDetails, err := c.dbClient.GetFieldDetailsByName(req.FieldName, firmware, docType)
	if err != nil {
		return nil, err
	}

	configField, err := c.dbClient.GetConfigByName(firmware, fieldDetails, req)
	if err != nil {
		return nil, err
	}

	return configField, nil
}

// GetNewConfigDoc Create blank config doc
func (c *Service) GetNewConfigDoc(token string, req *pb.Identifier) (*pb.ConfigDoc, error) {
	allowedRoles := []string{c.adminRole, c.superuserRole}
	_, err := authhelper.CheckToken(c.grpcAuthClient, token, allowedRoles)
	if err != nil {
		return nil, err
	}

	docType := nosql.DocTypeConfigSchema

	if req.Slot != 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}

	// get latest firmware version
	firmware, err := c.dbClient.GetLatestFirmware(docType)
	if err != nil {
		return nil, err
	}
	loggerhelper.WriteToLog("Got latest firmware: " + firmware)

	// get config schema
	configFields, err := c.dbClient.GetFieldDetails(firmware, docType)
	if err != nil {
		return nil, err
	}

	configDoc := &pb.ConfigDoc{
		Desired:  make(map[string]string),
		Reported: make(map[string]string),
	}

	for _, v := range configFields {
		configDoc.Desired[v.Name] = ""
		configDoc.Reported[v.Name] = ""
	}

	return configDoc, nil
}

// GetConfigByIndex get config by index
func (c *Service) GetConfigByIndex(token string, req *pb.GetConfigByIndexRequest) (*pb.ConfigField, error) {
	allowedRoles := []string{c.adminRole, c.installerRole, c.superuserRole}
	_, err := authhelper.CheckToken(c.grpcAuthClient, token, allowedRoles)
	if err != nil {
		return nil, err
	}

	return c.dbClient.GetConfigByIndex(req)
}

// GetDeviceConfig get all config for a device
func (c *Service) GetDeviceConfig(token string, req *pb.Identifier) (*pb.ConfigFields, error) {
	allowedRoles := []string{c.adminRole, c.installerRole, c.superuserRole}
	_, err := authhelper.CheckToken(c.grpcAuthClient, token, allowedRoles)
	if err != nil {
		return &pb.ConfigFields{}, err
	}

	return c.dbClient.GetDeviceConfig(req)
}

// UpdateFirmwareAllDevices update all devices to new firmware
func (c *Service) UpdateFirmwareAllDevices(token string) error {
	allowedRoles := []string{c.adminRole, c.superuserRole}
	_, err := authhelper.CheckToken(c.grpcAuthClient, token, allowedRoles)
	if err != nil {
		return err
	}

	return c.dbClient.UpdateFirmwareAllDevices()

}

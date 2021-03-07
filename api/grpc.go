package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/sukhajata/devicetwin/pkg/authhelper"
	"github.com/sukhajata/devicetwin/pkg/loggerhelper"

	"github.com/sukhajata/devicetwin/internal/consistency"
	"github.com/sukhajata/devicetwin/internal/core"
	pb "github.com/sukhajata/ppconfig"
	pbLogger "github.com/sukhajata/pplogger"
)

// GRPCServer implements ppconfig.ConfigServiceServer
type GRPCServer struct {
	configService      core.ConfigHandler
	consistencyService *consistency.Service
	loggerHelper       loggerhelper.Helper
	pb.UnimplementedConfigServiceServer
}

// NewGRPCConfigServer factory method
func NewGRPCConfigServer(configService core.ConfigHandler, consistencyService *consistency.Service, loggerHelper loggerhelper.Helper) *GRPCServer {
	return &GRPCServer{
		configService:      configService,
		consistencyService: consistencyService,
		loggerHelper:       loggerHelper,
	}
}

// GetNewConfigDoc Create blank config doc
func (s *GRPCServer) GetNewConfigDoc(ctx context.Context, req *pb.Identifier) (*pb.ConfigDoc, error) {
	token, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.configService.GetNewConfigDoc(token, req)
}

//SetDesired set a desired config value
func (s *GRPCServer) SetDesired(ctx context.Context, req *pb.SetDesiredRequest) (*pb.Response, error) {
	loggerhelper.WriteToLog(fmt.Sprintf("Received set desired request %v %v", req.Identifier, req.FieldName))
	token, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	return s.configService.SetDesired(token, req)

}

//CheckConsistency check consistency for a field
func (s *GRPCServer) CheckConsistency(ctx context.Context, req *pb.CheckConsistencyRequest) (*pb.Response, error) {
	_, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	return s.consistencyService.ProcessCheckConsistencyRequest(req)

}

//UpdateReported update reported config
func (s *GRPCServer) UpdateReported(ctx context.Context, req *pb.UpdateReportedRequest) (*pb.Response, error) {
	_, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return &pb.Response{
			Reply: "UNAUTHORIZED",
		}, err
	}

	if req.GetDeviceEUI() == "" {
		err = errors.New("missing identifier")
		s.loggerHelper.LogError("UpdateReported", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	if req.GetFieldIndex() == 0 {
		err = errors.New("missing field index")
		s.loggerHelper.LogError("UpdateReported", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	/*if len(req.GetFieldValue()) == 0 {
		helpers.LogError("UpdateReported3", "Missing field value", pbLogger.ErrorMessage_SEVERE)
		return &pb.Response{
			Reply: "NOT OK",
		}, errors.New("Missing field value")
	}*/

	return s.configService.UpdateReported(req)
}

// UpdateFirmware update descired firmware for all devices
func (s *GRPCServer) UpdateFirmware(ctx context.Context, req *pb.UpdateFirmwareRequest) (*pb.Response, error) {
	token, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// start on new goroutine so we can return quickly
	go func(token string, loggerHelper loggerhelper.Helper) {
		err := s.configService.UpdateFirmwareAllDevices(token)
		if err != nil {
			loggerHelper.LogError("UpdateFirmware", err.Error(), pbLogger.ErrorMessage_FATAL)
		}
	}(token, s.loggerHelper)

	return &pb.Response{
		Reply: "Started updating",
	}, nil
}

// GetConfigByName get config by name
func (s *GRPCServer) GetConfigByName(ctx context.Context, req *pb.GetConfigByNameRequest) (*pb.ConfigField, error) {
	token, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.configService.GetConfigByName(token, req)

}

// GetConfigByIndex get config by index
func (s *GRPCServer) GetConfigByIndex(ctx context.Context, req *pb.GetConfigByIndexRequest) (*pb.ConfigField, error) {
	token, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.configService.GetConfigByIndex(token, req)
}

// GetAllConfig get all config for a device
func (s *GRPCServer) GetAllConfig(ctx context.Context, req *pb.Identifier) (*pb.ConfigFields, error) {
	token, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	results, err := s.configService.GetDeviceConfig(token, req)
	if err != nil {
		loggerhelper.WriteToLog(err.Error())
	}

	return results, err
}

//AssignRadioOffset (Identifier) returns (Response) {}
func (s *GRPCServer) AssignRadioOffset(ctx context.Context, req *pb.Identifier) (*pb.Response, error) {
	token, err := authhelper.GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.configService.AssignRadioOffset(token, req)
}

package consistency

import (
	"fmt"
	db2 "github.com/sukhajata/devicetwin.git/internal/dbclient"
	"github.com/sukhajata/devicetwin.git/internal/dbclient/nosql"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/sukhajata/devicetwin.git/internal/dataapi"
	"github.com/sukhajata/devicetwin.git/internal/types"
	"github.com/sukhajata/devicetwin.git/internal/utility"
	"github.com/sukhajata/devicetwin.git/pkg/loggerhelper"
	pb "github.com/sukhajata/ppconfig"
	pbLogger "github.com/sukhajata/pplogger"
	"github.com/sukhajata/ppmessage/ppdownlink"
)

type ConsistencyChecker interface {
	ScheduleConsistencyCheckForField(req *pb.Identifier, fieldDetails types.ConfigFieldDetails, firmware string, numRetries int32)
	CheckConsistencyForField(field *pb.ConfigField, firmware string, req *pb.Identifier) error
	ProcessCheckConsistencyRequest(req *pb.CheckConsistencyRequest) (*pb.Response, error)
	ScheduleMessageSend(identifier string, downlink *ppdownlink.ConfigDownlinkMessage)
	RunScheduledConsistencyCheck()
	CheckConsistencyAllFieldsForDevice(req *pb.Identifier)
}

// Service provides config consistency checking
type Service struct {
	transmitChannel     chan<- *ppdownlink.ConfigDownlinkMessage
	dbClient            db2.Client
	dataAPIClient       dataapi.Client
	repeatCheckSchedule string
	loggerHelper        loggerhelper.Helper
}

// NewService factory method
func NewService(
	dbClient db2.Client,
	dataAPIClient dataapi.Client,
	transmitChannel chan<- *ppdownlink.ConfigDownlinkMessage,
	repeatCheckSchedule string,
	loggerHelper loggerhelper.Helper,
) *Service {
	return &Service{
		transmitChannel:     transmitChannel,
		dbClient:            dbClient,
		dataAPIClient:       dataAPIClient,
		repeatCheckSchedule: repeatCheckSchedule,
		loggerHelper:        loggerHelper,
	}
}

// ScheduleConsistencyCheckForField schedule a consistency check
func (s *Service) ScheduleConsistencyCheckForField(req *pb.Identifier, fieldDetails types.ConfigFieldDetails, firmware string, numRetries int32) {
	var timer1 *time.Timer

	if fieldDetails.Name == "installd" {
		// s11 command
		if numRetries == 0 {
			timer1 = time.NewTimer(30 * time.Second)
		} else if numRetries == 1 {
			timer1 = time.NewTimer(120 * time.Second)
		} else if numRetries == 2 {
			timer1 = time.NewTimer(240 * time.Second)
		} else {
			//leave to scheduled check
			return
		}
	} else {
		// this will run after the first message has been sent
		// scheduling involves waiting some period of time, then checking consistency, then resending if necessary. Resends will be sent in dlresmin
		// repeatCheckSchedule is expected to be a string describing the seconds to wait before retrying, separated by _, such as "25_540_3420"
		schedule := strings.Split(s.repeatCheckSchedule, "_")
		firstDelay, err := strconv.Atoi(schedule[0])
		if err != nil {
			s.loggerHelper.LogError("scheduleConsistencyCheckForField", err.Error(), pbLogger.ErrorMessage_SEVERE)
			firstDelay = 25
		}
		secondDelay, err := strconv.Atoi(schedule[1])
		if err != nil {
			s.loggerHelper.LogError("scheduleConsistencyCheckForField", err.Error(), pbLogger.ErrorMessage_SEVERE)
			secondDelay = 540
		}
		thirdDelay, err := strconv.Atoi(schedule[2])
		if err != nil {
			s.loggerHelper.LogError("scheduleConsistencyCheckForField", err.Error(), pbLogger.ErrorMessage_SEVERE)
			thirdDelay = 3420 // 57 minutes
		}
		// time to wait depends on numRetries
		if numRetries == 0 {
			timer1 = time.NewTimer(time.Duration(firstDelay) * time.Second)
		} else if numRetries < 3 {
			timer1 = time.NewTimer(time.Duration(secondDelay) * time.Second)
		} else if numRetries < 7 {
			timer1 = time.NewTimer(time.Duration(thirdDelay) * time.Second)
		} else {
			// leave to scheduled check
			return
		}
	}

	// wait
	<-timer1.C

	// check consistency
	configByNameRequest := pb.GetConfigByNameRequest{
		Identifier: req.Identifier,
		FieldName:  fieldDetails.Name,
		Slot:       req.Slot,
	}
	result, err := s.dbClient.GetConfigByName(firmware, fieldDetails, &configByNameRequest)
	if err != nil {
		s.loggerHelper.LogError("scheduleConsistencyCheckForField2", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return
	}

	// only do something if the desired field has been set, and does not match the reported
	if result.Desired != "" && result.Desired != result.Reported {
		downlink, err := utility.BuildDownlinkMessage(req.Identifier, fieldDetails, utility.GetFormattedValue(result.Desired), firmware, numRetries+1, uint32(req.Slot))
		if err != nil {
			s.loggerHelper.LogError("scheduleConsistencyCheckForField3", err.Error(), pbLogger.ErrorMessage_SEVERE)
			return
		}

		if numRetries > 1 {
			// send in dlresmin
			s.ScheduleMessageSend(req.Identifier, downlink)
		} else {
			// s11 control. Send immediately
			s.Send(downlink)
		}
		loggerhelper.WriteToLog(fmt.Sprintf("Resent message %s: retries: %d\n", fieldDetails.Name, numRetries+1))
	} else {
		loggerhelper.WriteToLog("No need to resend message")
	}

}

// CheckConsistencyForField - compare desired vs reported in the dbclient and schedule a message send if inconsistent
func (s *Service) CheckConsistencyForField(field *pb.ConfigField, firmware string, req *pb.Identifier) error {
	docType := nosql.DocTypeConfigSchema

	if req.Slot > 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}

	if field.Desired != "" && field.Desired != field.Reported {
		//mismatch. send a message
		//will fail if device does not support this config field
		fieldDetails, err := s.dbClient.GetFieldDetailsByIndex(field.Index, firmware, docType)
		if err != nil {
			s.loggerHelper.LogError("checkConsistencyForField", err.Error(), pbLogger.ErrorMessage_SEVERE)
			return err
		}

		loggerhelper.WriteToLog(fmt.Sprintf("Config %v with mismatch for %s slot %v, scheduling message, setting to %v", field.Name, req.Identifier, req.Slot, field.Desired))

		downlink, err := utility.BuildDownlinkMessage(req.Identifier, fieldDetails, field.Desired, firmware, 0, uint32(req.Slot))
		if err != nil {
			s.loggerHelper.LogError("checkConsistencyForField", err.Error(), pbLogger.ErrorMessage_SEVERE)
			return err
		}

		go s.ScheduleMessageSend(req.Identifier, downlink)
	}

	return nil
}

// ProcessCheckConsistencyRequest check consistency
func (s *Service) ProcessCheckConsistencyRequest(req *pb.CheckConsistencyRequest) (*pb.Response, error) {
	docType := nosql.DocTypeConfigSchema

	if req.Slot > 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}

	firmware, err := s.dbClient.GetLatestFirmware(docType)
	if err != nil {
		s.loggerHelper.LogError("CheckConsistency2", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return nil, err
	}

	fieldDetails, err := s.dbClient.GetFieldDetailsByIndex(req.FieldIndex, firmware, docType)
	if err != nil {
		s.loggerHelper.LogError("CheckConsistency3", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return nil, err
	}

	identifier := &pb.Identifier{
		Identifier: req.DeviceEUI,
		Slot:       req.Slot,
	}
	go s.ScheduleConsistencyCheckForField(identifier, fieldDetails, firmware, req.NumRetries)

	return &pb.Response{
		Reply: "OK",
	}, nil

}

// ScheduleMessageSend - send message in dlresmin
func (s *Service) ScheduleMessageSend(identifier string, downlink *ppdownlink.ConfigDownlinkMessage) {
	//when to send message?
	reservedMinutes, err := s.dbClient.GetDLResmin(identifier)
	if err != nil {
		s.loggerHelper.LogError("scheduleMessageSend1", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return
	}

	//find min and max times to send
	min := 10
	max := 0

	minutesArray := strings.Split(reservedMinutes, ",")

	for _, e := range minutesArray {
		element, err := strconv.Atoi(e)
		if err != nil {
			min = 6
			max = 8
			break
		}

		if element < min {
			min = element
		}

		if element > max {
			max = element
		}
	}

	//figure out how long to wait before sending
	_, currentMinutes, currentSeconds := time.Now().Clock()
	outOfTen := currentMinutes % 10
	loggerhelper.WriteToLog(fmt.Sprintf("min = %v, max = %v, current minute %v", min, max, outOfTen))
	var addMinutes int
	var addSeconds int
	rand.Seed(time.Now().UnixNano())
	if outOfTen == min || outOfTen == max {
		//we are currently in the sending window
		addMinutes = 0
		addSeconds = 0
	} else if outOfTen < min {
		loggerhelper.WriteToLog("current minute is less than min")
		//eg. current minute == 3, dlresmin == 6,8
		addSeconds = 60 - currentSeconds + (rand.Intn(60))
		addMinutes = min - outOfTen - 1
	} else if outOfTen < max {
		loggerhelper.WriteToLog("current minute is less than max")
		//eg. current minute == 7, dlresmin == 6,8
		addSeconds = 60 - currentSeconds + (rand.Intn(60))
		addMinutes = max - outOfTen - 1
	} else {
		loggerhelper.WriteToLog("current minute is greater than max")
		//eg. current minute == 9, dlresmin == 6,8
		addSeconds = 60 - currentSeconds + (rand.Intn(60))
		addMinutes = min + (10 - outOfTen) - 1
	}

	timeToSend := time.Now().Add(time.Minute * time.Duration(addMinutes)).Add(time.Second * time.Duration(addSeconds))

	timeToWait := timeToSend.Sub(time.Now())

	loggerhelper.WriteToLog(fmt.Sprintf("Going to wait %v", timeToWait))
	timer2 := time.NewTimer(timeToWait)
	<-timer2.C
	loggerhelper.WriteToLog("Thanks for your patience")
	s.Send(downlink)

}

func (s *Service) Send(downlink *ppdownlink.ConfigDownlinkMessage) {
	s.transmitChannel <- downlink

	// check consistency
	checkConsistencyRequest := &pb.CheckConsistencyRequest{
		DeviceEUI:  downlink.Deviceeui, //identifier,
		Slot:       int32(downlink.Slot),
		FieldIndex: int32(downlink.Index),
		NumRetries: int32(downlink.Numretries),
	}
	_, err := s.ProcessCheckConsistencyRequest(checkConsistencyRequest)
	if err != nil {
		s.loggerHelper.LogError("Send", err.Error(), pbLogger.ErrorMessage_SEVERE)
	}
}

// RunScheduledConsistencyCheck for non S11, periodic check of consistency of desired and reported config for all devices
func (s *Service) RunScheduledConsistencyCheck() {
	loggerhelper.WriteToLog("Running scheduled consistency check")
	inconsistent, err := s.dbClient.GetInconsistentDevices()
	if err != nil {
		s.loggerHelper.LogError("runScheduledConsistencyCheck", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return
	}

	for _, v := range inconsistent {
		minsSinceLastMsg, err := s.dataAPIClient.GetMinsSinceLastMsg(v)
		if err != nil {
			continue
		}
		// don't send to dead devices
		if minsSinceLastMsg < 40 {
			s.CheckConsistencyAllFieldsForDevice(&pb.Identifier{
				Identifier: v,
				Slot:       0,
			})
			// important! don't overload couchbase
			time.Sleep(time.Second * 2)
		}

	}
}

// CheckConsistencyAllFieldsForDevice check consistency for device
func (s *Service) CheckConsistencyAllFieldsForDevice(req *pb.Identifier) {
	docType := nosql.DocTypeConfigSchema

	if req.Slot > 0 {
		docType = nosql.DocTypeS11ConfigSchema
	}

	firmware, err := s.dbClient.GetLatestFirmware(docType)
	if err != nil {
		s.loggerHelper.LogError("CheckConsistencyAllFieldsForDevice", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return
	}

	configFields, err := s.dbClient.GetDeviceConfig(req)
	if err != nil {
		s.loggerHelper.LogError("CheckConsistencyAllFieldsForDevice", err.Error(), pbLogger.ErrorMessage_SEVERE)
		return
	}

	for _, field := range configFields.Fields {
		err = s.CheckConsistencyForField(field, firmware, req)
		if err != nil {
			s.loggerHelper.LogError("CheckConsistencyAllFieldsForDevice", err.Error(), pbLogger.ErrorMessage_SEVERE)
			return
		}
	}
}

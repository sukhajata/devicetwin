package utility

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"github.com/sukhajata/devicetwin.git/internal/types"
	"github.com/sukhajata/devicetwin.git/pkg/loggerhelper"
	pb "github.com/sukhajata/ppconfig"
	"github.com/sukhajata/ppmessage/ppdownlink"
)

//StringToInterface convert string to int or bool
func StringToInterface(fieldDetails types.ConfigFieldDetails, rawValue string) (interface{}, error) {

	if fieldDetails.Type == "i" || fieldDetails.Type == "t" {
		intValue, err := strconv.Atoi(rawValue)
		if err != nil {
			loggerhelper.WriteToLog("format error " + rawValue)
			return "", err
		}
		return intValue, nil
	} else if fieldDetails.Type == "b" {
		if rawValue == "1" {
			return 1, nil
		}
		return 0, nil
	} else {
		return rawValue, nil
	}

}

//DecodeFieldValue convert byte array to appropriate type
func DecodeFieldValue(fieldDetails types.ConfigFieldDetails, rawValue []byte) (result interface{}, err error) {
	reader := bytes.NewReader(rawValue)
	switch fieldDetails.Type {
	case "i":
		//int
		var intValue int32
		err := binary.Read(reader, binary.BigEndian, &intValue)
		if err != nil {
			return "", err
		}
		return intValue, nil
	case "t":
		//short
		var intValue int16
		err := binary.Read(reader, binary.BigEndian, &intValue)
		if err != nil {
			return "", err
		}
		return intValue, nil
	case "b":
		//2 byte bool
		//just convert to int16
		var intValue int16
		err := binary.Read(reader, binary.BigEndian, &intValue)
		if err != nil {
			return "", err
		}
		return intValue, nil
	default:
		//trim padding
		trimmed := bytes.Trim(rawValue, "\x00")
		return string(trimmed), nil
	}

}

//BuildDownlinkMessage build a config downlink message with the given details, checking that the value is valid
func BuildDownlinkMessage(deviceEUI string, fieldDetails types.ConfigFieldDetails, fieldValue string, firmwareVersion string, numRetries int32, slot uint32) (*ppdownlink.ConfigDownlinkMessage, error) {
	downlink := &ppdownlink.ConfigDownlinkMessage{
		Deviceeui:  deviceEUI,
		Slot:       slot,
		Index:      uint32(fieldDetails.Index),
		Firmware:   firmwareVersion,
		Numretries: uint32(numRetries),
	}

	// next write the value of the field in binary
	buf := new(bytes.Buffer)
	switch fieldDetails.Type {
	case "i":
		// 4 byte signed int
		intValue, err := strconv.Atoi(fieldValue)
		if err != nil {
			return downlink, err
		}
		// check range
		if fieldDetails.Min != nil {
			var minValInt int
			set := false

			switch v := fieldDetails.Min.(type) {
			case float64:
				minValInt = int(v)
				set = true
			case int32:
				minValInt = int(v)
				set = true
			case int:
				minValInt = v
				set = true
			case string:
				if v != "" {
					try, err := strconv.Atoi(v)
					if err != nil {
						loggerhelper.WriteToLog(fmt.Sprintf("Failed to convert %v min %v to int", fieldDetails.Name, fieldDetails.Min))
					} else {
						minValInt = int(try)
						set = true
					}
				}
			default:
				set = false
			}

			if set && intValue < minValInt {
				return downlink, fmt.Errorf("Value %d below minimum allowed %d", intValue, minValInt)
			}
		}
		if fieldDetails.Max != nil {
			var maxValInt int
			set := false

			switch v := fieldDetails.Max.(type) {
			case float64:
				maxValInt = int(v)
				set = true
			case int32:
				maxValInt = int(v)
				set = true
			case int:
				maxValInt = v
				set = true
			case string:
				if v != "" {
					try, err := strconv.Atoi(v)
					if err != nil {
						loggerhelper.WriteToLog(fmt.Sprintf("Failed to convert %v max %v to int", fieldDetails.Name, fieldDetails.Max))
					} else {
						maxValInt = int(try)
						set = true
					}
				}
			default:
				set = false
			}

			if set && intValue > maxValInt {
				return downlink, fmt.Errorf("Value %d above maximum allowed %d", intValue, maxValInt)
			}
		}

		err = binary.Write(buf, binary.BigEndian, int32(intValue))
		if err != nil {
			return downlink, err
		}
	case "t":
		// 2 byte signed int
		// ParseInt returns int64 but you can specify that it should fit into int16
		intValue, err := strconv.Atoi(fieldValue)
		if err != nil {
			return downlink, err
		}
		// check range
		if fieldDetails.Min != nil {
			var minValInt int
			set := false

			switch v := fieldDetails.Min.(type) {
			case float64:
				minValInt = int(v)
				set = true
			case int32:
				minValInt = int(v)
				set = true
			case int:
				minValInt = v
				set = true
			case string:
				if v != "" {
					try, err := strconv.Atoi(v)
					if err != nil {
						loggerhelper.WriteToLog(fmt.Sprintf("Failed to convert %v min %v to int", fieldDetails.Name, fieldDetails.Min))
						minValInt = 0
					} else {
						minValInt = int(try)
						set = true
					}
				}
			default:
				set = false
			}

			if set && intValue < minValInt {
				return downlink, fmt.Errorf("Value %d below minimum allowed %d", intValue, minValInt)
			}
		}
		if fieldDetails.Max != nil {
			var maxValInt int
			set := false

			switch v := fieldDetails.Max.(type) {
			case float64:
				maxValInt = int(v)
				set = true
			case int32:
				maxValInt = int(v)
				set = true
			case int:
				maxValInt = v
				set = true
			case string:
				if v != "" {
					try, err := strconv.Atoi(v)
					if err != nil {
						loggerhelper.WriteToLog(fmt.Sprintf("Failed to convert %v max %v to int", fieldDetails.Name, fieldDetails.Max))
						maxValInt = 0
					} else {
						maxValInt = int(try)
						set = true
					}
				}
			default:
				set = false
			}

			if set && intValue > maxValInt {
				return downlink, fmt.Errorf("Value %d above maximum allowed %d", intValue, maxValInt)
			}
		}

		int16Value := int16(intValue)
		err = binary.Write(buf, binary.BigEndian, int16Value)
		if err != nil {
			return downlink, err
		}
	case "b":
		// 2 byte bool
		boolValue, err := strconv.ParseBool(fieldValue)
		if err != nil {
			return downlink, err
		}
		// convert to 2 byte int
		var int16Value int16
		int16Value = 0
		if boolValue {
			int16Value = 1
		}
		err = binary.Write(buf, binary.BigEndian, int16Value)
		if err != nil {
			return downlink, err
		}
	default:
		// Type will be max length of string
		var length int
		set := false

		switch v := fieldDetails.Type.(type) {
		case float64:
			length = int(v)
			set = true
		case int:
			length = v
			set = true
		case string:
			if v != "" {
				intVal, err := strconv.Atoi(v)
				if err != nil {
					loggerhelper.WriteToLog(fmt.Sprintf("Failed to convert %v to int", intVal))
				} else {
					length = intVal
					set = true
				}
			}
		default:
			set = false
		}

		buf.Write([]byte(fieldValue))

		if set {
			// check length
			if buf.Len() > length {
				return downlink, errors.New("String too long for " + fieldDetails.Name + ", length: " +
					strconv.Itoa(buf.Len()) + ", allowed: " + strconv.Itoa(length))
			}
		}

		// padding
		diff := length - buf.Len()
		for i := 0; i < diff; i++ {
			buf.WriteByte(byte(0))
		}
	}

	//configMsg["value"] = buf.Bytes()
	//configMsg["numRetries"] = numRetries
	downlink.Value = buf.Bytes()

	return downlink, nil
}

// GetFormattedValue get a formatted string
func GetFormattedValue(val interface{}) string {
	switch val.(type) {
	case int, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Find find a config field in an array
func Find(fields []*pb.ConfigField, name string) *pb.ConfigField {
	for _, v := range fields {
		if v.GetName() == name {
			return v
		}
	}

	return nil
}

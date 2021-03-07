package utility

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sukhajata/devicetwin/internal/types"
)

func Test_StringToInterface_String(t *testing.T) {
	fieldDetails := types.ConfigFieldDetails{
		Index: 46,
		Name:  "firmware",
		Type:  32.0,
	}

	value, err := StringToInterface(fieldDetails, "1.1.4")
	s, ok := value.(string)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, "1.1.4", s)
}

func Test_StringToInterface_Int(t *testing.T) {
	fieldDetails := types.ConfigFieldDetails{
		Index: 3,
		Name:  "roffset",
		Type:  "i",
	}

	value, err := StringToInterface(fieldDetails, "2500")
	s, ok := value.(int)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, 2500, s)
}

func Test_DecodeFieldValue_Int(t *testing.T) {
	fieldDetails := types.ConfigFieldDetails{
		Index: 3,
		Name:  "roffset",
		Type:  "i",
	}

	result, err := DecodeFieldValue(fieldDetails, []byte{0x00, 0x00, 0x09, 0xc4})
	s, ok := result.(int32)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, int32(2500), s)
}

func Test_DecodeFieldValue_String(t *testing.T) {
	fieldDetails := types.ConfigFieldDetails{
		Index: 6,
		Name:  "mqttuser",
		Type:  "8",
	}

	result, err := DecodeFieldValue(fieldDetails, []byte{0x72, 0x75, 0x70, 0x65, 0x72, 0x74})
	s, ok := result.(string)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, "rupert", s)
}

func Test_BuildDownlinkMessage_String(t *testing.T) {
	fieldDetails := types.ConfigFieldDetails{
		Index: 6,
		Name:  "mqttuser",
		Type:  8.0,
	}
	fieldValue := "rupert"
	firmwareVersion := "1.1.4"

	downlink, err := BuildDownlinkMessage("123", fieldDetails, fieldValue, firmwareVersion, 0, 0)
	require.Nil(t, err)
	require.Equal(t, []byte{0x72, 0x75, 0x70, 0x65, 0x72, 0x74, 0x00, 0x00}, downlink.Value)
	require.Equal(t, uint32(6), downlink.Index)
}

func Test_BuildDownlinkMessage_Int(t *testing.T) {
	fieldDetails := types.ConfigFieldDetails{
		Index: 3,
		Name:  "roffset",
		Type:  "i",
		Min:   0,
		Max:   3000,
	}
	fieldValue := "2500"
	firmwareVersion := "1.1.4"

	downlink, err := BuildDownlinkMessage("123", fieldDetails, fieldValue, firmwareVersion, 0, 0)
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x09, 0xc4}, downlink.Value)
	require.Equal(t, uint32(3), downlink.Index)
}

func Test_BuildDownlinkMessage_Short(t *testing.T) {
	fieldDetails := types.ConfigFieldDetails{
		Index: 8,
		Name:  "mqttssl",
		Type:  "t",
		Min:   0,
		Max:   1,
	}
	fieldValue := "1"
	firmwareVersion := "1.1.4"

	downlink, err := BuildDownlinkMessage("123", fieldDetails, fieldValue, firmwareVersion, 0, 0)
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x01}, downlink.Value)
	require.Equal(t, uint32(8), downlink.Index)
}

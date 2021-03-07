package dataapi

import (
	"encoding/json"
	"fmt"
	"github.com/sukhajata/devicetwin/pkg/loggerhelper"
	"net/http"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client for data API
type Client struct {
	dataServiceAddress      string
	dataToken               string
	connectionsDataViewName string
	HTTPClient              HTTPClient
}

// NewClient factory method
func NewClient(dataServiceAddress string, dataToken string, connectionsDataViewName string, httpClient HTTPClient) Client {
	return Client{
		dataServiceAddress:      dataServiceAddress,
		dataToken:               dataToken,
		connectionsDataViewName: connectionsDataViewName,
		HTTPClient:              httpClient,
	}
}

// GetMinsSinceLastMsg mins since the last message was received
func (c *Client) GetMinsSinceLastMsg(deviceEUI string) (int32, error) {
	url := fmt.Sprintf("%s/%s?select=LASTRECEIVEDMESSAGE&DEVICEEUI=eq.%s", c.dataServiceAddress, c.connectionsDataViewName, deviceEUI)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.dataToken))

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return -1, err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			loggerhelper.WriteToLog(err.Error())
		}
	}()

	decoder := json.NewDecoder(response.Body)
	var data []map[string]interface{}
	err = decoder.Decode(&data)
	if err != nil {
		return -1, err
	}

	if len(data) == 0 {
		return -1, fmt.Errorf("No data for %s", deviceEUI)
	}

	value, ok := data[0]["LASTRECEIVEDMESSAGE"].(float64)
	if ok != true {
		return -1, fmt.Errorf("Failed to cast %v as float", data[0]["LASTRECEIVEDMESSAGE"])
	}

	return int32(value), nil
}

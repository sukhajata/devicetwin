package ppmqtt

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/sukhajata/devicetwin/pkg/loggerhelper"
)

// Client represents an mqtt client
type Client interface {
	Publish(message Message) error
	Subscribe(topic string) error
}

// PPClient implements Client
type PPClient struct {
	mqttClient       mqtt.Client
	ReceiveChan      chan Message
	ErrorChan        chan error
	callback         mqtt.MessageHandler
	onConnectionLost mqtt.ConnectionLostHandler
}

// Message represent an ppmqtt message
type Message struct {
	Topic   string
	Payload []byte
}

// Create tls.Config with desired tls properties
func getTLSConfig() *tls.Config {
	return &tls.Config{
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// NewClient - factory method for ppmqtt.Client
func NewClient(mqttBroker string, mqttUsername string, mqttPassword string, serviceName string) (*PPClient, error) {
	c := &PPClient{
		ReceiveChan: make(chan Message, 2),
		ErrorChan:   make(chan error, 2),
	}

	// callback for when a message is received
	c.callback = func(client mqtt.Client, msg mqtt.Message) {
		loggerhelper.WriteToLog(fmt.Sprintf("TOPIC: %s\n", msg.Topic()))
		c.ReceiveChan <- Message{msg.Topic(), msg.Payload()}
	}

	// send error on error chan when connection is lost
	c.onConnectionLost = func(client mqtt.Client, err error) {
		c.ErrorChan <- err
		client.Disconnect(0)
	}

	rand.Seed(time.Now().UnixNano())

	tlsConfig := getTLSConfig()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(serviceName + randStringRunes(8))
	opts.SetCleanSession(true)
	opts.SetUsername(mqttUsername)
	opts.SetPassword(mqttPassword)
	opts.SetDefaultPublishHandler(c.callback)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(2 * time.Second)
	opts.SetConnectionLostHandler(c.onConnectionLost)
	//opts.SetAutoReconnect(true)

	c.mqttClient = mqtt.NewClient(opts)

	failCount := 0
	for {
		if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
			failCount++
			if failCount > 5 {
				return c, token.Error()
			}
			fmt.Println(token.Error())
			time.Sleep(time.Second * 2)
			continue
		}
		break
	}

	return c, nil
}

// Subscribe to the given topic
func (c *PPClient) Subscribe(topic string) error {
	if c.mqttClient == nil {
		return errors.New("mqtt client not initialized")
	}
	// subscribe to the topic and request messages to be delivered
	// at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("subscribed to %s\n", topic)

	return nil

}

// Publish to mqtt
func (c *PPClient) Publish(message Message) error {
	if c.mqttClient == nil {
		return errors.New("mqtt client not initialized")
	}
	token := c.mqttClient.Publish(message.Topic, 0, false, message.Payload)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

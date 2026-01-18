package mqtt

import (
	"DataManager/internal/pb"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

const (
	broker   = "tcp://localhost:1884"
	clientID = "iot-mqtt-client"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func InitMqtt() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func PublishMessage(topic string, message string) error {
	token := mqttClient.Publish(topic, 0, false, message)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Message '%s' published!\n", message)
	return nil
}

func PublishReading(topic string, reading *pb.Reading) error {
	token := mqttClient.Publish(topic, 0, false, reading)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Reading '%+v' published!\n", reading)
	return nil
}

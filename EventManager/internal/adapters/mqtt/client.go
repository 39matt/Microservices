package mqtt

import (
	"EventManager/internal/data/limits"
	"EventManager/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

const (
	broker   = "tcp://mosquitto:1883"
	clientID = "iot-mqtt-subscriber"
)

var mqttMsgChan = make(chan mqtt.Message)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mqttMsgChan <- msg
}

func processMsg(ctx context.Context, input <-chan mqtt.Message) chan mqtt.Message {
	out := make(chan mqtt.Message)
	go func() {
		defer close(out)
		for {
			select {
			case msg, ok := <-input:
				if !ok {
					return
				}
				out <- msg
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func InitMqtt(ctx context.Context) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(messagePubHandler)
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

func PublishReading(topic string, reading []byte) error {
	token := mqttClient.Publish(topic, 0, false, reading)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Reading published to topic '%s'!\n", topic)
	return nil
}

func SubscribeToTopic(topic string) error {
	token := mqttClient.Subscribe(topic, 0, nil)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
	return nil
}

func UnsubscribeFromTopic(topic string) error {
	fmt.Println("Unsubscribing and disconnecting...")
	token := mqttClient.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	mqttClient.Disconnect(250)
	return nil
}

func WatchForPublications(ctx context.Context, topic string) error {
	if err := SubscribeToTopic(topic); err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		finalChan := processMsg(ctx, mqttMsgChan)
		for msg := range finalChan {
			if err := MessageHandler(msg); err != nil {
				slog.Error("Handler failed", "err", err)
				continue
			}
		}
	}()

	<-ctx.Done()
	if err := UnsubscribeFromTopic(topic); err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func MessageHandler(msg mqtt.Message) error {
	var reading *models.Reading
	var err error
	if reading, err = models.ConvertFromBytes(msg.Payload()); err != nil {
		return err
	}
	fmt.Printf("Received final message **%s** from topic **%s**\n", reading.ID, msg.Topic())

	if reading.Co > limits.CoLimit || reading.Lpg > limits.LpgLimit || reading.Smoke > limits.SmokeLimit ||
		reading.Humidity > limits.HumidityLimit || reading.Temperature > limits.TemperatureLimit {
		var data []byte
		data, err = json.Marshal(reading)
		if err != nil {
			return fmt.Errorf("marshalling failed: %w", err)
		}

		if err = PublishReading("/limit", data); err != nil {
			return err
		}
	}
	return nil
}

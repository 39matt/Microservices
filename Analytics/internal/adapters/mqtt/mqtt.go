package mqtt

import (
	"Analytics/internal/models"
	"context"
	"fmt"
	"log/slog"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker   = "tcp://mosquitto:1883"
	clientID = "iot-mqtt-subscriber-analytics"
)

var mqttClient mqtt.Client
var mqttMsgChan = make(chan mqtt.Message)
var readingsChan = make(chan []models.Reading)

//var readingsMu = new(sync.RWMutex)

func InitMqtt(readingsChannel chan []models.Reading) error {
	readingsChan = readingsChannel
	//readingsMu = readingsMutex

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

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mqttMsgChan <- msg
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
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

func MessageHandler(msg mqtt.Message) error {
	var reading *models.Reading
	var err error
	if reading, err = models.ConvertFromBytes(msg.Payload()); err != nil {
		return err
	}
	fmt.Printf("Received final message **%s** from topic **%s**\n", reading.ID, msg.Topic())

	readingsChan <- []models.Reading{*reading}
	return nil
}

//
//WatchForPublications increments WaitGroup and launches goroutine A
//
//Goroutine A calls processMsg, which launches goroutine B and returns finalChan immediately
//
//Goroutine A ranges over finalChan until goroutine B closes it
//
//When goroutine B exits (via ctx.Done() or input closed), it closes finalChan, which ends goroutine A's loop
//
//Goroutine A returns, defer wg.Done() runs, WaitGroup reaches zero, wg.Wait() unblocks

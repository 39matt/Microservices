package mqtt

import (
	"EventManager/internal/models"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
				fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
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

	childCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		finalChan := processMsg(childCtx, mqttMsgChan)
		for msg := range finalChan {
			if err := MessageHandler(msg); err != nil {
				slog.Error("Handler failed", "err", err)
				continue
			}
		}
	}()

	if err := SubscribeToTopic("/readings"); err != nil {
		cancel()
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-sigChan:
		fmt.Printf("Signal received: %v\n", sig)
	case <-ctx.Done():
		fmt.Println("Parent context cancelled")
	}

	cancel()
	wg.Wait()

	if err := UnsubscribeFromTopic("/readings"); err != nil {
		return err
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

func MessageHandler(msg mqtt.Message) error {
	var reading *models.Reading
	var err error
	if reading, err = models.ConvertFromBytes(msg.Payload()); err != nil {
		return err
	}
	fmt.Printf("Received final message **%s** from topic **%s**", reading.ID, msg.Topic())
	return nil
}

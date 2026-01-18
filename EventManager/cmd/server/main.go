package main

import (
	"EventManager/internal/adapters/mqtt"
	"context"
	"log"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mqtt.InitMqtt(ctx); err != nil {
		log.Fatal(err)
	}
}

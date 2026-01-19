package main

import (
	"EventManager/internal/adapters/mqtt"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mqtt.InitMqtt(ctx); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := mqtt.WatchForPublications(ctx, "/readings"); err != nil {
			slog.Error("Watch failed", "err", err)
		}
	}()

	select {
	case sig := <-sigChan:
		fmt.Printf("Signal received: %v\n", sig)
	case <-ctx.Done():
		fmt.Println("Context cancelled")
	}

}

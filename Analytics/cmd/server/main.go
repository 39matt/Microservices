package main

import (
	"Analytics/internal/adapters/http"
	"Analytics/internal/adapters/mqtt"
	"Analytics/internal/adapters/nats"
	"Analytics/internal/models"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var readings []models.Reading
var readingsChan = make(chan []models.Reading)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//readingsMutex = &sync.RWMutex{}
	var err error
	if err = mqtt.InitMqtt(readingsChan); err != nil {
		log.Fatal(err)
	}

	if err = nats.InitNats(); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err = mqtt.WatchForPublications(ctx, "/readings"); err != nil {
			slog.Error("Watch failed", "err", err)
		}
	}()

	for batch := range readingsChan {
		readings = append(readings, batch...)
		if len(readings) > 3 {
			temps, _ := http.PredictTemps(readings)
			log.Println("Predicted:", temps)
			if err = nats.PublishTemps(temps); err != nil {
				slog.Error("Publish failed", "err", err)
			}
		}
	}

	select {
	case sig := <-sigChan:
		fmt.Printf("Signal received: %v\n", sig)
	case <-ctx.Done():
		fmt.Println("Context cancelled")
	}

}

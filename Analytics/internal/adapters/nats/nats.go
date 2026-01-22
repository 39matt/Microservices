package nats // Add package name

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

var natsConn *nats.Conn
var mu sync.Mutex

func InitNats() error {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		return err
	}
	mu.Lock()
	natsConn = nc
	mu.Unlock()
	log.Println("Connected to NATS")
	return nil
}

func PublishTemps(temps []float32) error {
	data, err := json.Marshal(temps)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if natsConn.IsClosed() {
		return fmt.Errorf("NATS closed")
	}

	if err = natsConn.Publish("predictions", data); err != nil {
		return fmt.Errorf("publish: %w", err)
	}
	log.Printf("Published %d temps: %v", len(temps), temps)
	return nil
}

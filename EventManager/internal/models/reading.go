package models

import (
	"encoding/json"
)

type Reading struct {
	ID          string
	Timestamp   string
	DeviceID    string
	Co          float64
	Humidity    float32
	Light       bool
	Lpg         float64
	Motion      bool
	Smoke       float64
	Temperature float32
}

func ConvertFromBytes(payload []byte) (*Reading, error) {
	//var ts float64
	//fmt.Sscanf(reading.GetTimestamp(), "%f", &ts)
	//timestamp := time.Unix(int64(ts), int64((ts-float64(int64(ts)))*1e9))
	reading := &Reading{}
	if err := json.Unmarshal(payload, reading); err != nil {
		return nil, err
	}

	return reading, nil
}

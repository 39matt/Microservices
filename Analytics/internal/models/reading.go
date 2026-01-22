package models

import (
	"encoding/json"
)

type Reading struct {
	ID          string  `json:"id"`
	Timestamp   string  `json:"ts"`
	DeviceID    string  `json:"device"`
	Temperature float32 `json:"temp"`
	Humidity    float32 `json:"humidity"`
	Co          float32 `json:"co"`
	Lpg         float32 `json:"lpg"`
	Smoke       float32 `json:"smoke"`
	Light       bool    `json:"light"`
	Motion      bool    `json:"motion"`
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

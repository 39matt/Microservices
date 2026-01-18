package models

import (
	"DataManager/internal/pb"
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

func ConvertFromProto(reading *pb.Reading) *Reading {
	//var ts float64
	//fmt.Sscanf(reading.GetTimestamp(), "%f", &ts)
	//timestamp := time.Unix(int64(ts), int64((ts-float64(int64(ts)))*1e9))

	return &Reading{
		ID: reading.GetId(),
		//Timestamp:   timestamp.Format(time.RFC3339),
		Timestamp:   reading.GetTimestamp(),
		DeviceID:    reading.GetDeviceId(),
		Co:          reading.GetCo(),
		Humidity:    reading.GetHumidity(),
		Light:       reading.GetLight(),
		Lpg:         reading.GetLpg(),
		Motion:      reading.GetMotion(),
		Smoke:       reading.GetSmoke(),
		Temperature: reading.GetTemperature(),
	}
}

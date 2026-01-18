package models

type Reading struct {
	ID          int64
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

//
//func (Reading) ConvertFromRaw(protoReading ReadingRaw) Reading {
//	var ts float64
//	fmt.Sscanf(protoReading.Timestamp, "%f", &ts)
//	timestamp := time.Unix(int64(ts), int64((ts-float64(int64(ts)))*1e9))
//
//	return Reading{
//		ID:          protoReading.ID,
//		Timestamp:   timestamp.Format(time.RFC3339),
//		DeviceID:    protoReading.DeviceID,
//		Co:          protoReading.Co,
//		Humidity:    protoReading.Humidity,
//		Light:       protoReading.Light == "true",
//		Lpg:         protoReading.Lpg,
//		Motion:      protoReading.Motion == "true",
//		Smoke:       protoReading.Smoke,
//		Temperature: protoReading.Temperature,
//	}
//
//}

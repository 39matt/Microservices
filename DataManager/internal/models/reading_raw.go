package models

type ReadingRaw struct {
	ID          int64
	Timestamp   string
	DeviceID    string
	Co          float64
	Humidity    float32
	Light       string
	Lpg         float64
	Motion      string
	Smoke       float64
	Temperature float32
}

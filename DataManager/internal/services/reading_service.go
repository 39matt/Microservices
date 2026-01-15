package services

import (
	"DataManager/internal/database"
	"DataManager/internal/pb"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ReadingService struct {
	DB *sql.DB
}

func NewReadingService(db *sql.DB) *ReadingService {
	return &ReadingService{DB: db}
}

func (s *ReadingService) GetAllReadings(ctx context.Context) (*pb.GetAllReadingsResponse, error) {
	var query = `SELECT id, timestamp, device_id, co, 
    			humidity, light, lpg, motion, smoke, temperature  
				FROM readings`
	rows, err := database.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*pb.Reading
	for rows.Next() {
		r := &pb.Reading{}
		var light, motion bool
		if err = rows.Scan(
			&r.Id, &r.Timestamp, &r.DeviceId, &r.Co, &r.Humidity,
			&light, &r.Lpg, &motion, &r.Smoke, &r.Temperature,
		); err != nil {
			return nil, errors.New(fmt.Sprint("Error scanning row: ", err))
		}
		r.Motion = motion
		r.Light = light
		readings = append(readings, r)
	}
	return &pb.GetAllReadingsResponse{
		Readings: readings,
	}, nil
}

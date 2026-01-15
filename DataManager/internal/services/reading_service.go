package services

import (
	"DataManager/internal/pb"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ReadingService struct {
	DB *sql.DB
}

func NewReadingService(db *sql.DB) *ReadingService {
	return &ReadingService{DB: db}
}

func scanRow(row *sql.Row) (*pb.Reading, error) {
	r := &pb.Reading{}
	var light, motion bool
	err := row.Scan(
		&r.Timestamp, &r.DeviceId, &r.Co,
		&r.Humidity, &light, &r.Lpg,
		&motion, &r.Smoke, &r.Temperature,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("reading not found")
	}
	if err != nil {
		return nil, errors.New("failed to get reading")
	}

	r.Light = light
	r.Motion = motion
	return r, nil
}

func (s *ReadingService) GetAllReadings(ctx context.Context) (*pb.GetAllReadingsResponse, error) {
	query := `SELECT id, timestamp, device_id, co, humidity, light, lpg, motion, smoke, temperature 
              FROM readings ORDER BY id DESC`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
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
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		r.Light = light
		r.Motion = motion
		readings = append(readings, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &pb.GetAllReadingsResponse{Readings: readings}, nil
}

func (s *ReadingService) GetReading(ctx context.Context, request *pb.GetReadingRequest) (*pb.Reading, error) {
	idStr := strings.TrimSpace(request.GetId())
	if idStr == "" {
		return nil, errors.New("id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	query := `SELECT id, timestamp, device_id, co, humidity, light, lpg, motion, smoke, temperature
              FROM readings WHERE id = $1`

	row := s.DB.QueryRowContext(ctx, query, id)
	return scanRow(row)
}

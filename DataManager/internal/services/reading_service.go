package services

import (
	"DataManager/internal/database"
	"DataManager/internal/pb"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ReadingService struct {
	pb.UnimplementedReadingServiceServer
	DB *sql.DB
}

func NewReadingService(db *sql.DB) *ReadingService {
	return &ReadingService{DB: db}
}

func scanRow(row *sql.Row) (*pb.Reading, error) {
	r := &pb.Reading{}
	var light, motion bool
	err := row.Scan(
		&r.Id,
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

func (s *ReadingService) GetAllReadings(ctx context.Context, empty *pb.Empty) (*pb.GetAllReadingsResponse, error) {
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

func (s *ReadingService) CreateReading(ctx context.Context, request *pb.CreateReadingRequest) (*pb.CreateReadingResponse, error) {
	r := request.GetReading()
	if r == nil {
		return nil, errors.New("reading is required")
	}

	query := `INSERT INTO readings (timestamp, device_id, co, humidity, light, lpg, motion, smoke, temperature)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id`

	var idStr string
	err := s.DB.QueryRowContext(ctx, query, r.GetTimestamp(), r.GetDeviceId(), r.GetCo(),
		r.GetHumidity(), r.GetLight(), r.GetLpg(),
		r.GetMotion(), r.GetSmoke(), r.GetTemperature()).Scan(&idStr)
	if err != nil {
		return nil, fmt.Errorf("create reading failed: %w", err)
	}
	id, err := strconv.Atoi(idStr)
	log.Printf("Added reading ID: %d", id)
	return &pb.CreateReadingResponse{Id: int32(id)}, nil
}

func (s *ReadingService) RemoveReading(ctx context.Context, request *pb.RemoveReadingRequest) (*pb.Empty, error) {
	id := request.GetId()
	if len(id) == 0 {
		return &pb.Empty{}, errors.New("id is required")
	}

	query := `DELETE FROM readings WHERE id = $1`

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return &pb.Empty{}, fmt.Errorf("remove reading failed: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, errors.New("remove reading failed: reading not found")
	}

	log.Printf("Successfully removed reading ID: %d", id)
	return &pb.Empty{}, nil
}

func (s *ReadingService) UpdateReading(ctx context.Context, request *pb.UpdateReadingRequest) (*pb.Empty, error) {
	var r *pb.Reading
	if r = request.GetReading(); r == nil {
		return nil, errors.New("reading is required")
	}

	query := `UPDATE readings 
        SET timestamp = $1, device_id = $2, co = $3, humidity = $4, light = $5, 
            lpg = $6, motion = $7, smoke = $8, temperature = $9
        WHERE id = $10`

	result, err := database.DB.ExecContext(ctx, query,
		r.GetTimestamp(), r.GetDeviceId(), r.GetCo(),
		r.GetHumidity(), r.GetLight(), r.GetLpg(),
		r.GetMotion(), r.GetSmoke(), r.GetTemperature(), r.GetId(),
	)
	if err != nil {
		return nil, fmt.Errorf("update reading failed: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, errors.New("update reading failed: reading not found")
	}
	log.Printf("Added reading ID: %d", r.GetId())
	return &pb.Empty{}, nil
}

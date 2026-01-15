package services

import "database/sql"

type ReadingService struct {
	DB *sql.DB
}

func NewReadingService(db *sql.DB) *ReadingService {
	return &ReadingService{DB: db}
}

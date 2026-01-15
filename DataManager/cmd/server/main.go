package main

import (
	"DataManager/internal/config"
	"DataManager/internal/database"
	"DataManager/internal/pb"
	"DataManager/internal/services"
	"context"
	"log"
)

func main() {

	var err error
	var ctx = context.Background()

	if err = config.Init(); err != nil {
		log.Fatal(err)
	}
	if err = database.InitPostgres(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	readingService := services.NewReadingService(database.DB)
	var readings *pb.GetAllReadingsResponse
	readings, err = readingService.GetAllReadings(ctx)
	print(readings)
}

package main

import (
	"DataManager/internal/config"
	"DataManager/internal/database"
	"DataManager/internal/pb"
	"DataManager/internal/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	var err error

	// No need in container, only locally
	if err = config.Init(); err != nil {
		log.Fatal(err)
	}
	if err = database.InitPostgres(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listening on " + lis.Addr().String())

	grpcServer := grpc.NewServer()

	readingService := services.NewReadingService(database.DB)
	pb.RegisterReadingServiceServer(grpcServer, readingService)

	grpcServer.Serve(lis)
}

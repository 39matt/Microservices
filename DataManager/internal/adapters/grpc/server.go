package grpc

import (
	"DataManager/internal/database"
	"DataManager/internal/pb"
	"DataManager/internal/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

func InitGrpcServer() error {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	log.Println("Listening on " + lis.Addr().String())

	grpcServer := grpc.NewServer()

	// Initialize services
	readingService := services.NewReadingService(database.DB)
	pb.RegisterReadingServiceServer(grpcServer, readingService)

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

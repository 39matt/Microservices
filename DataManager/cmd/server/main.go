package main

import (
	"DataManager/internal/adapters/grpc"
	"DataManager/internal/adapters/mqtt"
	"DataManager/internal/config"
	"DataManager/internal/database"
	"database/sql"
	"log"
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
	defer func(DB *sql.DB) {
		err = DB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database.DB)

	if err = mqtt.InitMqtt(); err != nil {
		log.Fatal(err)
	}

	//if err = mqtt.PublishMessage("a", "Hello"); err != nil {
	//	log.Fatal(err)
	//}

	if err = grpc.InitGrpcServer(); err != nil {
		log.Fatal(err)
	}
}

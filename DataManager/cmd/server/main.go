package main

import (
	"DataManager/internal/config"
	"DataManager/internal/database"
	"log"
)

func main() {

	var err error

	err = config.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = database.InitPostgres()
	if err != nil {
		log.Fatal(err)
		return
	}

}

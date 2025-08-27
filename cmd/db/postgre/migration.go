package main

import (
	"fmt"
	"log"

	"github.com/Eatriceeveryday/data-pool-service/internal/config"
	"github.com/Eatriceeveryday/data-pool-service/internal/database/postgredb"
	"github.com/Eatriceeveryday/data-pool-service/internal/entities"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := postgredb.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&entities.MqttSensorKey{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("DB migration done")
}

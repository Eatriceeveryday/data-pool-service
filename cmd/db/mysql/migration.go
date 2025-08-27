package main

import (
	"fmt"
	"log"

	"github.com/Eatriceeveryday/data-pool-service/internal/config"
	"github.com/Eatriceeveryday/data-pool-service/internal/database/mysqldb"
	"github.com/Eatriceeveryday/data-pool-service/internal/entities"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := mysqldb.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&entities.User{}, &entities.Sensor{}, &entities.SensorReport{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("DB migration done")
}

package postgredb

import (
	"fmt"
	"log"

	"github.com/Eatriceeveryday/data-pool-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
		cfg.PSQLDBHost, cfg.PSQLDBUsername, cfg.PSQLDBPassword, cfg.PSQLDBName, cfg.PSQLDBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	return db, nil
}

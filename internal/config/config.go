package config

import (
	"os"

	"github.com/joho/godotenv"
)

var cfg *Config

type Config struct {
	MYSQLDBUsername string
	MYSQLDBPassword string
	MYSQLDBHost     string
	MYSQLDBPort     string
	MYSQLDBName     string

	PSQLDBUsername string
	PSQLDBPassword string
	PSQLDBHost     string
	PSQLDBPort     string
	PSQLDBName     string

	AccessKey string

	EMQX_host     string
	EMQX_port     string
	EMQX_username string
	EMQX_password string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		MYSQLDBUsername: os.Getenv("MYSQL_DB_USERNAME"),
		MYSQLDBPassword: os.Getenv("MYSQL_DB_PASSWORD"),
		MYSQLDBHost:     os.Getenv("MYSQL_DB_HOST"),
		MYSQLDBPort:     os.Getenv("MYSQL_DB_PORT"),
		MYSQLDBName:     os.Getenv("MYSQL_DB_NAME"),

		PSQLDBUsername: os.Getenv("PSQL_DB_USERNAME"),
		PSQLDBPassword: os.Getenv("PSQL_DB_PASSWORD"),
		PSQLDBHost:     os.Getenv("PSQL_DB_HOST"),
		PSQLDBPort:     os.Getenv("PSQL_DB_PORT"),
		PSQLDBName:     os.Getenv("PSQL_DB_NAME"),

		AccessKey: os.Getenv("ACCESS_KEY"),

		EMQX_host:     os.Getenv("EMQX_HOST"),
		EMQX_port:     os.Getenv("EMQX_PORT"),
		EMQX_username: os.Getenv("EMQX_USERNAME"),
		EMQX_password: os.Getenv("EMQX_PASSWORD"),
	}
	return cfg, nil
}

package main

import (
	"log"
	"net/http"

	"github.com/Eatriceeveryday/data-pool-service/http/middleware"
	"github.com/Eatriceeveryday/data-pool-service/http/v1/sensor"
	"github.com/Eatriceeveryday/data-pool-service/http/v1/user"
	"github.com/Eatriceeveryday/data-pool-service/internal/config"
	"github.com/Eatriceeveryday/data-pool-service/internal/database/mysqldb"
	"github.com/Eatriceeveryday/data-pool-service/internal/database/postgredb"
	"github.com/Eatriceeveryday/data-pool-service/internal/emqx"
	"github.com/Eatriceeveryday/data-pool-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	mdb, err := mysqldb.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	pdb, err := postgredb.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := emqx.ConnectToClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	validator := validator.New()

	userService := service.NewUserService(mdb)
	userHandler := user.NewUserHandler(userService, validator)

	sensorService := service.NewSensorService(mdb, pdb)
	sensorHandler := sensor.NewSensorHandler(sensorService, validator)

	emqxService := service.NewEmqxService(client, sensorService)
	emqxService.Subscribe("sensors")

	e := echo.New()

	e.GET("/", func(ctx echo.Context) error {
		data := "Hello from /index"
		return ctx.String(http.StatusOK, data)
	})

	e.POST("/register", userHandler.CreatUser)
	e.POST("/login", userHandler.Login)

	sensorGroup := e.Group("/sensors")
	sensorGroup.Use(middleware.AuthenticateToken)

	sensorGroup.POST("", sensorHandler.CreateSensor)
	sensorGroup.GET("/reports", sensorHandler.GetSensorReportByID)
	sensorGroup.GET("/reports/duration", sensorHandler.GetSensorReportByDuration)

	e.Start(":9000")
}

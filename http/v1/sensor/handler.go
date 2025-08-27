package sensor

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Eatriceeveryday/data-pool-service/internal/entities"
	"github.com/Eatriceeveryday/data-pool-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type SensorHandler struct {
	ss *service.SensorService
	v  *validator.Validate
}

func NewSensorHandler(ss *service.SensorService, v *validator.Validate) *SensorHandler {
	return &SensorHandler{ss: ss, v: v}
}

func (h *SensorHandler) CreateSensor(c echo.Context) error {
	var req CreateSensorRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	newSensor, err := h.ss.CreateSensor(entities.Sensor{SensorType: req.SensorType, ID1: req.ID1, ID2: req.ID2, UserID: c.Get("id").(uint)})
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}
	fmt.Println("Sensor ", newSensor)

	return c.JSON(http.StatusCreated, map[string]any{"status": "Success", "data": newSensor})
}

func (h *SensorHandler) GetSensorReportByID(c echo.Context) error {
	var req GetSensorReportRequestById
	page, err := strconv.Atoi(c.QueryParam("p"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request params"})
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println("here 1")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		fmt.Println("here 2")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sensorId, err := h.ss.GetSensor(req.ID1, req.ID2, c.Get("id").(uint))
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	reports, newPage, err := h.ss.GetReportWithId(sensorId, page)

	return c.JSON(http.StatusOK, map[string]any{"status": "Success", "total": newPage, "data": reports})
}

func (h *SensorHandler) GetSensorReportByIDandDuration(c echo.Context) error {
	var req getSensorReportRequestByIDandDuration
	page, err := strconv.Atoi(c.QueryParam("p"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request params"})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sensorId, err := h.ss.GetSensor(req.ID1, req.ID2, c.Get("id").(uint))
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	reports, newPage, err := h.ss.GetReportWithDuration(sensorId, req.Start, req.End, page)

	return c.JSON(http.StatusOK, map[string]any{"status": "Success", "total": newPage, "data": reports})
}

func (h *SensorHandler) GetSensorReportByDuration(c echo.Context) error {
	var req getSensorReportRequestByDuration
	page, err := strconv.Atoi(c.QueryParam("p"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request params"})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sensors, err := h.ss.GetAllUserSensor(c.Get("id").(uint))
	var ids []uint

	for _, sensor := range sensors {
		ids = append(ids, sensor.SensorID)
	}

	reports, newPage, err := h.ss.GetReportByDuration(ids, page, req.Start, req.End)

	return c.JSON(http.StatusOK, map[string]any{"status": "Success", "total": newPage, "data": reports})
}

func (h *SensorHandler) UpdateSensorValueById(c echo.Context) error {
	var req EditSensorReportRequestById

	if err := c.Bind(&req); err != nil {
		fmt.Println("here 1")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		fmt.Println("here 2")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sensorId, err := h.ss.GetSensor(req.ID1, req.ID2, c.Get("id").(uint))
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	err = h.ss.UpdateSensorValueById(sensorId, req.Value)

	return c.JSON(http.StatusOK, map[string]any{"status": "Success"})
}

func (h *SensorHandler) UpdateSensorValueByDuration(c echo.Context) error {
	var req EditSensorReportRequestByDuration

	if err := c.Bind(&req); err != nil {
		fmt.Println("here 1")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sensors, err := h.ss.GetAllUserSensor(c.Get("id").(uint))
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	var ids []uint

	for _, sensor := range sensors {
		ids = append(ids, sensor.SensorID)
	}

	err = h.ss.UpdateSensorValueByDuration(ids, req.Start, req.End, req.Value)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "Success"})
}

func (h *SensorHandler) UpdateSensorValueByIDandDuration(c echo.Context) error {
	var req EditSensorReportRequestByIDandDuration

	if err := c.Bind(&req); err != nil {
		fmt.Println("here 1")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sensorId, err := h.ss.GetSensor(req.ID1, req.ID2, c.Get("id").(uint))
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	err = h.ss.UpdateSensorValueByIDandDuration(sensorId, req.Start, req.End, req.Value)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "Success"})
}

func (h *SensorHandler) DeleteSensorReportById(c echo.Context) error {
	var req DeleteSensorReportByIDRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sensorId, err := h.ss.GetSensor(req.ID1, req.ID2, c.Get("id").(uint))
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	err = h.ss.DeleteSensorReportByID(sensorId)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "Success"})
}

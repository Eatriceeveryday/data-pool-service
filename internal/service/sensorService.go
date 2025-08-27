package service

import (
	"fmt"
	"os"
	"time"

	"github.com/Eatriceeveryday/data-pool-service/internal/entities"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type SensorService struct {
	mdb *gorm.DB //MysqlDB main db
	pdb *gorm.DB //PostgreDB for emqx auth
}

func NewSensorService(mdb *gorm.DB, pdb *gorm.DB) *SensorService {
	return &SensorService{mdb: mdb, pdb: pdb}
}

func (s *SensorService) CreateSensor(sensor entities.Sensor) (entities.Sensor, error) {
	newSensor := entities.Sensor{
		UserID:     sensor.UserID,
		ID1:        sensor.ID1,
		ID2:        sensor.ID2,
		SensorType: sensor.SensorType,
	}

	if err := s.mdb.Create(&newSensor).Error; err != nil {
		return entities.Sensor{}, err
	}

	key, err := createApiKey(newSensor.SensorID)
	if err != nil {
		return entities.Sensor{}, fmt.Errorf("failed to generate sensor key: %w", err)
	}

	if err := s.pdb.Create(&entities.MqttSensorKey{
		SensorKey: key,
	}).Error; err != nil {
		return entities.Sensor{}, err
	}
	newSensor.SensorKey = key
	fmt.Println("New Sensor : ", newSensor)
	return newSensor, nil
}

func (s *SensorService) CreateReport(msg entities.Message) error {
	parsedTime, err := time.Parse(time.RFC3339, msg.Timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(msg.Key, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("invalid jwt: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid claims format")
	}

	sensorId := uint(claims["id"].(float64))

	newReport := entities.SensorReport{
		SensorValue: float32(msg.Value),
		Timestamp:   parsedTime,
		SensorID:    sensorId,
	}

	if err := s.mdb.Create(&newReport).Error; err != nil {
		return err
	}

	return nil
}

func (s *SensorService) GetSensor(id1 string, id2 int, userId uint) (uint, error) {
	var sensor entities.Sensor
	err := s.mdb.Where("id1 = ? AND id2 = ? AND user_id = ?", id1, id2, userId).First(&sensor).Error
	if err != nil {
		return 0, err
	}

	return sensor.SensorID, nil
}

func (s *SensorService) GetAllUserSensor(userId uint) ([]entities.Sensor, error) {
	var sensors []entities.Sensor

	err := s.mdb.Select("sensor_id").Where("user_id = ?", userId).Find(&sensors).Error
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func (s *SensorService) GetReportWithDuration(sensorId uint, start time.Time, end time.Time, page int) ([]entities.SensorReport, int64, error) {
	var reports []entities.SensorReport
	var total int64

	start = start.UTC()
	end = end.UTC()

	fmt.Println("Querying reports between", start, "and", end)
	err := s.mdb.Model(&entities.SensorReport{}).
		Where("sensor_id = ? AND timestamp BETWEEN ? AND ?", sensorId, start, end).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * 25
	err = s.mdb.Where("sensor_id = ? AND timestamp BETWEEN ? AND ?", sensorId, start, end).
		Order("timestamp ASC").
		Limit(10).
		Offset(offset).
		Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (s *SensorService) GetReportWithId(sensorId uint, page int) ([]entities.SensorReport, int64, error) {
	var reports []entities.SensorReport
	var total int64

	err := s.mdb.Model(&entities.SensorReport{}).
		Where("sensor_id = ? ", sensorId).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * 25
	err = s.mdb.Where("sensor_id = ? ", sensorId).
		Order("timestamp ASC").
		Limit(10).
		Offset(offset).
		Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (s *SensorService) GetReportByDuration(sensorId []uint, page int, start time.Time, end time.Time) ([]entities.SensorReport, int64, error) {
	var reports []entities.SensorReport
	var total int64

	start = start.UTC()
	end = end.UTC()

	err := s.mdb.Model(&entities.SensorReport{}).
		Where("sensor_id IN ? AND timestamp BETWEEN ? AND ?", sensorId, start, end).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * 25
	err = s.mdb.Where("sensor_id IN ? AND timestamp BETWEEN ? AND ?", sensorId, start, end).
		Order("timestamp ASC").
		Limit(10).
		Offset(offset).
		Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (s *SensorService) UpdateSensorValueById(sensorId uint, newValue float32) error {
	result := s.mdb.Model(&entities.SensorReport{}).Where("sensor_id = ?", sensorId).Order("timestamp DESC").Update("sensor_value", newValue)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no report found for sensor_id %d", sensorId)
	}

	return nil
}

func (s *SensorService) UpdateSensorValueByDuration(sensorId []uint, start time.Time, end time.Time, newValue float32) error {
	start = start.UTC()
	end = end.UTC()
	result := s.mdb.Model(&entities.SensorReport{}).Where("sensor_id IN ? AND timestamp BETWEEN ? AND ?", sensorId, start, end).Order("timestamp DESC").Update("sensor_value", newValue)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *SensorService) UpdateSensorValueByIDandDuration(sensorId uint, start time.Time, end time.Time, newValue float32) error {
	start = start.UTC()
	end = end.UTC()

	result := s.mdb.Model(&entities.SensorReport{}).Where("sensor_id = ? AND timestamp BETWEEN ? AND ?", sensorId, start, end).Order("timestamp DESC").Update("sensor_value", newValue)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func createApiKey(sensorId uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": sensorId,
	})

	token, err := claims.SignedString([]byte(os.Getenv("ACCESS_KEY")))
	if err != nil {
		return "", err
	}

	return token, nil
}

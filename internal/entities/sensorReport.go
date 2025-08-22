package entities

import "time"

type SensorReport struct {
	ReportID    int       `gorm:"primaryKey,autoIncrement"`
	SensorValue float32   `gorm:"not null"`
	Timestamp   time.Time `gorm:"not null"`
	SensorID    int       `gorm:"not null;index"`
}

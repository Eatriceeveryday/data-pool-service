package entities

import (
	"time"
)

type SensorReport struct {
	ReportID    uint      `gorm:"primaryKey;autoIncrement;index"`
	SensorValue float32   `gorm:"not null"`
	Timestamp   time.Time `gorm:"not null"`

	//Relation
	SensorID uint `gorm:"not null;index"`
}

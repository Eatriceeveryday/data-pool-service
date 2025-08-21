package entities

import "time"

type Sensor struct {
	ID1         string `gorm:"type:Varchar(1);not null"`
	ID2         int
	SensorType  string
	SensorValue float32
	Timestamp   time.Time
}

package entities

type Sensor struct {
	SensorID   int    `gorm:"primaryKey,autoIncrement"`
	ID1        string `gorm:"type:Varchar(1);not null;index"`
	ID2        int    `gorm:"not null;index"`
	SensorType string `gorm:"not null"`
}

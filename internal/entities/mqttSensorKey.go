package entities

type MqttSensorKey struct {
	SensorID  uint   `gorm:"primaryKey;autoIncrement;index"`
	SensorKey string `gorm:"not null"`
}

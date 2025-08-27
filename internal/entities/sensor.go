package entities

type Sensor struct {
	SensorID   uint   `gorm:"primaryKey;autoIncrement;index" json:"sensorId"`
	ID1        string `gorm:"type:Varchar(1);not null;index" json:"id1"`
	ID2        int    `gorm:"not null;index" json:"id2"`
	SensorType string `gorm:"not null" json:"sensorType"`
	SensorKey  string `gorm:"not null" json:"sensorKey"`

	//Relation
	UserID  uint           `gorm:"not null;index"  json:"userId"`
	Reports []SensorReport `gorm:"foreignKey:SensorID;references:SensorID"`
}

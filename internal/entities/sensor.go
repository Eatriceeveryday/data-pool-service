package entities

type Sensor struct {
	SensorID   uint   `gorm:"primaryKey;autoIncrement;index"`
	ID1        string `gorm:"type:Varchar(1);not null;index"`
	ID2        int    `gorm:"not null;index"`
	SensorType string `gorm:"not null"`

	//Relation
	UserID  uint           `gorm:"not null;index"`
	Reports []SensorReport `gorm:"foreignKey:SensorID;references:SensorID"`
}

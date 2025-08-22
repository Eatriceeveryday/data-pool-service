package entities

type User struct {
	UserID   uint   `gorm:"primaryKey;autoIncrement;index;unique"`
	FullName string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Password string `gorm:"not null"`

	//Relation
	Sensors []Sensor `gorm:"foreignKey:UserID;references:UserID"`
}

package entities

type User struct {
	UserID   int    `gorm:"primaryKey,autoIncrement"`
	FullName string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Password string `gorm:"not null"`
}

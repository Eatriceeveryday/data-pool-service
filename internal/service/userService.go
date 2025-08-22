package service

import (
	"github.com/Eatriceeveryday/data-pool-service/internal/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(user entities.User) (uint, error) {
	hashedPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	user = entities.User{FullName: user.FullName, Email: user.Email, Password: hashedPassword}

	result := s.db.Create(&user)

	if result.Error != nil {
		return 0, err
	}

	return user.UserID, nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

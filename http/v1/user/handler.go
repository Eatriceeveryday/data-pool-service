package user

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Eatriceeveryday/data-pool-service/internal/entities"
	"github.com/Eatriceeveryday/data-pool-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	us *service.UserService
	v  *validator.Validate
}

func NewUserHandler(us *service.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{us: us, v: v}
}

func (h *UserHandler) CreatUser(c echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if _, err := h.us.CreateUser(entities.User{FullName: req.FullName, Email: req.Email, Password: req.Password}); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email already in use"})
		}
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "Success"})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	user, err := h.us.GetUser(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid login request"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid login request"})
	}

	token, err := createToken(user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})

}

func createToken(userId uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  fmt.Sprintf("%d", userId),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("ACCESS_KEY")))
	if err != nil {
		return "", err
	}

	return token, nil
}

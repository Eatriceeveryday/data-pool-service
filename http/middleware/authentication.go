package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthenticateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization format"})
		}

		tokenStr := parts[1]

		id, err := validateToken(tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
		}

		uid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id in token"})
		}

		c.Set("id", uint(uid))

		return next(c)
	}
}

func validateToken(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_KEY")), nil
	})

	if err != nil {
		return "", err
	}

	if !tkn.Valid {
		return "", errors.New("Invalid Token")
	}

	claims := tkn.Claims.(jwt.MapClaims)
	exp := claims["exp"].(float64)
	if exp <= float64(time.Now().Unix()) {
		return "", errors.New("Expired Token")
	}

	return claims["id"].(string), nil
}

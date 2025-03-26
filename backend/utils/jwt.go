package utils

import (
	"my-rest-api/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(config.JWTSecret)
}

func SetTokenCookie(c *gin.Context, token string) {
	c.SetCookie("token", token, 3600*24, "/", "", false, true)
}

func ClearTokenCookie(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
}

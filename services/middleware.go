package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abe27/api/crypto/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashingPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func CheckPasswordHashing(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(user models.User) models.AuthSession {
	var obj models.AuthSession
	secret_key := os.Getenv("SECRET_KEY")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = obj.JwtToken
	claims["name"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenKey, _ := token.SignedString([]byte(secret_key))

	obj.Header = "Authorization"
	obj.JwtType = "Bearer"
	obj.JwtToken = tokenKey
	obj.User = &user
	return obj
}

func ValidateToken(tokenKey string) (interface{}, error) {
	token, err := jwt.Parse(tokenKey, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}
	return claims["name"], nil
}

func AuthorizationRequired(c *fiber.Ctx) error {
	var r models.Response
	s := c.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")
	if token == "" {
		r.Message = "Token required"
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	_, er := ValidateToken(token)
	if er != nil {
		r.Message = "Token is Expired"
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}
	return c.Next()
}

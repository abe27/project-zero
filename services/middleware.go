package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abe27/api/crypto/configs"
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

func CreateToken(user models.User) (models.AuthSession, error) {
	db := configs.Store
	var obj models.AuthSession
	/// Create User Login
	// fmt.Println(user.ID)
	var userLogin models.UserLogin
	db.First(&userLogin, models.UserLogin{UserID: user.ID})
	if userLogin.ID != "" {
		db.Delete(&userLogin)
	}

	userLogin.UserID = user.ID
	if err := configs.Store.Create(&userLogin).Error; err != nil {
		return obj, fmt.Errorf(err.Error())
	}

	secret_key := os.Getenv("SECRET_KEY")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = obj.JwtToken
	claims["uid"] = user.ID
	claims["pass"] = userLogin.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenKey, _ := token.SignedString([]byte(secret_key))
	obj.Header = "Authorization"
	obj.JwtType = "Bearer"
	obj.JwtToken = tokenKey
	obj.User = &user

	return obj, nil
}

func ValidateToken(tokenKey string) (interface{}, error) {
	// fmt.Println(len(tokenKey))
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

	/// Check User is Login
	var userLogin models.UserLogin
	if err := configs.Store.First(&userLogin, &models.UserLogin{ID: fmt.Sprintf("%s", claims["pass"])}).Error; err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return claims["uid"], nil
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
		r.Message = er.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}
	return c.Next()
}

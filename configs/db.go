package configs

import (
	"fmt"

	"github.com/abe27/api/crypto/models"
	"github.com/abe27/api/crypto/services"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

var (
	Store   *gorm.DB
	Session *session.Store
)

func SeedDB() {
	if !Store.Migrator().HasTable(&models.User{}) {
		Store.Migrator().CreateTable(&models.User{})

		password := services.HashingPassword("ADSads123")
		fmt.Println(len(password))
		// isMatch := services.CheckPasswordHashing("ADSads123", password)
		user := &models.User{
			Name:     "Administrator",
			Password: password,
			Email:    "krumii.it@gmail.com",
			IsAdmin:  true,
			IsActive: true,
		}

		if err := Store.Create(&user).Error; err != nil {
			panic(err)
		}
	}

	if !Store.Migrator().HasTable(&models.Logs{}) {
		Store.Migrator().CreateTable(&models.Logs{})
	}
}

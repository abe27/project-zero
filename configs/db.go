package configs

import (
	"github.com/abe27/api/crypto/models"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	Store   *gorm.DB
	Session *session.Store
)

func SeedDB() {
	if !Store.Migrator().HasTable(&models.User{}) {
		Store.Migrator().CreateTable(&models.User{})

		password, _ := bcrypt.GenerateFromPassword([]byte("ADSads123"), bcrypt.DefaultCost)
		// fmt.Println(len(password))
		// isMatch := services.CheckPasswordHashing("ADSads123", password)
		user := &models.User{
			Name:     "Administrator",
			Password: string(password),
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

	if !Store.Migrator().HasTable(&models.UserLogin{}) {
		Store.Migrator().CreateTable(&models.UserLogin{})
	}
}

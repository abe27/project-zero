package configs

import (
	"github.com/abe27/api/crypto/models"
	"gorm.io/gorm"
)

var (
	Store *gorm.DB
)

func SeedDB() {
	if !Store.Migrator().HasTable(&models.User{}) {
		Store.Migrator().CreateTable(&models.User{})
	}
}

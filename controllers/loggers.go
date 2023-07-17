package controllers

import (
	"github.com/abe27/api/crypto/configs"
	"github.com/abe27/api/crypto/models"
)

func CreateLogger(title, message string) error {
	var logs models.Logs
	logs.Title = title
	logs.Message = message
	if err := configs.Store.Create(&logs).Error; err != nil {
		return err
	}
	return nil
}

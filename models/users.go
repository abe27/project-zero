package models

import (
	"time"

	g "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:21;" json:"id"`
	Name      string    `gorm:"not null;size:255;" json:"name" binding:"required"`
	Password  string    `gorm:"size:255" json:"password" binding:"required"`
	Email     string    `gorm:"size:255;uniqueIndex" json:"email" binding:"required"`
	IsAdmin   bool      `gorm:"null" json:"is_admin" default:"false"`
	IsActive  bool      `gorm:"null" json:"is_active" default:"false"`
	CreatedAt time.Time `json:"created_at" default:"now"`
	UpdatedAt time.Time `json:"updated_at" default:"now"`
}

func (User) TableName() string {
	return "tbt_users"
}

func (obj *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := g.New()
	obj.ID = id
	return
}
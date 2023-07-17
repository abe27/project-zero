package models

import (
	"time"

	g "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:21;" json:"id,omitempty"`
	Name      string    `gorm:"not null;size:255;" json:"name" binding:"required"`
	Password  string    `gorm:"not null;size:60" json:"password" binding:"required"`
	Email     string    `gorm:"not null;size:255;uniqueIndex" json:"email" binding:"required"`
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

type AuthSession struct {
	ID       string `json:"id"`
	Header   string `json:"header"`
	JwtType  string `json:"type"`
	JwtToken string `json:"token"`
	User     *User  `json:"user"`
}

type UserLoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

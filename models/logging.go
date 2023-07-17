package models

import (
	"time"

	g "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Logs struct {
	ID        string    `json:"id,omitempty"`
	Title     string    `gorm:"not null;size:50;" json:"title,omitempty"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at" default:"now"`
}

func (Logs) TableName() string {
	return "tbt_logs"
}

func (obj *Logs) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := g.New()
	obj.ID = id
	return
}

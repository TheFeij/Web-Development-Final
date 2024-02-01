package services

import "gorm.io/gorm"

type ChatServices struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *ChatServices {
	return &ChatServices{
		DB: db,
	}
}

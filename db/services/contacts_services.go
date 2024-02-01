package services

import (
	"gorm.io/gorm"
)

type ContactsServices struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *ContactsServices {
	return &ContactsServices{
		DB: db,
	}
}

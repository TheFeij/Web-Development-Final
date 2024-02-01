package services

import "gorm.io/gorm"

type Services struct {
	ChatServices     *ChatServices
	ContactsServices *ContactsServices
	UserServices     *UserServices
}

func New(db *gorm.DB) *Services {
	return &Services{
		ChatServices: &ChatServices{
			DB: db,
		},
		ContactsServices: &ContactsServices{
			DB: db,
		},
		UserServices: &UserServices{
			DB: db,
		},
	}
}

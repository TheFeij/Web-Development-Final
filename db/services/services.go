package services

import "gorm.io/gorm"

type Services struct {
	ChatServices     *ChatServices
	ContactsServices *ContactsServices
	UserServices     *UserServices
	GroupServices    *GroupServices
	ChannelServices  *ChannelServices
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
		GroupServices: &GroupServices{
			DB: db,
		},
		ChannelServices: &ChannelServices{
			DB: db,
		},
	}
}

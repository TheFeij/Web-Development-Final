package handlers

import (
	"Messenger/db/services"
	"gorm.io/gorm"
)

type Handler struct {
	UserHandler    *UserHandler
	ContactHandler *ContactHandler
	GroupHandler   *GroupHandler
	ChannelHandler *ChannelHandler
	ChatHandler    *ChatHandler
}

func NewHandler(db *gorm.DB, services *services.Services) *Handler {
	return &Handler{
		UserHandler: &UserHandler{
			db:       db,
			services: services.UserServices,
		},
		ContactHandler: &ContactHandler{
			db:       db,
			services: services.ContactsServices,
		},
		GroupHandler: &GroupHandler{
			db:       db,
			services: services.GroupServices,
		},
		ChannelHandler: &ChannelHandler{
			db:       db,
			services: services.ChannelServices,
		},
		ChatHandler: &ChatHandler{
			db:       db,
			services: services.ChatServices,
		},
	}
}

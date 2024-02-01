package api

import (
	"Messenger/api/handlers"
	"Messenger/api/middleware"
	"Messenger/db/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	server := &Server{
		router: gin.Default(),
		db:     db,
	}

	server.router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Welcome to our shop!",
		})
	})

	handler := handlers.NewHandler(server.db, services.New(db))
	server.router.POST("/api/register", handler.UserHandler.RegisterUser)
	server.router.GET("/api/refresh-token", handler.UserHandler.GetAccessToken)
	server.router.POST("/api/login", handler.UserHandler.Login)
	usersRouter := server.router.Group("/api/users")
	{
		usersRouter.Use(middleware.AuthMiddleware())
		usersRouter.GET("/:user_id", handler.UserHandler.GetUserInformation)
		usersRouter.POST("/set_profile_image", handler.UserHandler.SetProfilePicture)
		usersRouter.PATCH("/", handler.UserHandler.UpdateUser)
		usersRouter.DELETE("/", handler.UserHandler.DeleteUser)

		usersRouter.GET("/contacts", handler.ContactHandler.GetContacts)
		usersRouter.POST("/contacts", handler.ContactHandler.AddContact)
		usersRouter.DELETE("/contacts/:contact_id", handler.ContactHandler.DeleteContact)
	}

	chatsRouter := server.router.Group("/api/chats")
	{
		chatsRouter.Use(middleware.AuthMiddleware())
		chatsRouter.POST("/", handler.ChatHandler.CreateChat)
		chatsRouter.GET("/", handler.ChatHandler.GetChats)
		chatsRouter.GET("/:chat_id", handler.ChatHandler.GetChatContent)
		chatsRouter.DELETE("/:chat_id", handler.ChatHandler.DeleteChat)
		chatsRouter.DELETE("/:chat_id/:message_id", handler.ChatHandler.DeleteChatMessage)
	}

	groupsRouter := server.router.Group("/api/groups")
	{
		groupsRouter.Use(middleware.AuthMiddleware())
		groupsRouter.POST("/", handler.GroupHandler.CreateGroup)
		groupsRouter.DELETE("/:group_id", handler.GroupHandler.DeleteGroup)
		groupsRouter.PATCH("/:group_id", handler.GroupHandler.AddMember)
		groupsRouter.DELETE("/:group_id/:user_id", handler.GroupHandler.DeleteMember)
	}

	channelsRouter := server.router.Group("/api/channels")
	{
		channelsRouter.Use(middleware.AuthMiddleware())
		channelsRouter.POST("/", handler.ChannelHandler.CreateChannel)
		channelsRouter.DELETE("/:channels_id", handler.ChannelHandler.DeleteChannel)
		channelsRouter.PATCH("/:channels_id", handler.ChannelHandler.AddMember)
		channelsRouter.DELETE("/:channels_id/:user_id", handler.ChannelHandler.DeleteMember)
	}

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

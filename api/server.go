package api

import (
	"Messenger/api/middleware"
	"Messenger/api/user"
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

	handler := user.NewHandler(server.db)
	server.router.POST("/api/register", handler.RegisterUser)
	protectedGroup := server.router.Group("/api/user")
	protectedGroup.Use(middleware.AuthMiddleware())
	server.router.GET("/api/user/:user_id", handler.GetUserInformation)

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

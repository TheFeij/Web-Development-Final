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
	server.router.GET("/api/refresh-token", handler.GetAccessToken)
	server.router.POST("/api/login", handler.Login)
	protectedRoute := server.router.Group("/api/user")
	{
		protectedRoute.Use(middleware.AuthMiddleware())
		protectedRoute.GET("/", handler.GetUserInformation)
		protectedRoute.POST("/set_profile_image", handler.SetProfilePicture)
		protectedRoute.PATCH("/", handler.UpdateUser)
		protectedRoute.DELETE("/", handler.DeleteUser)
	}

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

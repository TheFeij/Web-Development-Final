package api

import (
	"Messenger/api/middleware"
	"Messenger/api/user"
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

	service := services.New(db)
	userHandler := user.NewHandler(server.db, service.UserServices)
	server.router.POST("/api/register", userHandler.RegisterUser)
	server.router.GET("/api/refresh-token", userHandler.GetAccessToken)
	server.router.POST("/api/login", userHandler.Login)
	protectedRoute := server.router.Group("/api/user")
	{
		protectedRoute.Use(middleware.AuthMiddleware())
		protectedRoute.GET("/:user_id", userHandler.GetUserInformation)
		protectedRoute.POST("/set_profile_image", userHandler.SetProfilePicture)
		protectedRoute.PATCH("/", userHandler.UpdateUser)
		protectedRoute.DELETE("/", userHandler.DeleteUser)
	}

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

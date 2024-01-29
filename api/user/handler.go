package user

import (
	"Messenger/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h Handler) RegisterUser(context *gin.Context) {

	var req requests.RegisterUser
	if err := context.ShouldBind(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	// Handle image upload
	file, err := context.FormFile("image")
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".png" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension. Supported formats: jpg, png"})
		return
	}

	imagePath := "./data/profile_images/" + req.Username + ext
	req.Image = imagePath

	//save the user

	if err := context.SaveUploadedFile(file, imagePath); err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

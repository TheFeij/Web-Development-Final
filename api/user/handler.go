package user

import (
	"Messenger/db/services"
	"Messenger/requests"
	"Messenger/responses"
	"Messenger/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"time"
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
	service := services.New(h.db)
	var res responses.RegisterUserResponse
	res, err = service.RegisterUser(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err := context.SaveUploadedFile(file, imagePath); err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	refreshToken, err := utils.NewToken(utils.UserClaims{
		ID: res.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	accessToken, err := utils.NewToken(utils.UserClaims{
		ID: res.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.Header("Authorization", accessToken)
	context.Header("Refresh-Token", refreshToken)
	context.JSON(http.StatusOK, res)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

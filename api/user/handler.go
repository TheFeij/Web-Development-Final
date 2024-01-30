package user

import (
	"Messenger/db/services"
	"Messenger/requests"
	"Messenger/responses"
	"Messenger/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	var res responses.UserInformation
	res, err = service.RegisterUser(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err := context.SaveUploadedFile(file, imagePath); err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	refreshToken, err := utils.NewRefreshToken(res.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	accessToken, err := utils.NewAccessToken(res.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.Header("Authorization", accessToken)
	context.Header("Refresh-Token", refreshToken)
	context.JSON(http.StatusOK, res)
}

func (h Handler) GetUserInformation(context *gin.Context) {
	service := services.New(h.db)

	userID, err := strconv.ParseUint(strings.TrimSpace(context.Param("user_id")), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := service.GetUserInfo(uint(userID))
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
	}

	context.JSON(http.StatusOK, user)
}

func (h Handler) GetAccessToken(context *gin.Context) {
	refreshToken := context.GetHeader("Refresh-Token")

	if refreshToken == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}

	parsedRefreshToken := utils.ParseToken(refreshToken)
	if parsedRefreshToken == nil || parsedRefreshToken.Valid == false {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not valid"})
		return
	}

	// Attach the user claims to the context for use in the handlers
	claims, ok := parsedRefreshToken.Claims.(*utils.UserClaims)
	if !ok {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token claims"})
		return
	}

	accessToken, err := utils.NewAccessToken(claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.Header("Authorization", accessToken)
	context.Header("Refresh-Token", refreshToken)
	context.Status(http.StatusOK)
}

func (h Handler) Login(context *gin.Context) {
	var req requests.LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	service := services.New(h.db)
	res, err := service.CheckLogin(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	refreshToken, err := utils.NewRefreshToken(res.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	accessToken, err := utils.NewAccessToken(res.ID)
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

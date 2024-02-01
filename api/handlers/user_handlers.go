package handlers

import (
	"Messenger/db/services"
	"Messenger/requests"
	"Messenger/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type UserHandler struct {
	db       *gorm.DB
	services *services.UserServices
}

func NewHandler(db *gorm.DB, services *services.UserServices) *UserHandler {
	return &UserHandler{
		db:       db,
		services: services,
	}
}

func (h UserHandler) RegisterUser(context *gin.Context) {

	var req requests.RegisterUser
	if err := context.ShouldBind(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	res, err := h.services.RegisterUser(req)
	if err != nil {
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

func (h UserHandler) SetProfilePicture(context *gin.Context) {

	claims, err := GetClaims(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
	}

	file, err := context.FormFile("image")
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
	}
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".png" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension. Supported formats: jpg, png"})
		return
	}

	imagePath := "./data/profile_images/" + strconv.FormatUint(uint64(claims.ID), 10) + ext

	if err := h.services.SetProfileImage(claims.ID, imagePath); err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err := context.SaveUploadedFile(file, imagePath); err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
}

func (h UserHandler) GetUserInformation(context *gin.Context) {
	userID, err := GetUserIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
	}

	claims, err := GetClaims(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
	}

	user, err := h.services.GetUserInfo(claims.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
	}

	if uint64(claims.ID) != userID {
		if !user.DisplayPhone {
			user.Phone = ""
		}
	}
	context.JSON(http.StatusOK, user)
}

func (h UserHandler) GetAccessToken(context *gin.Context) {
	refreshToken := context.GetHeader("Refresh-Token")

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
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

func (h UserHandler) Login(context *gin.Context) {
	var req requests.LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	res, err := h.services.CheckLogin(req)
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

func (h UserHandler) DeleteUser(context *gin.Context) {
	claims, err := GetClaims(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
	}

	res, err := h.services.DeleteUser(claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func (h UserHandler) UpdateUser(context *gin.Context) {
	var req requests.RegisterUser
	if err := context.ShouldBind(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	claims, err := GetClaims(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
	}

	res, err := h.services.UpdateUser(req, claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func GetClaims(context *gin.Context) (utils.UserClaims, error) {
	var claims utils.UserClaims
	value, _ := context.Get("userClaims")
	if userClaims, ok := value.(utils.UserClaims); ok {
		claims = userClaims
	} else {
		context.Status(http.StatusInternalServerError)
		return utils.UserClaims{}, errors.New("could not get the claims")
	}

	return claims, nil
}

func GetUserIDParam(context *gin.Context) (uint64, error) {
	userID, err := strconv.ParseUint(strings.TrimSpace(context.Param("user_id")), 10, 64)
	if err != nil {
		return 0, errors.New("invalid handlers ID")
	}
	return userID, nil
}

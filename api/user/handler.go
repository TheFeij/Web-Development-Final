package user

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

	//save the user
	service := services.New(h.db)
	res, err := service.RegisterUser(req)
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

func (h Handler) SetProfilePicture(context *gin.Context) {

	claims, err := getClaims(context)
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

	service := services.New(h.db)
	if err := service.SetProfileImage(claims.ID, imagePath); err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err := context.SaveUploadedFile(file, imagePath); err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
}

func (h Handler) GetUserInformation(context *gin.Context) {

	claims, err := getClaims(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
	}

	service := services.New(h.db)
	user, err := service.GetUserInfo(claims.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
	}

	context.JSON(http.StatusOK, user)
}

func (h Handler) GetAccessToken(context *gin.Context) {
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

func (h Handler) DeleteUser(context *gin.Context) {
	claims, err := getClaims(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
	}

	service := services.New(h.db)
	res, err := service.DeleteUser(claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func (h Handler) UpdateUser(context *gin.Context) {
	var req requests.RegisterUser
	if err := context.ShouldBind(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	claims, err := getClaims(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
	}

	service := services.New(h.db)
	res, err := service.UpdateUser(req, claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func getClaims(context *gin.Context) (utils.UserClaims, error) {
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

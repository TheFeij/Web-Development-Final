package handlers

import (
	"Messenger/db/services"
	"Messenger/requests"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type ChannelHandler struct {
	db       *gorm.DB
	services *services.ChannelServices
}

func (h *ChannelHandler) CreateChannel(context *gin.Context) {
	claims, _ := GetClaims(context)

	var req requests.CreateChannel
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	newChannel, err := h.services.CreateChannel(req, claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, newChannel)
}

func (h *ChannelHandler) DeleteChannel(context *gin.Context) {
	claims, _ := GetClaims(context)

	channelID, err := getChannelIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	deletedChannel, err := h.services.DeleteChannel(claims.ID, uint(channelID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, deletedChannel)
}

func (h *ChannelHandler) AddMember(context *gin.Context) {
	claims, _ := GetClaims(context)

	channelID, err := getChannelIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var req requests.AddMember
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	addedMember, err := h.services.AddMember(req, claims.ID, uint(channelID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, addedMember)
}

func (h *ChannelHandler) DeleteMember(context *gin.Context) {
	claims, _ := GetClaims(context)

	channelID, err := getChannelIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	userID, err := getUserIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	deletedMember, err := h.services.DeleteMember(uint(userID), claims.ID, uint(channelID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, deletedMember)
}

func getChannelIDParam(context *gin.Context) (uint64, error) {
	channelID, err := strconv.ParseUint(strings.TrimSpace(context.Param("channel_id")), 10, 64)
	if err != nil {
		return 0, errors.New("invalid channel ID")
	}
	return channelID, nil
}

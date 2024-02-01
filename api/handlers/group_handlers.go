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

type GroupHandler struct {
	db       *gorm.DB
	services *services.GroupServices
}

func (h *GroupHandler) CreateGroup(context *gin.Context) {
	claims, _ := GetClaims(context)

	var req requests.CreateGroup
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	newGroup, err := h.services.CreateGroup(req, claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, newGroup)
}

func (h *GroupHandler) DeleteGroup(context *gin.Context) {
	claims, _ := GetClaims(context)

	groupID, err := getGroupIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	deletedGroup, err := h.services.DeleteGroup(claims.ID, uint(groupID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, deletedGroup)
}

func (h *GroupHandler) AddMember(context *gin.Context) {
	claims, _ := GetClaims(context)

	groupID, err := getGroupIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var req requests.AddMember
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	addedMember, err := h.services.AddMember(req, claims.ID, uint(groupID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, addedMember)
}

func (h *GroupHandler) DeleteMember(context *gin.Context) {
	claims, _ := GetClaims(context)

	groupID, err := getGroupIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	userID, err := getUserIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	deletedMember, err := h.services.DeleteMember(uint(userID), claims.ID, uint(groupID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, deletedMember)
}

func getGroupIDParam(context *gin.Context) (uint64, error) {
	groupID, err := strconv.ParseUint(strings.TrimSpace(context.Param("group_id")), 10, 64)
	if err != nil {
		return 0, errors.New("invalid group ID")
	}
	return groupID, nil
}

func getUserIDParam(context *gin.Context) (uint64, error) {
	userID, err := strconv.ParseUint(strings.TrimSpace(context.Param("user_id")), 10, 64)
	if err != nil {
		return 0, errors.New("invalid user ID")
	}
	return userID, nil
}

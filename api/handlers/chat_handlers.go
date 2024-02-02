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

type ChatHandler struct {
	db       *gorm.DB
	services *services.ChatServices
}

func (h *ChatHandler) CreateChat(context *gin.Context) {
	claims, _ := GetClaims(context)

	var req requests.CreateChat
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	newChat, err := h.services.CreateChat(req, claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, newChat)
}

func (h *ChatHandler) GetChats(context *gin.Context) {
	claims, _ := GetClaims(context)

	chatsList, err := h.services.GetChatsList(claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, chatsList)
}

func (h *ChatHandler) GetChatContent(context *gin.Context) {
	claims, _ := GetClaims(context)

	chatID, err := getChatIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	chatContent, err := h.services.GetChatContent(claims.ID, uint(chatID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, chatContent)
}

func (h *ChatHandler) DeleteChat(context *gin.Context) {
	claims, _ := GetClaims(context)

	chatID, err := getChatIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	deletedChat, err := h.services.DeleteChat(uint(chatID), claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, deletedChat)
}

func (h *ChatHandler) DeleteChatMessage(context *gin.Context) {
	claims, _ := GetClaims(context)

	chatID, err := getChatIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	messageID, err := getMessageIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	deletedMessage, err := h.services.DeleteChatMessage(claims.ID, uint(chatID), uint(messageID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, deletedMessage)
}

func getChatIDParam(context *gin.Context) (uint64, error) {
	chatID, err := strconv.ParseUint(strings.TrimSpace(context.Param("chat_id")), 10, 64)
	if err != nil {
		return 0, errors.New("invalid chat ID")
	}
	return chatID, nil
}

func getMessageIDParam(context *gin.Context) (uint64, error) {
	messageID, err := strconv.ParseUint(strings.TrimSpace(context.Param("message_id")), 10, 64)
	if err != nil {
		return 0, errors.New("invalid message ID")
	}
	return messageID, nil
}

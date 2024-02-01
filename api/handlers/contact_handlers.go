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

type ContactHandler struct {
	db       *gorm.DB
	services *services.ContactsServices
}

func (h *ContactHandler) GetContacts(context *gin.Context) {
	claims, _ := GetClaims(context)

	contacts, err := h.services.GetUserContacts(claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) AddContact(context *gin.Context) {
	claims, _ := GetClaims(context)

	var req requests.AddContact
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	newContact, err := h.services.AddContacts(req, claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, newContact)
}

func (h *ContactHandler) DeleteContact(context *gin.Context) {
	claims, _ := GetClaims(context)

	contactID, err := getContactIDParam(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	deletedContact, err := h.services.DeleteContacts(uint(contactID), claims.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	context.JSON(http.StatusOK, deletedContact)
}

func getContactIDParam(context *gin.Context) (uint64, error) {
	contactID, err := strconv.ParseUint(strings.TrimSpace(context.Param("contact_id")), 10, 64)
	if err != nil {
		return 0, errors.New("invalid contact ID")
	}
	return contactID, nil
}

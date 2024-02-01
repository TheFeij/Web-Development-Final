package services

import (
	"Messenger/db/models"
	"Messenger/responses"
	"gorm.io/gorm"
)

type ContactsServices struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *ContactsServices {
	return &ContactsServices{
		DB: db,
	}
}

func (contactsServices *ContactsServices) getUserContacts(userID uint) (responses.ContactsList, error) {
	var contactsList responses.ContactsList

	// Fetch specific attributes (ContactID, ContactName) for the given userID from the database
	if err := contactsServices.DB.Model(&models.Contact{}).
		Select("contact_id, contact_name").
		Where("user_id = ?", userID).
		Find(&contactsList.Contacts).Error; err != nil {
		return contactsList, err
	}

	return contactsList, nil
}

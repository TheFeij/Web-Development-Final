package services

import (
	"Messenger/db/models"
	"Messenger/requests"
	"Messenger/responses"
	"gorm.io/gorm"
)

type ContactsServices struct {
	DB *gorm.DB
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

func (contactsServices *ContactsServices) AddContacts(req requests.AddContact, userID uint) (responses.Contact, error) {
	newContact := models.Contact{
		UserID:      userID,
		ContactID:   req.ContactID,
		ContactName: req.ContactName,
	}

	if err := contactsServices.DB.Create(&newContact).Error; err != nil {
		return responses.Contact{}, err
	}

	return responses.Contact{
		ContactID:   newContact.ContactID,
		ContactName: newContact.ContactName,
	}, nil
}

func (contactsServices *ContactsServices) DeleteContacts(contactID uint, userID uint) (responses.Contact, error) {
	var contact models.Contact

	if err := contactsServices.DB.
		Where("user_id = ? and contact_id = ?", userID, contactID).
		First(&contact).Error; err != nil {
		return responses.Contact{}, err
	}

	if err := contactsServices.DB.Delete(&contact).Error; err != nil {
		return responses.Contact{}, err
	}

	return responses.Contact{
		ContactID:   contact.ContactID,
		ContactName: contact.ContactName,
	}, nil
}

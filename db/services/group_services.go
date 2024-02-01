package services

import (
	"Messenger/db/models"
	"Messenger/requests"
	"Messenger/responses"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupServices struct {
	DB *gorm.DB
}

func (groupServices *GroupServices) CreateGroup(req requests.CreatGroup, userID uint) (responses.Group, error) {

	newGroup := models.Groups{
		ID:    uint(uuid.New().ID()),
		Name:  req.Name,
		Owner: userID,
	}

	if err := groupServices.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newGroup).Error; err != nil {
			return err
		}

		groupParticipant := models.GroupParticipant{
			GroupID: newGroup.ID,
			UserID:  userID,
		}
		if err := tx.Create(&groupParticipant).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.Group{}, err
	}

	return responses.Group{
		ID:        newGroup.ID,
		Name:      newGroup.Name,
		CreatedAt: newGroup.CreatedAt,
		Owner:     newGroup.Owner,
	}, nil
}

func (groupServices *GroupServices) AddMember(req requests.AddMember, userID uint, groupID uint) (responses.Member, error) {
	ownerCheck := models.Groups{}
	if err := groupServices.DB.
		Where("id = ? AND owner = ?", groupID, userID).
		First(&ownerCheck).
		Error; err != nil {
		return responses.Member{}, errors.New("user is not the owner of the group")
	}

	groupParticipant := models.GroupParticipant{
		GroupID: groupID,
		UserID:  req.UserID,
	}

	if err := groupServices.DB.Create(&groupParticipant).Error; err != nil {
		return responses.Member{}, err
	}

	return responses.Member{
		UserID: groupParticipant.UserID,
	}, nil
}

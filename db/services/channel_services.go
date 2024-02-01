package services

import (
	"Messenger/db/models"
	"Messenger/requests"
	"Messenger/responses"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChannelServices struct {
	DB *gorm.DB
}

func (channelServices *ChannelServices) CreateChannel(req requests.CreateChannel, userID uint) (responses.Channel, error) {
	newChannel := models.Channel{
		ID:    uint(uuid.New().ID()),
		Name:  req.Name,
		Owner: userID,
	}

	if err := channelServices.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newChannel).Error; err != nil {
			return err
		}

		channelParticipant := models.ChannelParticipant{
			ChannelID: newChannel.ID,
			UserID:    userID,
		}
		if err := tx.Create(&channelParticipant).Error; err != nil {
			return err
		}

		channelAdmin := models.ChannelAdmin{
			ChannelID: newChannel.ID,
			UserID:    userID,
		}
		if err := tx.Create(&channelAdmin).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.Channel{}, err
	}

	return responses.Channel{
		ID:        newChannel.ID,
		Name:      newChannel.Name,
		CreatedAt: newChannel.CreatedAt,
		Owner:     newChannel.Owner,
	}, nil
}

func (channelServices *ChannelServices) AddMember(req requests.AddMember, userID uint, channelID uint) (responses.Member, error) {
	if err := channelServices.isOwner(userID, channelID); err != nil {
		return responses.Member{}, err
	}

	if err := channelServices.isMember(req.UserID, channelID); err == nil {
		return responses.Member{}, errors.New("user is already a member of the channel")
	}

	channelParticipant := models.ChannelParticipant{
		ChannelID: channelID,
		UserID:    req.UserID,
	}

	if err := channelServices.DB.Create(&channelParticipant).Error; err != nil {
		return responses.Member{}, err
	}

	return responses.Member{
		UserID: channelParticipant.UserID,
	}, nil
}

func (channelServices *ChannelServices) DeleteMember(memberID, userID, channelID uint) (responses.Member, error) {
	if err := channelServices.isOwner(userID, channelID); err != nil {
		return responses.Member{}, err
	}

	if err := channelServices.isMember(memberID, channelID); err != nil {
		return responses.Member{}, err
	}

	var deletedMember models.ChannelParticipant
	if err := channelServices.DB.Where("channel_id = ? AND user_id = ?", channelID, memberID).
		Delete(&deletedMember).
		Error; err != nil {
		return responses.Member{}, err
	}

	return responses.Member{
		UserID: deletedMember.UserID,
	}, nil
}

func (channelServices *ChannelServices) DeleteChannel(userID, channelID uint) (responses.Channel, error) {
	if err := channelServices.isOwner(userID, channelID); err != nil {
		return responses.Channel{}, err
	}

	var deletedChannel models.Channel

	if err := channelServices.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("channel_id = ?", channelID).
			Delete(models.ChannelParticipant{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where("channel_id = ?", channelID).
			Delete(models.ChannelAdmin{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where("id = ?", channelID).
			Delete(&deletedChannel).
			Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.Channel{}, err
	}

	return responses.Channel{
		ID:        deletedChannel.ID,
		Name:      deletedChannel.Name,
		CreatedAt: deletedChannel.CreatedAt,
		Owner:     deletedChannel.Owner,
	}, nil
}

func (channelServices *ChannelServices) isOwner(userID, channelID uint) error {
	var ownerCheck models.Channel
	if err := channelServices.DB.
		Where("id = ? AND owner_id = ?", channelID, userID).
		First(&ownerCheck).
		Error; err != nil {
		return errors.New("user is not the owner of the channel")
	}

	return nil
}

func (channelServices *ChannelServices) isMember(memberID, channelID uint) error {
	var memberCheck models.ChannelParticipant
	if err := channelServices.DB.
		Where("channel_id = ? AND user_id = ?", channelID, memberID).
		First(&memberCheck).
		Error; err != nil {
		return errors.New("user is not a member of the channel")
	}

	return nil
}

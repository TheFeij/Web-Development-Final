package services

import (
	"Messenger/db/models"
	"Messenger/requests"
	"Messenger/responses"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatServices struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *ChatServices {
	return &ChatServices{
		DB: db,
	}
}

func (chatServices *ChatServices) CreateChat(req requests.CreatChat, userID uint) (responses.Chat, error) {
	var newChat models.Chat

	// Check if a chat already exists between the two users
	existingChat := models.Chat{}
	err := chatServices.DB.
		Where("user_id = ? AND chat_id IN (SELECT chat_id FROM chat_participants WHERE user_id = ?)",
			userID, req.ParticipantID).
		First(&existingChat).Error
	if err == nil {
		return responses.Chat{}, errors.New("chat already exists")
	}

	if err := chatServices.DB.Transaction(func(tx *gorm.DB) error {
		// If no chat exists, create a new chat and add participants
		newChat := models.Chat{
			ID:     uint(uuid.New().ID()),
			IsDead: false,
		}

		if err := tx.Create(&newChat).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.ChatParticipant{
			ChatID: newChat.ID,
			UserID: userID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.ChatParticipant{
			ChatID: newChat.ID,
			UserID: req.ParticipantID,
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.Chat{}, err
	}

	return responses.Chat{
		ID:        newChat.ID,
		IsDead:    newChat.IsDead,
		CreatedAt: newChat.CreatedAt,
	}, nil
}

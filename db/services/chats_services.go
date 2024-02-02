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

func (chatServices *ChatServices) CreateChat(req requests.CreateChat, userID uint) (responses.Chat, error) {
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
		newChat = models.Chat{
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

func (chatServices *ChatServices) DeleteChat(chatID, userID uint) (responses.Chat, error) {
	var deletedChat models.Chat

	if err := chatServices.isUserInChat(userID, chatID); err != nil {
		return responses.Chat{}, err
	}

	if err := chatServices.DB.First(&deletedChat, chatID).Error; err != nil {
		return responses.Chat{}, err
	}

	// If the handlers is a participant, delete the chat and associated participants
	if err := chatServices.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("chat_id = ?", chatID).Delete(models.ChatParticipant{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", chatID).Delete(&deletedChat).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.Chat{}, err
	}

	return responses.Chat{
		ID:        deletedChat.ID,
		IsDead:    deletedChat.IsDead,
		CreatedAt: deletedChat.CreatedAt,
	}, nil
}

func (chatServices *ChatServices) GetChatsList(userID uint) (responses.ChatsList, error) {
	var chatsList responses.ChatsList

	if err := chatServices.DB.Model(&models.ChatParticipant{}).
		Joins("JOIN chats ON chat_participants.chat_id = chats.id").
		Where("chat_participants.user_id = ?", userID).
		Select("chats.id, chats.created_at, chats.is_dead").
		Find(&chatsList.Chats).Error; err != nil {
		return chatsList, err
	}

	return chatsList, nil
}

func (chatServices *ChatServices) GetChatContent(userID uint, chatID uint) (responses.ChatContent, error) {
	var chatContent responses.ChatContent

	if err := chatServices.isUserInChat(userID, chatID); err != nil {
		return responses.ChatContent{}, err
	}

	// If the handlers is a participant, retrieve the chat and its messages
	if err := chatServices.DB.Model(&models.Chat{}).
		Where("id = ?", chatID).
		First(&chatContent.Chat).
		Error; err != nil {
		return responses.ChatContent{}, err
	}

	// Retrieve messages for the chat
	if err := chatServices.DB.Model(&models.ChatMessage{}).
		Where("chat_id = ?", chatID).
		Order("created_at").
		Find(&chatContent.Messages).
		Error; err != nil {
		return responses.ChatContent{}, err
	}

	return chatContent, nil
}

func (chatServices *ChatServices) DeleteChatMessage(userID uint, chatID uint, messageID uint) (responses.ChatMessage, error) {
	var deletedMessage models.ChatMessage

	if err := chatServices.isUserInChat(userID, chatID); err != nil {
		return responses.ChatMessage{}, err
	}

	if err := chatServices.DB.Where("id = ? AND chat_id = ?", messageID, chatID).
		First(&deletedMessage).
		Error; err != nil {
		return responses.ChatMessage{}, err
	}

	if err := chatServices.DB.Where("id = ? AND chat_id = ?", messageID, chatID).
		Delete(&deletedMessage).
		Error; err != nil {
		return responses.ChatMessage{}, err
	}

	return responses.ChatMessage{
		ID:                deletedMessage.ID,
		ChatID:            deletedMessage.ChatID,
		SourceSenderID:    deletedMessage.SourceSenderID,
		OriginalMessageID: deletedMessage.OriginalMessageID,
		SenderID:          deletedMessage.SenderID,
		ReceiverID:        deletedMessage.ReceiverID,
		Content:           deletedMessage.Content,
		CreatedAt:         deletedMessage.CreatedAt,
	}, nil
}

func (chatServices *ChatServices) isUserInChat(userID uint, chatID uint) error {
	// Check if the handlers is a participant of the chat
	participant := models.ChatParticipant{}
	if err := chatServices.DB.
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		First(&participant).
		Error; err != nil {
		return errors.New("user is not a participant of the chat")
	}

	return nil
}

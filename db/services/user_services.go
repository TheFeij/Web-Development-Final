package services

import (
	"Messenger/db/models"
	"Messenger/requests"
	"Messenger/responses"
	"Messenger/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserServices struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *UserServices {
	return &UserServices{
		DB: db,
	}
}

func (accountServices *UserServices) RegisterUser(req requests.RegisterUser) (responses.UserInformation, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return responses.UserInformation{}, err
	}

	var newUser models.User
	newUser = models.User{
		ID:        uint(uuid.New().ID()),
		Username:  req.Username,
		Password:  hashedPassword,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Bio:       req.Bio,
		Image:     req.Image,
		Phone:     req.Phone,
	}

	if err := accountServices.DB.Create(&newUser).Error; err != nil {
		return responses.UserInformation{}, err
	}

	return responses.UserInformation{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Firstname: newUser.Firstname,
		Lastname:  newUser.Lastname,
		Bio:       newUser.Bio,
		Phone:     newUser.Phone,
	}, nil
}

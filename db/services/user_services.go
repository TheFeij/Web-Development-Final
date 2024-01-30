package services

import (
	"Messenger/db/models"
	"Messenger/requests"
	"Messenger/responses"
	"Messenger/utils"
	"errors"
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

func (userServices *UserServices) RegisterUser(req requests.RegisterUser) (responses.UserInformation, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return responses.UserInformation{}, err
	}

	var newUser models.User
	newUser = models.User{
		ID:           uint(uuid.New().ID()),
		Username:     req.Username,
		Password:     hashedPassword,
		Firstname:    req.Firstname,
		Lastname:     req.Lastname,
		Bio:          req.Bio,
		Phone:        req.Phone,
		DisplayImage: req.DisplayProfilePicture,
		DisplayPhone: req.DisplayPhone,
	}

	if err := userServices.DB.Create(&newUser).Error; err != nil {
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

func (userServices *UserServices) SetProfileImage(userID uint, imagePath string) error {
	if err := userServices.
		DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("image", imagePath).
		Error; err != nil {
		return err
	}

	return nil
}

func (userServices *UserServices) GetUserInfo(userID uint) (responses.UserInformation, error) {

	var user models.User
	if err := userServices.DB.First(&user, userID).Error; err != nil {
		return responses.UserInformation{}, err
	}

	return responses.UserInformation{
		ID:        user.ID,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Bio:       user.Bio,
		Phone:     user.Phone,
	}, nil
}

func (userServices *UserServices) CheckLogin(req requests.LoginRequest) (responses.UserInformation, error) {
	var user models.User

	if err := userServices.DB.First(&user, "username = ?", req.Username).Error; err != nil {
		return responses.UserInformation{}, errors.New("user not found")
	}

	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		return responses.UserInformation{}, errors.New("wrong Password")
	}

	return responses.UserInformation{
		ID:        user.ID,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Bio:       user.Bio,
		Phone:     user.Phone,
	}, nil
}

func (userServices *UserServices) DeleteUser(userID uint) (responses.UserInformation, error) {
	var user models.User

	if err := userServices.DB.Delete(&user, userID).Error; err != nil {
		return responses.UserInformation{}, err
	}

	return responses.UserInformation{
		ID:        user.ID,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Bio:       user.Bio,
		Phone:     user.Phone,
	}, nil
}

func (userServices *UserServices) UpdateUser(req requests.RegisterUser, userID uint) (responses.UserInformation, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return responses.UserInformation{}, err
	}

	user := models.User{
		ID:           userID,
		Username:     req.Username,
		Password:     hashedPassword,
		Firstname:    req.Firstname,
		Lastname:     req.Lastname,
		Bio:          req.Bio,
		Phone:        req.Phone,
		DisplayImage: req.DisplayProfilePicture,
		DisplayPhone: req.DisplayPhone,
	}

	if err := userServices.DB.Model(&user).Updates(&user).Error; err != nil {
		return responses.UserInformation{}, err
	}

	return responses.UserInformation{
		ID:                    user.ID,
		Username:              user.Username,
		Firstname:             user.Firstname,
		Lastname:              user.Lastname,
		Bio:                   user.Bio,
		Phone:                 user.Phone,
		DisplayPhone:          user.DisplayPhone,
		DisplayProfilePicture: user.DisplayImage,
	}, nil
}

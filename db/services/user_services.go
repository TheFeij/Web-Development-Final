package services

import (
	"Messenger/db/models"
	"Messenger/requests"
	"Messenger/responses"
	"Messenger/utils"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"unicode"
)

type UserServices struct {
	DB *gorm.DB
}

func (userServices *UserServices) RegisterUser(req requests.RegisterUser) (responses.UserInformation, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return responses.UserInformation{}, err
	}

	var newUser models.User
	newUser = models.User{
		ID:                    uint(uuid.New().ID()),
		Username:              req.Username,
		Password:              hashedPassword,
		Firstname:             req.Firstname,
		Lastname:              req.Lastname,
		Bio:                   req.Bio,
		Phone:                 req.Phone,
		DisplayProfilePicture: req.DisplayProfilePicture,
		DisplayNumber:         req.DisplayPhone,
	}

	if err := userServices.DB.Create(&newUser).Error; err != nil {
		return responses.UserInformation{}, err
	}

	return responses.UserInformation{
		ID:                    newUser.ID,
		Username:              newUser.Username,
		Firstname:             newUser.Firstname,
		Lastname:              newUser.Lastname,
		Bio:                   newUser.Bio,
		Phone:                 newUser.Phone,
		DisplayPhone:          newUser.DisplayNumber,
		DisplayProfilePicture: newUser.DisplayProfilePicture,
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

func (userServices *UserServices) GetProfileImagePath(userID uint) (string, error) {
	var user models.User

	if err := userServices.DB.Model(&models.User{}).
		Where("id = ?", userID).
		First(&user).
		Error; err != nil {
		return "", err
	}

	return user.Image, nil
}

func (userServices *UserServices) GetUserInfo(userID uint) (responses.UserInformation, error) {

	var user models.User
	if err := userServices.DB.First(&user, userID).Error; err != nil {
		return responses.UserInformation{}, err
	}

	return responses.UserInformation{
		ID:                    user.ID,
		Username:              user.Username,
		Firstname:             user.Firstname,
		Lastname:              user.Lastname,
		Bio:                   user.Bio,
		Phone:                 user.Phone,
		DisplayPhone:          user.DisplayNumber,
		DisplayProfilePicture: user.DisplayProfilePicture,
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
		ID:                    user.ID,
		Username:              user.Username,
		Firstname:             user.Firstname,
		Lastname:              user.Lastname,
		Bio:                   user.Bio,
		Phone:                 user.Phone,
		DisplayPhone:          user.DisplayNumber,
		DisplayProfilePicture: user.DisplayProfilePicture,
	}, nil
}

func (userServices *UserServices) DeleteUser(userID uint) (responses.UserInformation, error) {
	var user models.User

	if err := userServices.DB.Transaction(func(tx *gorm.DB) error {
		if err := userServices.DB.First(&user, userID).Error; err != nil {
			return err
		}

		if err := userServices.DB.Where("user_id = ? OR contact_id = ?", userID, userID).
			Delete(models.Contact{}).Error; err != nil {
			return err
		}

		if err := userServices.DB.Where("user_id = ?", userID).
			Delete(models.ChatParticipant{}).Error; err != nil {
			return err
		}

		if err := userServices.DB.Where("user_id = ?", userID).
			Delete(models.GroupParticipant{}).Error; err != nil {
			return err
		}

		if err := userServices.DB.Where("user_id = ?", userID).
			Delete(models.ChannelParticipant{}).Error; err != nil {
			return err
		}

		if err := userServices.DB.Where("user_id = ?", userID).
			Delete(models.ChannelAdmin{}).Error; err != nil {
			return err
		}

		if err := userServices.DB.Delete(&user, userID).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.UserInformation{}, err
	}

	return responses.UserInformation{
		ID:                    user.ID,
		Username:              user.Username,
		Firstname:             user.Firstname,
		Lastname:              user.Lastname,
		Bio:                   user.Bio,
		Phone:                 user.Phone,
		DisplayPhone:          user.DisplayNumber,
		DisplayProfilePicture: user.DisplayProfilePicture,
	}, nil
}

func (userServices *UserServices) UpdateUser(req requests.RegisterUser, userID uint) (responses.UserInformation, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return responses.UserInformation{}, err
	}

	user := models.User{
		ID:                    userID,
		Username:              req.Username,
		Password:              hashedPassword,
		Firstname:             req.Firstname,
		Lastname:              req.Lastname,
		Bio:                   req.Bio,
		Phone:                 req.Phone,
		DisplayProfilePicture: req.DisplayProfilePicture,
		DisplayNumber:         req.DisplayPhone,
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
		DisplayPhone:          user.DisplayNumber,
		DisplayProfilePicture: user.DisplayProfilePicture,
	}, nil
}

func (userServices *UserServices) SearchUsers(keyword string) (responses.UsersSearch, error) {
	var results responses.UsersSearch

	if !isAlphanumeric(keyword) {
		return responses.UsersSearch{}, errors.New("keyword should be alphanumeric")
	}

	if err := userServices.DB.Model(&models.User{}).
		Where("username LIKE ?", "%"+keyword+"%").
		Pluck("username", &results.Usernames).
		Error; err != nil {
		return responses.UsersSearch{}, err
	}

	return results, nil
}

func isAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}

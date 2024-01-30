package requests

type RegisterUser struct {
	Username              string `form:"username" validate:"required;min=4;max=128"`
	Password              string `form:"password" validate:"required;min=8;max=64"`
	Firstname             string `form:"firstname" validate:"required;min=1;max=50"`
	Lastname              string `form:"lastname" validate:"required;min=1;max=50"`
	Phone                 string `form:"phone" validate:"required;size=11"`
	Bio                   string `form:"bio" validate:"max=100"`
	Image                 string
	DisplayPhone          bool `form:"display_phone"`
	DisplayProfilePicture bool `form:"display_profile_picture"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required;min=4;max=128"`
	Password string `json:"password" validate:"required;min=8;max=64"`
}

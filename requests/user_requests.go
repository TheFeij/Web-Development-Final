package requests

type RegisterUser struct {
	Username              string `json:"username" validate:"required;min=4;max=128"`
	Password              string `json:"password" validate:"required;min=8;max=64"`
	Firstname             string `json:"firstname" validate:"required;min=1;max=50"`
	Lastname              string `json:"lastname" validate:"required;min=1;max=50"`
	Phone                 string `json:"phone" validate:"required;size=11"`
	Bio                   string `json:"bio" validate:"max=100"`
	DisplayPhone          bool   `json:"display_phone"`
	DisplayProfilePicture bool   `json:"display_profile_picture"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required;min=4;max=128"`
	Password string `json:"password" validate:"required;min=8;max=64"`
}

package models

type User struct {
	ID                    uint
	Firstname             string
	Lastname              string
	Phone                 string
	Username              string
	Password              string
	Image                 string
	Bio                   string
	DisplayProfilePicture bool
	DisplayNumber         bool
}

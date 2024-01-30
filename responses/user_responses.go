package responses

type UserInformation struct {
	ID                    uint   `json:"id"`
	Username              string `json:"username"`
	Firstname             string `json:"firstname"`
	Lastname              string `json:"lastname"`
	Phone                 string `json:"phone"`
	Bio                   string `json:"bio"`
	DisplayPhone          bool   `json:"display_phone"`
	DisplayProfilePicture bool   `json:"display_profile_picture"`
}

type UsersSearch struct {
	Usernames []string `json:"usernames"`
}

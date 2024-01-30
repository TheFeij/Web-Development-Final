package responses

type UserInformation struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
	Bio       string `json:"bio"`
}

type UsersSearch struct {
	Usernames []string `json:"usernames"`
}

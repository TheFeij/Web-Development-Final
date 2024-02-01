package requests

type CreatGroup struct {
	Name    string `json:"name" validate:"required;min=1;max=64"`
	OwnerID uint   `json:"owner_id" validate:"required"`
}

type AddMember struct {
	UserID uint `json:"user_id" validate:"required"`
}

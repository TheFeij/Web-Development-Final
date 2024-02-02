package requests

type CreateGroup struct {
	Name string `json:"name" validate:"required;min=1;max=64"`
}

type AddMember struct {
	UserID uint `json:"user_id" validate:"required"`
}

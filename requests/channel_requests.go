package requests

type CreateChannel struct {
	Name string `json:"name" validate:"required;min=1;max=64"`
}

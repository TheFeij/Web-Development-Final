package requests

type AddContact struct {
	ContactID   uint   `json:"contact_id" validate:"required"`
	ContactName string `json:"contact_name" validate:"required;min=1;max=100"`
}

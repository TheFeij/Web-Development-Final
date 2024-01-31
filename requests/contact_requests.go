package requests

type AddContact struct {
	ContactID   uint   `json:"contact_id"`
	ContactName string `json:"contact_name"`
}

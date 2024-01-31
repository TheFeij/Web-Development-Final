package responses

type Contact struct {
	ContactID   uint   `json:"contact_id"`
	ContactName string `json:"contact_name"`
}

type ContactsList struct {
	Contacts []Contact `json:"contacts"`
}

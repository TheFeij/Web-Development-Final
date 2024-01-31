package requests

type CreatGroup struct {
	Name    string `json:"name" validator:"required;min=1;max=64"`
	OwnerID uint   `json:"owner_id" validator:"required"`
}

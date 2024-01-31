package requests

type CreatChat struct {
	ParticipantID uint `json:"participant_id" validate:"required"`
}

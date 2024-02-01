package requests

type CreateChat struct {
	ParticipantID uint `json:"participant_id" validate:"required"`
}

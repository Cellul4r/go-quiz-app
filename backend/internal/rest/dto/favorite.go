package dto

import "github.com/google/uuid"

type FavoriteRequest struct {
	ProfileID uuid.UUID `json:"profile_id" validate:"required,uuid4"`
	QuizID    uuid.UUID `json:"quiz_id" validate:"required,uuid4"`
}

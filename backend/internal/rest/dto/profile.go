package dto

import (
	"time"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type ProfileUpdateRequest struct {
	Username  string `json:"username,omitempty" validate:"omitempty,min=1,max=20,alphanum"`
	FullName  string `json:"full_name" validate:"omitempty,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

type ProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *ProfileUpdateRequest) ToDomain(id uuid.UUID) *domain.Profile {
	return &domain.Profile{
		ID:        id,
		Username:  r.Username,
		FullName:  r.FullName,
		AvatarURL: r.AvatarURL,
	}
}

func ToProfileResponse(p *domain.Profile) *ProfileResponse {
	return &ProfileResponse{
		ID:        p.ID,
		Username:  p.Username,
		FullName:  p.FullName,
		AvatarURL: p.AvatarURL,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

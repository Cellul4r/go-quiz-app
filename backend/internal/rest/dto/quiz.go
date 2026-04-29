package dto

import (
	"time"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type QuizCreateRequest struct {
	Title            string `json:"title" validate:"required,min=1,max=100"`
	Description      string `json:"description" validate:"omitempty,max=200"`
	VisibilityStatus string `json:"visibility_status" validate:"required,oneof=public private"`
}

type QuizUpdateRequest struct {
	Title            *string `json:"title" validate:"omitempty,min=1,max=100"`
	Description      string  `json:"description" validate:"omitempty,max=200"`
	VisibilityStatus *string `json:"visibility_status" validate:"omitempty,oneof=public private"`
}

type QuizResponse struct {
	ID               uuid.UUID `json:"id"`
	ProfileID        uuid.UUID `json:"profile_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	VisibilityStatus string    `json:"visibility_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (r *QuizCreateRequest) ToDomain() *domain.Quiz {
	return &domain.Quiz{
		Title:            r.Title,
		Description:      r.Description,
		VisibilityStatus: domain.VisibilityStatus(r.VisibilityStatus),
	}
}

func (r *QuizUpdateRequest) ToDomain(quizID uuid.UUID) *domain.Quiz {
	quiz := &domain.Quiz{
		ID:          quizID,
		Description: r.Description,
	}

	if r.Title != nil {
		quiz.Title = *r.Title
	}

	if r.VisibilityStatus != nil {
		quiz.VisibilityStatus = domain.VisibilityStatus(*r.VisibilityStatus)
	}

	return quiz
}

func ToQuizResponse(q *domain.Quiz) *QuizResponse {
	return &QuizResponse{
		ID:               q.ID,
		ProfileID:        q.ProfileID,
		Title:            q.Title,
		Description:      q.Description,
		VisibilityStatus: string(q.VisibilityStatus),
		CreatedAt:        q.CreatedAt,
		UpdatedAt:        q.UpdatedAt,
	}
}

func ToQuizResponseList(quizzes *[]domain.Quiz) *[]QuizResponse {
	responses := make([]QuizResponse, len(*quizzes))
	for _, quiz := range *quizzes {
		responses = append(responses, *ToQuizResponse(&quiz))
	}
	return &responses
}

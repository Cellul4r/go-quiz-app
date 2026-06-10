package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type QuestionCreateRequest struct {
	QuestionType     string          `json:"question_type" validate:"required,oneof=multiple_choice true_false typed_answer"`
	Content          string          `json:"content" validate:"required"`
	TimeLimitSeconds int             `json:"time_limit_seconds" validate:"min=0"`
	ImageURL         string          `json:"image_url" validate:"omitempty,url"`
	SortOrder        int             `json:"sort_order" validate:"min=0"`
	Options          json.RawMessage `json:"options" validate:"required" swaggertype:"object"`
}

type QuestionUpdateRequest struct {
	QuestionType     *string          `json:"question_type" validate:"omitempty,oneof=multiple_choice true_false typed_answer"`
	Content          *string          `json:"content" validate:"omitempty"`
	TimeLimitSeconds *int             `json:"time_limit_seconds" validate:"omitempty,min=0"`
	ImageURL         *string          `json:"image_url" validate:"omitempty,url"`
	SortOrder        *int             `json:"sort_order" validate:"omitempty,min=0"`
	Options          *json.RawMessage `json:"options" swaggertype:"object"`
}

type QuestionResponse struct {
	ID               string    `json:"id"`
	QuizID           string    `json:"quiz_id"`
	QuestionType     string    `json:"question_type"`
	Content          string    `json:"content"`
	TimeLimitSeconds int       `json:"time_limit_seconds"`
	ImageURL         string    `json:"image_url"`
	SortOrder        int       `json:"sort_order"`
	Options          any       `json:"options"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type QuestionMultipleChoiceOptions struct {
	Choices        []domain.Choice `json:"choices"`
	CorrectChoices []string        `json:"correct_choices"`
}

type QuestionTrueFalseOptions struct {
	CorrectAnswer bool `json:"correct_answer"`
}

type QuestionTypedAnswerOptions struct {
	CorrectAnswer   string   `json:"correct_answer"`
	AcceptedAnswers []string `json:"accepted_answers"`
	CaseSensitive   bool     `json:"case_sensitive"`
}

func (r *QuestionCreateRequest) ToDomain(quizID uuid.UUID) (*domain.Question, error) {
	options, err := parseQuestionOptions(domain.QuestionType(r.QuestionType), r.Options)
	if err != nil {
		return nil, err
	}

	return &domain.Question{
		QuizID:           quizID.String(),
		QuestionType:     domain.QuestionType(r.QuestionType),
		Content:          r.Content,
		TimeLimitSeconds: r.TimeLimitSeconds,
		ImageURL:         r.ImageURL,
		SortOrder:        r.SortOrder,
		Options:          options,
	}, nil
}

func (r *QuestionUpdateRequest) ToDomain(questionID string) (*domain.Question, error) {
	question := &domain.Question{ID: questionID}

	if r.Content != nil {
		question.Content = *r.Content
	}
	if r.TimeLimitSeconds != nil {
		question.TimeLimitSeconds = *r.TimeLimitSeconds
	}
	if r.ImageURL != nil {
		question.ImageURL = *r.ImageURL
	}
	if r.SortOrder != nil {
		question.SortOrder = *r.SortOrder
	}
	if r.QuestionType != nil {
		question.QuestionType = domain.QuestionType(*r.QuestionType)
	}
	if r.Options != nil {
		if r.QuestionType == nil {
			return nil, fmt.Errorf("question_type is required when options is provided")
		}
		options, err := parseQuestionOptions(domain.QuestionType(*r.QuestionType), *r.Options)
		if err != nil {
			return nil, err
		}
		question.Options = options
	}

	return question, nil
}

func ToQuestionResponse(q *domain.Question) *QuestionResponse {
	return &QuestionResponse{
		ID:               q.ID,
		QuizID:           q.QuizID,
		QuestionType:     string(q.QuestionType),
		Content:          q.Content,
		TimeLimitSeconds: q.TimeLimitSeconds,
		ImageURL:         q.ImageURL,
		SortOrder:        q.SortOrder,
		Options:          q.Options,
		CreatedAt:        q.CreatedAt,
		UpdatedAt:        q.UpdatedAt,
	}
}

func ToQuestionResponseList(questions []domain.Question) []QuestionResponse {
	responses := make([]QuestionResponse, 0, len(questions))
	for _, question := range questions {
		responses = append(responses, *ToQuestionResponse(&question))
	}
	return responses
}

func parseQuestionOptions(questionType domain.QuestionType, raw json.RawMessage) (any, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("options is required")
	}

	switch questionType {
	case domain.MultipleChoice:
		var opts domain.MultipleChoiceOptions
		if err := json.Unmarshal(raw, &opts); err != nil {
			return nil, fmt.Errorf("invalid multiple choice options: %w", err)
		}
		return opts, nil
	case domain.TrueFalse:
		var opts domain.TrueFalseOptions
		if err := json.Unmarshal(raw, &opts); err != nil {
			return nil, fmt.Errorf("invalid true false options: %w", err)
		}
		return opts, nil
	case domain.TypedAnswer:
		var opts domain.TypedAnswerOptions
		if err := json.Unmarshal(raw, &opts); err != nil {
			return nil, fmt.Errorf("invalid typed answer options: %w", err)
		}
		return opts, nil
	default:
		return nil, fmt.Errorf("unknown question type: %s", questionType)
	}
}

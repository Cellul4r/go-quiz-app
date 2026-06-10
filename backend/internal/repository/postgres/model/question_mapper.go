package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"gorm.io/gorm"
)

func ToDomainQuestion(q *QuestionModel) (*domain.Question, error) {
	opts, err := parseOptions(domain.QuestionType(q.QuestionType), q.Options.raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse question options: %w", err)
	}

	var deletedAt *time.Time
	if q.DeletedAt.Valid {
		deletedAt = &q.DeletedAt.Time
	}

	return &domain.Question{
		ID:               q.ID,
		QuizID:           q.QuizID,
		QuestionType:     domain.QuestionType(q.QuestionType),
		Content:          q.Content,
		TimeLimitSeconds: q.TimeLimitSeconds,
		ImageURL:         q.ImageURL,
		SortOrder:        q.SortOrder,
		Options:          opts,
		CreatedAt:        q.CreatedAt,
		UpdatedAt:        q.UpdatedAt,
		DeletedAt:        deletedAt,
	}, nil
}

func ToModelQuestion(q *domain.Question) (*QuestionModel, error) {
	opts, err := marshalOptions(q.Options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal question options: %w", err)
	}

	var deletedAt gorm.DeletedAt
	if q.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *q.DeletedAt, Valid: true}
	}

	return &QuestionModel{
		ID:               q.ID,
		QuizID:           q.QuizID,
		QuestionType:     string(q.QuestionType),
		Content:          q.Content,
		TimeLimitSeconds: q.TimeLimitSeconds,
		ImageURL:         q.ImageURL,
		SortOrder:        q.SortOrder,
		Options:          QuestionOptionsJSON{raw: opts},
		CreatedAt:        q.CreatedAt,
		UpdatedAt:        q.UpdatedAt,
		DeletedAt:        deletedAt,
	}, nil
}

func parseOptions(qt domain.QuestionType, raw json.RawMessage) (any, error) {
	if raw == nil {
		return nil, nil
	}
	switch qt {
	case domain.MultipleChoice:
		var opts domain.MultipleChoiceOptions
		if err := json.Unmarshal(raw, &opts); err != nil {
			return nil, fmt.Errorf("failed to unmarshal multiple choice options: %w", err)
		}
		return opts, nil
	case domain.TrueFalse:
		var opts domain.TrueFalseOptions
		if err := json.Unmarshal(raw, &opts); err != nil {
			return nil, fmt.Errorf("failed to unmarshal true/false options: %w", err)
		}
		return opts, nil
	case domain.TypedAnswer:
		var opts domain.TypedAnswerOptions
		if err := json.Unmarshal(raw, &opts); err != nil {
			return nil, fmt.Errorf("failed to unmarshal typed answer options: %w", err)
		}
		return opts, nil
	default:
		return nil, fmt.Errorf("unknown question type: %s", qt)
	}
}

func marshalOptions(options any) (json.RawMessage, error) {
	if options == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal options: %w", err)
	}
	return json.RawMessage(bytes), nil
}

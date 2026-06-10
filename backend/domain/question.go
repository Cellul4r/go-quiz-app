package domain

import (
	"fmt"
	"time"
)

type QuestionType string

const (
	MultipleChoice QuestionType = "multiple_choice"
	TrueFalse      QuestionType = "true_false"
	TypedAnswer    QuestionType = "typed_answer"
)

type MultipleChoiceOptions struct {
	Choices        []Choice `json:"choices"`
	CorrectChoices []string `json:"correct_choices"`
}

type Choice struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type TrueFalseOptions struct {
	CorrectAnswer bool `json:"correct_answer"`
}

type TypedAnswerOptions struct {
	CorrectAnswer  string   `json:"correct_answer"`
	AceptedAnswers []string `json:"accepted_answers"`
	CaseSensitive  bool     `json:"case_sensitive"`
}

type Question struct {
	ID               string
	QuizID           string
	QuestionType     QuestionType
	Content          string
	TimeLimitSeconds int
	ImageURL         string
	SortOrder        int
	Options          any
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

func (q *Question) MultipleChoiceOptions() (*MultipleChoiceOptions, error) {
	opts, ok := q.Options.(MultipleChoiceOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options type for multiple choice question")
	}
	return &opts, nil
}

func (q *Question) TrueFalseOptions() (*TrueFalseOptions, error) {
	opts, ok := q.Options.(TrueFalseOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options type for true/false question")
	}
	return &opts, nil
}

func (q *Question) TypedAnswerOptions() (*TypedAnswerOptions, error) {
	opts, ok := q.Options.(TypedAnswerOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options type for typed answer question")
	}
	return &opts, nil
}

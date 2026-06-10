package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type QuestionOptionsJSON struct {
	raw json.RawMessage
}

type QuestionModel struct {
	ID               string              `gorm:"type:uuid;primaryKey"`
	QuizID           string              `gorm:"type:uuid;not null;index"`
	QuestionType     string              `gorm:"type:varchar(50);not null"`
	Content          string              `gorm:"type:text;not null"`
	TimeLimitSeconds int                 `gorm:"not null;default:0"`
	ImageURL         string              `gorm:"type:text"`
	SortOrder        int                 `gorm:"not null;default:0"`
	Options          QuestionOptionsJSON `gorm:"type:jsonb;not null"`

	CreatedAt time.Time      `gorm:"autoCreateTime;default:now()"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (QuestionModel) TableName() string {
	return "questions"
}

func (o QuestionOptionsJSON) Value() (driver.Value, error) {
	if o.raw == nil {
		return nil, nil
	}
	return []byte(o.raw), nil
}

func (o *QuestionOptionsJSON) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan QuestionOptionsJSON: expected []byte, got %T", value)
	}
	o.raw = json.RawMessage(bytes)
	return nil
}

package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VisibilityStatus string

const (
	VisibilityPublic  VisibilityStatus = "public"
	VisibilityPrivate VisibilityStatus = "private"
)

type Quiz struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	ProfileID        uuid.UUID        `json:"profile_id" gorm:"type:uuid;not null"`
	Title            string           `json:"title" gorm:"not null"`
	Description      string           `json:"description"`
	VisibilityStatus VisibilityStatus `json:"visibility_status" gorm:"not null;default:'private'"`
	CreatedAt        time.Time        `json:"created_at" gorm:"autoCreateTime;default:now()"`
	UpdatedAt        time.Time        `json:"updated_at" gorm:"autoUpdateTime;default:now()"`
	DeletedAt        gorm.DeletedAt   `json:"deleted_at" gorm:"index"`
}

type QuizFilter struct {
	ProfileID  *uuid.UUID
	OnlyPublic bool
}

func (v *VisibilityStatus) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to scan VisibilityStatus: expected string, got ", value))
	}
	*v = VisibilityStatus(str)
	return nil
}

func (v VisibilityStatus) Value() (driver.Value, error) {
	if !v.IsValid() {
		return nil, errors.New(fmt.Sprint("Invalid VisibilityStatus value: ", v))
	}
	return string(v), nil
}

func (v VisibilityStatus) IsValid() bool {
	switch v {
	case VisibilityPublic, VisibilityPrivate:
		return true
	default:
		return false
	}
}

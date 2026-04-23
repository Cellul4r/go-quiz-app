package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null" validate:"required,min=1,max=20,alphanum"`
	FullName  string         `json:"full_name" validate:"omitempty,max=100"`
	AvatarURL string         `json:"avatar_url" validate:"omitempty,url"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;default:now()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;default:now()"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

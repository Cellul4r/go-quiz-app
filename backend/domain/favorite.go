package domain

import (
	"time"

	"github.com/google/uuid"
)

type Favorite struct {
	ProfileID uuid.UUID
	QuizID    uuid.UUID
	CreatedAt time.Time
}

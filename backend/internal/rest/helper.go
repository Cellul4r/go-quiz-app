package rest

import (
	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// helpers
func getProfileID(c fiber.Ctx) uuid.UUID {
	user := c.Locals("user").(*domain.User)
	return user.ID
}

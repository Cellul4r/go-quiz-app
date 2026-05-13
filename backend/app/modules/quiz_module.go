package modules

import (
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	postgresRepo "github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest"
	"github.com/Cellul4r/go-quiz-app/backend/quiz"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterQuizModule(db *gorm.DB, appLogger *slog.Logger,
	api fiber.Router, config config.Config) {
	// Prepare repositories
	quizRepo := postgresRepo.NewQuizRepository(db, appLogger.With("layer", "repository", "component", "quiz"))

	// BUild services layer
	quizService := quiz.NewService(quizRepo, appLogger.With("layer", "service", "component", "quiz"))
	rest.NewQuizHandler(api, &config, quizService)
}

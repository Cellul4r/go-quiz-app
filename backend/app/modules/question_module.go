package modules

import (
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	postgresRepo "github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest"
	"github.com/Cellul4r/go-quiz-app/backend/question"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterQuestionModule(db *gorm.DB, appLogger *slog.Logger, api fiber.Router, config config.Config) {
	questionRepo := postgresRepo.NewQuestionRepository(db, appLogger.With("layer", "repository", "component", "question"))
	quizRepo := postgresRepo.NewQuizRepository(db, appLogger.With("layer", "repository", "component", "quiz"))

	questionService := question.NewService(questionRepo, quizRepo, appLogger.With("layer", "service", "component", "question"))
	rest.NewQuestionHandler(api, &config, questionService)
}

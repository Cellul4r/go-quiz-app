package main

// @title GoQuiz API
// @version 1.0
// @description Go Quiz - A dynamic, community-driven learning platform
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"
import (
	"log"
	"os"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/database"
	postgresRepo "github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest"
	"github.com/Cellul4r/go-quiz-app/backend/profile"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func configureFiberLog(cfg config.Config) {
	fiberlog.SetOutput(os.Stdout)

	if cfg.AppEnv == "development" || cfg.AppDebug {
		fiberlog.SetLevel(fiberlog.LevelDebug)
		return
	}

	fiberlog.SetLevel(fiberlog.LevelInfo)
}

func main() {
	// Load environment variables from .env file
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	db, dbErr := database.ConnectDatabase(config)
	if dbErr != nil {
		log.Fatal("cannot connect to database: ", dbErr)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}

	// Defer the closing of the database connection pool
	defer sqlDB.Close()

	app := fiber.New()
	configureFiberLog(config)
	if config.AppEnv == "development" {
		log.Println("Running in development mode")
		app.Use(logger.New())
	}

	api := app.Group("/api/v1")

	// Prepare repositories
	profileRepo := postgresRepo.NewProfileRepository(db)

	// BUild services layer
	profileService := profile.NewService(profileRepo)
	rest.NewProfileHandler(api, &config, profileService)
	fiberlog.Info("Server is running on port " + config.Port)
	port := config.Port
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

package main

import (
	"log"

	"github.com/Cellul4r/go-quiz-app/backend/app/modules"
	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/database"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
)

// Swagger metadata consumed by swag tooling to generate the API document.
// @title Go Quiz App API
// @version 1.0
// @description API for Go Quiz, a dynamic community-driven learning platform
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter JWT token with Bearer prefix. Example: Bearer {token}
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

	app := fiber.New(fiber.Config{
		StructValidator: domain.NewStructValidator(),
	})
	configureFiberLog(config)
	configureDevelopmentMiddleware(app, config)
	appLogger := newAppLogger(config)

	api := app.Group("/api/v1")

	// Register modules
	modules.RegisterProfileModule(db, appLogger, api, config)
	modules.RegisterQuizModule(db, appLogger, api, config)
	modules.RegisterFavoriteModule(db, appLogger, api, config)

	fiberlog.Info("Server is running on port " + config.Port)
	port := config.Port
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

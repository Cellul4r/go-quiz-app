package main

import (
	"log"

	"github.com/Cellul4r/go-quiz-app/backend/api/routes"
	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/database"
	"github.com/gofiber/fiber/v3"
)

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
	routes.Register(app)

	port := config.Port
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

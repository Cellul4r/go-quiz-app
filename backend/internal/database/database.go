package database

import (
	"fmt"
	"log"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// only log in development
	if cfg.AppEnv == "development" {
		log.Printf("Connecting to database: host=%s port=%s dbname=%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}

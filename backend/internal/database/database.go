package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var gormLogger logger.Interface
	// only log in development
	if cfg.AppEnv == "development" {
		log.Printf("Connecting to database: host=%s port=%s dbname=%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info, // log all queries
				IgnoreRecordNotFoundError: true,        // don't log ErrRecordNotFound
				Colorful:                  true,        // colored output in terminal
			},
		)
	} else {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn, // only log errors, warns in production
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err == nil {
		return db, nil
	}

	if cfg.DBHost == "host.docker.internal" && strings.Contains(err.Error(), "lookup host.docker.internal") {
		fallbackDSN := buildDSN("localhost", cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		return gorm.Open(postgres.Open(fallbackDSN), &gorm.Config{
			Logger: gormLogger,
		})
	}

	return db, err
}

func buildDSN(host, port, user, password, dbName string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, port, user, password, dbName)
}

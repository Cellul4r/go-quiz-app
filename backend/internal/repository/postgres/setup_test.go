// go:build integration
package postgres_test

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	profilePostgresRepo "github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pg_test "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	ctx := context.Background()

	container, err := pg_test.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		pg_test.WithDatabase("testdb"),
		pg_test.WithUsername("postgres"),
		pg_test.WithPassword("postgres"),
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() { container.Terminate(ctx) })

	constr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(constr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	db.AutoMigrate(&domain.Profile{})
	return db
}

func setupProfileRepository(t *testing.T) (*gorm.DB, *profilePostgresRepo.ProfileRepository) {
	t.Helper()

	db := setupTestDB(t)
	profileRepo := profilePostgresRepo.NewProfileRepository(db, nil)

	return db, profileRepo
}

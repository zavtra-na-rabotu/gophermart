package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

var ErrMigrationsFailed = errors.New("migrations failed")

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		zap.L().Info("Failed to create migration driver", zap.Error(err))
		return ErrMigrationsFailed
	}

	migration, err := migrate.NewWithDatabaseInstance("file://db/migrations", "public", driver)
	if err != nil {
		zap.L().Info("Failed to create migrate instance", zap.Error(err))
		return ErrMigrationsFailed
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		zap.L().Info("Failed to run migrations", zap.Error(err))
		return ErrMigrationsFailed
	}

	if errors.Is(err, migrate.ErrNoChange) {
		zap.L().Info("No migrations to run")
	} else {
		zap.L().Info("Successfully ran migrations")
	}

	return nil
}

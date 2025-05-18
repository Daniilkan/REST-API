package migrations

// Package migrations provides functionality to run database migrations.
// It uses the golang-migrate library to apply migrations to the PostgreSQL database.

import (
	"TestRest/internal/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// New initializes the application configuration.
// @Summary Initialize configuration
// @Description Loads the application configuration from environment variables or YAML files.
// @Tags config
// @Success 200 {object} config.Config "Application configuration"
// @Failure 500 {string} string "Failed to load configuration"
// @Router /config/new [get]
func RunMigrations(db *pgxpool.Pool, migrationsPath string, config config.Config) error {
	connString := db.Config().ConnString()
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		connString,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

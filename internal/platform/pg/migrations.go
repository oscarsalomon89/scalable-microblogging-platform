package db

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/config"
)

func RunMigrations(cfg config.Database) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Sslmode,
	)

	migrationsPath := fmt.Sprintf("file://%s/migrations", os.Getenv(config.AppPathEnv))
	m, err := migrate.New(
		migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying migrations: %w", err)
	}
	return nil
}

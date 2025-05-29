package postgres

import (
	"fmt"

	"github.com/pressly/goose/v3"
)

func (p *Postgres) RunMigrations(migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(p.db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (p *Postgres) DownMigrations(migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(p.db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run down migrations: %w", err)
	}

	return nil
}

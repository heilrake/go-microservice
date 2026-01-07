package bootstrap

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type MigratorConfig struct {
	MigrationsFS  embed.FS
	MigrationsDir string
	DatabaseURL   string
	ServiceName   string 
}

func RunMigrator(cfg MigratorConfig) error {
	source, err := iofs.New(cfg.MigrationsFS, cfg.MigrationsDir)
	if err != nil {
		return fmt.Errorf("[%s] create migration source: %w", cfg.ServiceName, err)
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		source,
		cfg.DatabaseURL,
	)
	if err != nil {
		return fmt.Errorf("[%s] create migrator: %w", cfg.ServiceName, err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[%s] run migrations: %w", cfg.ServiceName, err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("[%s] get migration version: %w", cfg.ServiceName, err)
	}

	if dirty {
		log.Printf("[%s] ⚠️ database is dirty at version %d", cfg.ServiceName, version)
	} else {
		log.Printf("[%s] ✅ migrations applied, version %d", cfg.ServiceName, version)
	}

	return nil
}

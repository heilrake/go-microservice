package bootstrap

import (
	"log"

	"ride-sharing/shared/types"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm(cfg *types.PostgresConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to postgres with GORM: %v", err)
	}

	// Get underlying *sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(int(cfg.MaxConns))
	sqlDB.SetMaxIdleConns(int(cfg.MinConns))

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}

	log.Println("Postgres connected with GORM")

	return db
}

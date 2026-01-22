package database

import (
	"fmt"
	"lmbd-digital-push-notifications/internal/shared/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPgDatabaseConnection(cfg *config.Config) (*gorm.DB, error) {

	sslMode := "disable"
	if cfg.Db.SSL {
		sslMode = "require"
	}

	fmt.Printf("Connecting to database: host=%s port=%s dbname=%s ssl=%s\n",
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.DBName,
		sslMode,
	)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.DBName,
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB instance: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	fmt.Println("✓ Database connected successfully")

	return db, nil
}

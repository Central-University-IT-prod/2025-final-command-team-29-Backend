package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/infrastructure"
	"backend/internal/infrastructure/logger"
	// "backend/internal/infrastructure/logger"
)

func NewPostgresDB(config *infrastructure.Config, log *logger.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Name,
		config.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(
			"failed to connect to database with dsn: %s; exited with error: %s", dsn, err.Error(),
		)
		return nil, err
	}
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("failed to automigrate; exited with error: %s", err.Error())
		return nil, err
	}

	return db, nil
}

package repository

import "gorm.io/gorm"

type StatRepository struct {
	db *gorm.DB
}

func NewStatRepository(db *gorm.DB) *StatRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &StatRepository{db: db}
}

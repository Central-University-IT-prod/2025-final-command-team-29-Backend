package repository

import (
	"gorm.io/gorm"
)

type PartnerRepository struct {
	db *gorm.DB
}

func NewPartnerRepository(db *gorm.DB) *PartnerRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &PartnerRepository{db: db}
}

package repository

import "gorm.io/gorm"

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &ClientRepository{db: db}
}

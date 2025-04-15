package models

import (
	"github.com/google/uuid"
)

type Client struct {
	ID           uuid.UUID `json:"id"`
	PhoneNumber  string    `json:"phone_number" validate:"required,e164"`
	Password     string    `json:"password" gorm:"-"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	Name         string    `json:"name" validate:"required,min=3,max=50"`
}

type Partner struct {
	ID           uuid.UUID `json:"id"`
	PhoneNumber  string    `json:"phone_number" validate:"required,e164"`
	Password     string    `json:"password" gorm:"-"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	Title        string    `json:"title" validate:"required,min=3,max=50"`
	Description  string    `json:"description" validate:"required,min=3,max=500"`
	ImageID      int       `json:"image_id" validate:"required"`
}

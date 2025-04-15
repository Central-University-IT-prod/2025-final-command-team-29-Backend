package repository

import (
	"gorm.io/gorm"

	"backend/internal/models"
)

func (r *ClientRepository) CreateClient(client *models.Client) error {
	return r.db.Table("client").Create(client).Error
}

func (r *ClientRepository) GetClientByPhoneNumber(phoneNumber string) (*models.Client, error) {
	var client models.Client
	err := r.db.Table("client").
		Where("phone_number = ? OR phone_number = ?", phoneNumber, phoneNumber[1:]).
		First(&client).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &client, nil
}

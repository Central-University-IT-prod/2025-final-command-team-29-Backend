package repository

import (
	"gorm.io/gorm"

	"backend/internal/models"
)

func (r *PartnerRepository) CreatePartner(partner *models.Partner) error {
	return r.db.Table("partner").Create(partner).Error
}

func (r *PartnerRepository) GetPartnerByPhoneNumber(phoneNumber string) (*models.Partner, error) {
	var partner models.Partner
	err := r.db.Table("partner").
		Where("phone_number = ? OR phone_number = ?", phoneNumber, phoneNumber[1:]).
		First(&partner).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &partner, nil
}

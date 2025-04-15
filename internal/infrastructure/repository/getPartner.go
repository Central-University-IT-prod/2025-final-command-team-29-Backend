package repository

import "backend/internal/models"

func (r *PartnerRepository) GetPartnersList(page, limit int) ([]models.Partner, error) {
	var partners []models.Partner
	offset := (page - 1) * limit
	err := r.db.Table("partner").Limit(limit).Offset(offset).Find(&partners).Error
	return partners, err
}

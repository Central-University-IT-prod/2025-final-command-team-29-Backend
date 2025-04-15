package usecase

import (
	"backend/internal/infrastructure/repository"
	"backend/internal/models"
)

type PartnerService struct {
	partnerRepo *repository.PartnerRepository
}

func NewPartnerService(repo *repository.PartnerRepository) *PartnerService {
	return &PartnerService{partnerRepo: repo}
}

func (s *PartnerService) GetPartnersList(page, limit int) ([]models.Partner, error) {
	return s.partnerRepo.GetPartnersList(page, limit)
}

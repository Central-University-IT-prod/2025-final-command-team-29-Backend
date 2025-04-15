package usecase

import (
	"backend/internal/infrastructure/repository"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ActivationService struct {
	activationRepo *repository.ActivationRepository
	promoRepo      *repository.PromoRepository
}

func NewActivationService(activationRepo *repository.ActivationRepository, promoRepo *repository.PromoRepository) *ActivationService {
	return &ActivationService{activationRepo: activationRepo, promoRepo: promoRepo}
}

func (s *ActivationService) Activate(p *models.ActivationRequest, c echo.Context, partnerID uuid.UUID) error {
	if err := s.activationRepo.CheckGlobalConditions(p); err != nil {
		return c.JSON(406, utils.Err{Message: err.Error()})
	}

	promo, err := s.promoRepo.GetConditionByID(p.PromoID, partnerID)
	if err != nil {
		return c.JSON(500, utils.Err{Message: err.Error()})
	}
	typeOfCondition := promo.Condition.Type
	var activation *models.Activation
	switch typeOfCondition {
	case "COMMON":
		if err := s.activationRepo.ActivateCommon(p); err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(204, nil)
	case "UNIQUE":
		activation, err = s.activationRepo.CheckUniqueConditions(p)
		if err != nil {
			return c.JSON(406, utils.Err{Message: err.Error()})
		}
		if err := s.activationRepo.Create(activation); err != nil {
			return c.JSON(500, err)
		}
		return c.JSON(204, nil)
	case "BUNDLE":
		activation, err = s.activationRepo.GetActivation(p)
		if err != nil {
			return c.JSON(500, utils.Err{Message: err.Error()})
		}
		if activation.Buy == nil {
			activation.Buy = promo.Condition.Buy
		}
		s.activationRepo.ActivateBundle(p, activation)
	}
	_ = s.activationRepo.Activate(activation)
	return c.JSON(204, nil)
}

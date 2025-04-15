package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"backend/internal/models"
	"backend/pkg/hash"
	tokenjwt "backend/pkg/token_jwt"
)

func (s *PartnerService) SignUpPartner(u *models.Partner, c echo.Context) (string, error) {
	pswHsh, err := hash.GenerateHash(u.Password)
	if err != nil {
		return "", err
	}
	u.PasswordHash = pswHsh
	u.ID = uuid.New()
	token, _ := tokenjwt.GenerateJWT(u.ID)

	return token, s.partnerRepo.CreatePartner(u)
}

func (s *PartnerService) SignInPartner(u *models.Partner, c echo.Context) (string, error) {
	partner, err := s.partnerRepo.GetPartnerByPhoneNumber(u.PhoneNumber)
	if err != nil {
		return "", err
	}
	if partner == nil {
		return "", errors.New("Not Found")
	}
	if err := hash.ComparePassword(u.Password, partner.PasswordHash); err != nil {
		return "", err
	}

	token, _ := tokenjwt.GenerateJWT(partner.ID)
	*u = *partner
	return token, nil
}

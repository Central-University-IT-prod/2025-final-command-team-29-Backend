package usecase

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"backend/internal/models"
	"backend/pkg/hash"
)

func (s *ClientService) SignUp(u *models.Client, c echo.Context) error {
	pswHsh, err := hash.GenerateHash(u.Password)
	if err != nil {
		return err
	}
	u.PasswordHash = pswHsh
	u.ID = uuid.New()

	return s.clientRepo.CreateClient(u)
}

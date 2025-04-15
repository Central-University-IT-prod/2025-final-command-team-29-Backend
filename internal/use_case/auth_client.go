package usecase

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"backend/internal/models"
	"backend/pkg/hash"
	tokenjwt "backend/pkg/token_jwt"
)

func (s *ClientService) SignUpClient(u *models.Client, c echo.Context) (string, error) {
	pswHsh, err := hash.GenerateHash(u.Password)
	if err != nil {
		return "", err
	}
	u.PasswordHash = pswHsh
	u.ID = uuid.New()
	token, _ := tokenjwt.GenerateJWT(u.ID)

	return token, s.clientRepo.CreateClient(u)
}

func (s *ClientService) SignInClient(u *models.Client, c echo.Context) (string, error) {
	client, err := s.clientRepo.GetClientByPhoneNumber(u.PhoneNumber)
	if err != nil {
		return "", err
	}
	if client == nil {
		return "", errors.New("Not Found")
	}
	if err := hash.ComparePassword(u.Password, client.PasswordHash); err != nil {
		return "", err
	}

	token, _ := tokenjwt.GenerateJWT(client.ID)

	*u = *client

	fmt.Println(token)
	return token, nil
}

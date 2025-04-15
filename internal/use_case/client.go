package usecase

import "backend/internal/infrastructure/repository"

type ClientService struct {
	clientRepo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{clientRepo: repo}
}

package handlers

import (
	"backend/internal/models"
)

type clientSignResponse struct {
	models.Client
	Token string `json:"token"`
}

type partnerSignResponse struct {
	models.Partner
	Token string `json:"token"`
}

func newClientSignResponse(client *models.Client, token string) *clientSignResponse {
	return &clientSignResponse{
		Client: *client,
		Token:  token,
	}
}

func newPartnerSignResponse(partner *models.Partner, token string) *partnerSignResponse {
	return &partnerSignResponse{
		Partner: *partner,
		Token:   token,
	}
}

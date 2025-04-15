package models

import "github.com/google/uuid"

type ModerationBody struct {
	PromoID   uuid.UUID `json:"promo_id"`
	PartnerID uuid.UUID `json:"partner_id"`
	Verdict   bool      `json:"verdict"`
}

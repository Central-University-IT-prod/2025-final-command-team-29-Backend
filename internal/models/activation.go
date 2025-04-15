package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivationRequest struct {
	ClientID uuid.UUID `json:"client_id"`
	PromoID  uuid.UUID `json:"promo_id"`
}

type Activation struct {
	ClientID  uuid.UUID `json:"-"`
	PromoID   uuid.UUID `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Buy       *int      `json:"buy"`
	Current   *int      `json:"current"`
	Overall   int       `json:"overall"`
}

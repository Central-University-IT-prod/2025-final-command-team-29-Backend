package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Promo struct {
	ID              uuid.UUID `json:"promo_id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title           string    `json:"title"`
	Description     *string   `json:"description,omitempty"`
	ConditionString *string   `json:"condition_string,omitempty"` // ✅ ИСПРАВЛЕНО ЗДЕСЬ!
	Condition       Condition `json:"condition" gorm:"type:jsonb;serializer:json" validate:"required"`
	UsageLimit      *uint     `json:"usage_limit,omitempty"`
	TimeLimitStart  *int64    `json:"time_limit_start,omitempty"`
	TimeLimitEnd    *int64    `json:"time_limit_end,omitempty"`
	PartnerID       uuid.UUID `json:"partner_id"`
	IsActive        bool      `json:"is_active"`
	IsApproved      *bool     `json:"is_approved"`
	ImageID         int       `json:"image_id" validate:"required"`
}

type Condition struct {
	Type string `json:"type" validate:"oneof=COMMON UNIQUE BUNDLE"`
	Buy  *int   `json:"buy"`
}

type ApprovePromoResponce struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a Condition) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Condition) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type PromoCLientView struct {
	Promo
	*Activation
}

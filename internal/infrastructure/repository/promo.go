package repository

import (
	"backend/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) *PromoRepository {
	return &PromoRepository{
		db: db,
	}
}

func (p *PromoRepository) Create(promo *models.Promo) error {
	promo.ID = uuid.New()
	isActive := true
	if promo.TimeLimitStart != nil && promo.TimeLimitEnd != nil {
		isActive = *promo.TimeLimitStart < time.Now().Unix() && *promo.TimeLimitEnd > time.Now().Unix()
	} else if promo.TimeLimitEnd != nil {
		isActive = *promo.TimeLimitEnd > time.Now().Unix()
	} else if promo.TimeLimitStart != nil {
		isActive = *promo.TimeLimitStart < time.Now().Unix()
	}
	promo.IsActive = isActive
	return p.db.Model(&models.Promo{}).Create(promo).Error
}

func (p *PromoRepository) GetByIDForClient(promoID, partnerID, client_id uuid.UUID) (*models.PromoCLientView, error) {
	promo, err := p.GetByID(promoID, partnerID)
	if err != nil {
		return nil, err
	}
	if promo == nil {
		return nil, nil
	}
	activation := new(models.Activation)
	if err := p.db.Model(&models.Activation{}).First(activation, "client_id = ? AND promo_id = ?", client_id, promoID).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	resp := &models.PromoCLientView{
		Promo:      *promo,
		Activation: activation,
	}
	return resp, nil
}

func (p *PromoRepository) GetByID(promoID, partnerID uuid.UUID) (*models.Promo, error) {
	promo := new(models.Promo)
	if err := p.db.Model(&models.Promo{}).First(promo, "id = ? AND partner_id = ?", promoID, partnerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return promo, nil
}

func (p *PromoRepository) GetByPartnerID(id uuid.UUID) ([]*models.Promo, error) {
	promos := new([]*models.Promo)
	status := true
	pr := &models.Promo{PartnerID: id, IsApproved: &status}
	if err := p.db.Model(&models.Promo{}).Where(pr).Find(promos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return *promos, nil
}

func (p *PromoRepository) GetByPartnerIDForClient(clientID, partnerID uuid.UUID) ([]*models.PromoCLientView, error) {
	var promos []*models.Promo
	status := true
	if err := p.
		db.Model(&models.Promo{}).Where(&models.Promo{PartnerID: partnerID, IsApproved: &status}).Find(&promos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	promoIDs := make([]uuid.UUID, len(promos))
	for k, promo := range promos {
		promoIDs[k] = promo.ID
	}
	var activations []*models.Activation
	fmt.Print(clientID)
	if err := p.db.Model(&models.Activation{}).Where(&models.Activation{ClientID: clientID}).Find(&activations, "promo_id IN ? AND client_id = ?", promoIDs, clientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			return nil, err
		}
	}
	fmt.Print(activations)
	mapActivasioon := make(map[uuid.UUID]*models.Activation)
	for _, activation := range activations {
		mapActivasioon[activation.PromoID] = activation
	}
	response := make([]*models.PromoCLientView, len(promos))
	for k, promo := range promos {
		response[k] = &models.PromoCLientView{
			Promo: *promo,
		}
		activ, ok := mapActivasioon[promo.ID]
		if !ok {
			response[k].Activation = nil
		} else {
			response[k].Activation = activ
		}
	}

	return response, nil
}

func (p *PromoRepository) DeleteByID(promoID, partnerID uuid.UUID) error {
	if err := p.db.Delete(&models.Promo{}, "id = ? AND partner_id = ?", promoID, partnerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (p *PromoRepository) GetConditionByID(promoID, partnerID uuid.UUID) (*models.Promo, error) {
	promo := new(models.Promo)
	if err := p.db.Model(&models.Promo{}).First(promo, "id = ? AND partner_id = ?", promoID, partnerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	val, _ := promo.Condition.Value()
	if err := promo.Condition.Scan(val); err != nil {
		return nil, fmt.Errorf("failed to scan condition: %w", err)
	}
	return promo, nil
}

func (p *PromoRepository) GetNotApprovedPromo() (*models.Promo, error) {
	promo := new(models.Promo)
	if err := p.db.Raw(`SELECT * FROM "promos" WHERE is_approved IS NULL ORDER BY "promos"."id" LIMIT 1`).Scan(promo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return promo, nil
}

func (p *PromoRepository) ModeratePromo(promoID, partnerID uuid.UUID, verdict bool) error {
	return p.db.Model(&models.Promo{}).Where("id = ? AND partner_id = ?", promoID, partnerID).Update("is_approved", verdict).Error
}

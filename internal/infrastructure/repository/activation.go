package repository

import (
	"backend/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrGlobalConditionFailed = errors.New("activation failed on global conditions")
	ErrUniqueConditionFailed = errors.New("activation failed on bundle conditions")
	ErrBundleConditionFailed = errors.New("activation failed on unique conditions")
)

type ActivationRepository struct {
	db *gorm.DB
}

func NewActivationRepository(db *gorm.DB) *ActivationRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &ActivationRepository{db: db}
}

func (r *ActivationRepository) GetActivation(p *models.ActivationRequest) (*models.Activation, error) {
	activation := new(models.Activation)
	err := r.db.FirstOrCreate(activation, models.Activation{ClientID: p.ClientID, PromoID: p.PromoID}, models.Activation{ClientID: p.ClientID, PromoID: p.PromoID}).Error
	return activation, err
}

func (r *ActivationRepository) CheckGlobalConditions(p *models.ActivationRequest) error {
	var promo models.Promo
	err := r.db.First(&promo, p.PromoID).Error
	if err != nil {
		return err
	}
	if promo.TimeLimitStart != nil && time.Now().Unix() < *promo.TimeLimitStart {
		return ErrGlobalConditionFailed
	}
	if promo.TimeLimitEnd != nil && time.Now().Unix() > *promo.TimeLimitEnd {
		return ErrGlobalConditionFailed
	}
	var usages int64
	r.db.Model(&models.Activation{}).Where("promo_id = ?", p.PromoID).Count(&usages)
	if promo.UsageLimit != nil && usages >= int64(*promo.UsageLimit) {
		return ErrGlobalConditionFailed
	}
	return nil
}

func (r *ActivationRepository) CheckUniqueConditions(p *models.ActivationRequest) (*models.Activation, error) {
	activation := new(models.Activation)
	err := r.db.First(activation, "client_id = ? AND promo_id = ?", p.ClientID, p.PromoID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUniqueConditionFailed
	}
	activation = &models.Activation{
		ClientID: p.ClientID,
		PromoID:  p.PromoID,
	}
	return activation, nil
}

func (r *ActivationRepository) ActivateCommon(p *models.ActivationRequest) error {
	return r.db.Model(&models.Activation{}).FirstOrCreate(nil, &models.Activation{PromoID: p.PromoID, ClientID: p.ClientID}).Where(&models.Activation{PromoID: p.PromoID, ClientID: p.ClientID}).UpdateColumn("overall", gorm.Expr("overall + ?", 1)).Error
}

func (r *ActivationRepository) ActivateBundle(p *models.ActivationRequest, a *models.Activation) {
	if a.Current != nil && a.Buy != nil {
		if *a.Current+1 == *a.Buy {
			a.Overall++
			zero := 0
			a.Current = &zero
		} else {
			value := *a.Current + 1
			a.Current = &value
		}
	}
	if a.Current == nil {
		value := 1
		a.Current = &value
	}
}

func (r *ActivationRepository) Activate(a *models.Activation) error {
	a.Overall += 1
	return r.db.Model(&models.Activation{}).Where("promo_id = ? AND client_id = ?", a.PromoID, a.ClientID).Save(a).Error
}

func (r *ActivationRepository) Create(a *models.Activation) error {
	return r.db.Create(a).Error
}

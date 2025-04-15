package handlers

import (
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Moderation godoc
//
//	@Summary		Get not approved endpoint
//	@Description	Get list for moderation
//	@Tags			Moderation
//	@Produce		json
//	@Success		204
//	@Success		200	{object}	models.Promo
//	@Success		400	{object}	utils.Err
//	@Failure		500	{object}	utils.Err
//	@Router			/moderation [get]
func (p *PromoHandler) GetNotApprovedPromo(c echo.Context) error {
	// var body models.ModerationBody
	// if err := c.Bind(body); err != nil {
	// 	return c.JSON(400, utils.Err{Message: err.Error()})
	// }
	promo, err := p.promoService.GetNotApprovedPromo()
	if promo == nil || promo.ID == uuid.Nil {
		return c.JSON(204, nil)
	}
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, promo)
}

// Moderation godoc
//
//	@Summary		Moderate promo
//	@Description	Change moderation status for promo
//	@Tags			Moderation
//	@Consume		json
//	@Produce		json
//	@Param			Verdict	body	models.ModerationBody	true	"Verdict of moderator"
//	@Success		204
//	@Success		400	{object}	utils.Err
//	@Failure		500	{object}	utils.Err
//	@Router			/moderation [post]
func (p *PromoHandler) ModeratePromo(c echo.Context) error {
	var body models.ModerationBody
	if err := c.Bind(&body); err != nil {
		c.Logger().Error(err)
		return c.JSON(400, utils.Err{Message: err.Error()})
	}
	if err := p.promoService.ModeratePromo(body.PromoID, body.PartnerID, body.Verdict); err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(204, nil)
}

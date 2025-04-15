package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"backend/internal/models"
	usecase "backend/internal/use_case"
	"backend/internal/utils"
)

type ActivationHandler struct {
	service  *usecase.ActivationService
	validate *validator.Validate
}

func NewActivationHandler(s *usecase.ActivationService) *ActivationHandler {
	return &ActivationHandler{
		service:  s,
		validate: validator.New(),
	}
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Activate promo for client
//	@Tags			PromosForPartners
//	@Produce		json
//	@Param			promo			body		models.ActivationRequest	true	"Activation Request"
//	@Param			Authorization	header		string						true	"Authorization header must be set for access"
//	@Success		201				{object}	models.Promo
//	@Failure		400				{object}	utils.Err
//	@Failure		500				{object}	utils.Err
//	@Router			/promos/activate [post]
func (h *ActivationHandler) Activate(c echo.Context) error {

	id, err := getIDFromToken(c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.Err{Message: "Unauthorized"})
	}
	fmt.Print(id)
	var body models.ActivationRequest

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return h.service.Activate(&body, c, id)
}

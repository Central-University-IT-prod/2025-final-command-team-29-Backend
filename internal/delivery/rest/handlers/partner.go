package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	usecase "backend/internal/use_case"
	"backend/internal/utils"
)

type PartnerHandler struct {
	partnerService *usecase.PartnerService
	validate       *validator.Validate
}

type query struct {
	Page *int `query:"page"`
}

func NewPartnerHandler(partnerService *usecase.PartnerService) *PartnerHandler {
	return &PartnerHandler{
		partnerService: partnerService,
		validate:       validator.New(),
	}
}

// Partners godoc
//
//	@Summary		Partners endpoint
//	@Description	Returning list of partners for client
//	@Tags			Partners
//	@Produce		json
//	@Param			page			query		int		false	"Page pagination"
//	@Param			Authorization	header		string	true	"Authorization header must be set for access"
//	@Success		200				{object}	[]models.Partner
//	@Failure		400				{object}	utils.Err
//	@Failure		500				{object}	utils.Err
//	@Router			/partners [get]
func (h *PartnerHandler) GetPartnersList(c echo.Context) error {
	limit := 10
	q := new(query)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Err{Message: err.Error()})
	}
	pageInt := 1
	if q.Page != nil {
		pageInt = *q.Page
	}
	partners, err := h.partnerService.GetPartnersList(pageInt, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, partners)
}

package handlers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	usecase "backend/internal/use_case"
)

type AuthHandler struct {
	clientService  *usecase.ClientService
	partnerService *usecase.PartnerService
	validate       *validator.Validate
}

func NewAuthHandler(clientService *usecase.ClientService, partnerService *usecase.PartnerService) *AuthHandler {
	return &AuthHandler{
		clientService:  clientService,
		partnerService: partnerService,
		validate:       validator.New(),
	}
}

// ClientAuth godoc
//
//	@Summary		Create client
//	@Description	create client with given data
//	@Tags			Clients
//	@Accept			json
//	@Produce		json
//	@Param			client	body		models.Client	true	"Credentials to use"
//	@Success		201		{object}	handlers.clientSignResponse
//	@Failure		400
//	@Failure		500
//	@Router			/clients/auth/sign-up [post]
func (h *AuthHandler) SignUpClient(c echo.Context) error {
	var client models.Client
	if err := c.Bind(&client); err != nil {
		return c.JSON(utils.BadRequestError())
	}

	if err := h.validate.Struct(client); err != nil {
		return c.JSON(utils.BadRequestError())
	}

	token, err := h.clientService.SignUpClient(&client, c)
	if err != nil {
		return c.JSON(utils.BadRequestError())
	}

	return c.JSON(http.StatusCreated, newClientSignResponse(&client, token))
}

// ClientAuth godoc
//
//	@Summary		Sign-in for client
//	@Description	sign-in in client with given data
//	@Tags			Clients
//	@Accept			json
//	@Produce		json
//	@Param			client	body		handlers.signInRequest	true	"Credentials to use"
//	@Success		200		{object}	handlers.clientSignResponse
//	@Failure		400		{object}	utils.Err
//	@Failure		500		{object}	utils.Err
//	@Router			/clients/auth/sign-in [post]
func (h *AuthHandler) SignInClient(c echo.Context) error {
	var client models.Client
	req := new(signInRequest)
	if err := req.validate(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Err{Message: err.Error()})
	}

	client.Password = req.Password
	client.PhoneNumber = req.PhoneNumber

	token, err := h.clientService.SignInClient(&client, c)
	if err != nil {
		if err.Error() == "Not Found" {
			return c.JSON(utils.NotFoundError())
		}
		return c.JSON(utils.BadRequestError())
	}

	return c.JSON(http.StatusOK, newClientSignResponse(&client, token))
}

// Partner

// PartnerAuth godoc
//
//	@Summary		Sign-up partner
//	@Description	sign-up partner with given data
//	@Tags			Partners
//	@Accept			json
//	@Produce		json
//	@Param			partner	body		models.Partner	true	"Credentials to use"
//	@Success		201		{object}	handlers.partnerSignResponse
//	@Failure		400		{object}	utils.Err
//	@Failure		500		{object}	utils.Err
//	@Router			/partners/auth/sign-up [post]
func (h *AuthHandler) SignUpPartner(c echo.Context) error {
	var partner models.Partner
	if err := c.Bind(&partner); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error()) // 400
	}

	if err := h.validate.Struct(partner); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error()) // 400
	}

	token, err := h.partnerService.SignUpPartner(&partner, c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error()) // 400
	}

	return c.JSON(http.StatusCreated, newPartnerSignResponse(&partner, token))
}

// PartnerAuth godoc
//
//	@Summary		Sign-in for partners
//	@Description	sign-in in partners with given data
//	@Tags			Partners
//	@Accept			json
//	@Produce		json
//	@Param			partner	body		handlers.signInRequest	true	"Credentials to use"
//	@Success		200		{object}	handlers.partnerSignResponse
//	@Failure		400		{object}	utils.Err
//	@Failure		500		{object}	utils.Err
//	@Router			/partners/auth/sign-in [post]
func (h *AuthHandler) SignInPartner(c echo.Context) error {
	var partner models.Partner
	req := new(signInRequest)
	if err := req.validate(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Err{Message: err.Error()})
	}

	partner.Password = req.Password
	partner.PhoneNumber = req.PhoneNumber

	token, err := h.partnerService.SignInPartner(&partner, c)
	if err != nil {
		if err.Error() == "Not Found" {
			return c.JSON(utils.NotFoundError())
		}
		return c.JSON(utils.BadRequestError())
	}
	return c.JSON(http.StatusOK, newPartnerSignResponse(&partner, token))
}

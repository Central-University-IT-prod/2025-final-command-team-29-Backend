package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"backend/internal/models"
	usecase "backend/internal/use_case"
	"backend/internal/utils"
	tokenjwt "backend/pkg/token_jwt"
)

type PromoHandler struct {
	promoService *usecase.PromoService
}

func NewPromoHandler(promoService *usecase.PromoService) *PromoHandler {
	return &PromoHandler{
		promoService: promoService,
	}
}

func getIDFromToken(c echo.Context) (uuid.UUID, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return uuid.Nil, &utils.Err{Message: "Unauthorized"}
	}
	claims, ok := token.Claims.(*tokenjwt.Claims)
	if !ok {
		return uuid.Nil, &utils.Err{Message: "Unauthorized"}
	}
	return claims.UserID, nil
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Returning promo for client
//	@Tags			PromosForClient
//	@Produce		json
//	@Param			partnerID		path		string	true	"Partner ID"	format	uuid
//	@Param			promoID			path		string	true	"Partner ID"	format	uuid
//	@Param			Authorization	header		string	true	"Authorization header must be set for access"
//	@Success		200				{object}	models.Promo
//	@Failure		400				{object}	utils.Err
//	@Failure		401				{object}	utils.Err
//	@Failure		404				{object}	utils.Err
//	@Failure		500				{object}	utils.Err
//	@Router			/partners/{partnerID}/promos/{promoID} [get]
func (p *PromoHandler) GetPromo(c echo.Context) error {
	clientID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	partnerID, err := uuid.Parse(c.Param("partnerID"))
	if err != nil {
		return c.JSON(utils.BadRequestError())
	}
	promoID, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		return c.JSON(utils.BadRequestError())
	}
	promo, err := p.promoService.GetPromo(promoID, partnerID, clientID)
	if err != nil {
		return c.JSON(utils.BadRequestError())
	}
	if promo == nil {
		return c.JSON(utils.NotFoundError())
	}
	return c.JSON(http.StatusOK, promo)
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Returning list of promos of specific partner for client
//	@Tags			PromosForClient
//	@Produce		json
//	@Param			partnerID		path		string	true	"Partner ID"	format	uuid
//	@Param			Authorization	header		string	true	"Authorization header must be set for access"
//	@Success		200				{object}	[]models.Promo
//	@Failure		400				{object}	utils.Err
//	@Failure		500				{object}	utils.Err
//	@Router			/partners/{partnerID}/promos [get]
func (p *PromoHandler) GetPartnerPromos(c echo.Context) error {
	clientID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	partnerID, err := uuid.Parse(c.Param("partnerID"))
	if err != nil {
		return c.JSON(utils.BadRequestError())
	}
	promos, err := p.promoService.GetPartnerPromosForClient(clientID, partnerID)
	if err != nil {
		return c.JSON(utils.BadRequestError())
	}
	return c.JSON(http.StatusOK, promos)
}

// GetPromoImage godoc
//
//	@Summary		Get promo image
//	@Description	Returns the image associated with a specific promo
//	@Tags			PromosForClient
//	@Produce		octet-stream
//	@Param			partnerID		path		string	true	"Partner ID"	format	uuid
//	@Param			promoID			path		string	true	"Promo ID"	format	uuid
//	@Param			Authorization	header		string	true	"Authorization header must be set for access"
//	@Success		200				{file}		binary	"Image file data"
//	@Failure		400				{object}	map[string]string	"error: Bad request"
//	@Failure		401				{object}	map[string]string	"error: Unauthorized"
//	@Failure		404				{object}	map[string]string	"error: Изображение не найдено"
//	@Router			/partners/{partnerID}/promos/{promoID}/img [get]
func (p *PromoHandler) GetPromoImage(c echo.Context) error {
	_, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	partnerID := c.Param("partnerID")
	promoID := c.Param("promoID")

	// Получаем изображение из сервиса
	imageData, contentType, err := p.promoService.GetImage(promoID, partnerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Изображение не найдено"})
	}

	// Устанавливаем заголовок Content-Type и отправляем изображение
	return c.Blob(http.StatusOK, contentType, imageData)
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Create promo
//	@Tags			PromosForPartners
//	@Produce		json
//	@Param			promo			body		models.Promo	true	"Promo"
//	@Param			Authorization	header		string			true	"Authorization header must be set for access"
//	@Success		201				{object}	models.Promo
//	@Failure		400				{object}	utils.Err
//	@Failure		500				{object}	utils.Err
//	@Router			/promos [post]
func (p *PromoHandler) CreatePromo(c echo.Context) error {
	partnerID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	fmt.Println(partnerID)
	promo := &models.Promo{PartnerID: partnerID}
	if err := c.Bind(promo); err != nil {
		return c.JSON(utils.BadRequestError())
	}
	if err := c.Validate(promo); err != nil {
		return c.JSON(utils.BadRequestError())
	}
	if err := p.promoService.CreatePromo(promo); err != nil {
		return c.JSON(utils.BadRequestError())
	}
	return c.JSON(http.StatusCreated, promo)
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Delete promo
//	@Tags			PromosForPartners
//	@Produce		json
//	@Param			promoID			path	string	true	"Promo ID"
//	@Param			Authorization	header	string	true	"Authorization header must be set for access"
//	@Success		204
//	@Failure		400	{object}	utils.Err
//	@Failure		401	{object}	utils.Err
//	@Failure		500	{object}	utils.Err
//	@Router			/promos/{promoID} [delete]
func (p *PromoHandler) DeletePromo(c echo.Context) error {
	partnerID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	promoID, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		return c.JSON(utils.BadRequestError())
	}

	if err := p.promoService.DeletePromo(promoID, partnerID); err != nil {
		return c.JSON(utils.BadRequestError())
	}
	return c.JSON(http.StatusNoContent, nil)
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Get promo by partner
//	@Tags			PromosForPartners
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization header must be set for access"
//	@Success		200				{object}	[]models.Promo
//	@Failure		400				{object}	utils.Err
//	@Failure		401				{object}	utils.Err
//	@Router			/promos [get]
func (p *PromoHandler) GetParnterPromosProtected(c echo.Context) error {
	partnerID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	promos, err := p.promoService.GetPartnerPromos(partnerID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, promos)
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Get promo by partner
//	@Tags			PromosForPartners
//	@Produce		json
//	@Param			promoID			path		string	true	"Promo ID"
//	@Param			Authorization	header		string	true	"Authorization header must be set for access"
//	@Success		200				{object}	models.Promo
//	@Failure		400				{object}	utils.Err
//	@Failure		401				{object}	utils.Err
//	@Router			/promos/{promoID} [get]
func (p *PromoHandler) GetPromoForPartner(c echo.Context) error {
	partnerID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}
	promoID, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	promo, err := p.promoService.GetPromoForPartner(promoID, partnerID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	if promo == nil {
		return c.JSON(utils.NotFoundError())
	}
	return c.JSON(http.StatusOK, promo)
}

// Promo godoc
//
//	@Summary		Promo endpoint
//	@Description	Get promo by partner
//	@Tags			PromosForPartners
//	@Produce		json
//	@Param			promoID			path		string	true	"Promo ID"
//	@Param			clientID		path		string	true	"Client ID"
//	@Param			Authorization	header		string	true	"Authorization header must be set for access"
//	@Success		200				{object}	models.PromoCLientView
//	@Failure		400				{object}	utils.Err
//	@Failure		401				{object}	utils.Err
//	@Router			/promos/{promoID}/{clientID} [get]
func (p *PromoHandler) GetPromoForPartnerWithClientData(c echo.Context) error {
	partnerID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}
	promoID, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	clientID, err := uuid.Parse(c.Param("clientID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	promo, err := p.promoService.GetPromo(promoID, partnerID, clientID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	if promo == nil {
		return c.JSON(utils.NotFoundError())
	}
	return c.JSON(http.StatusOK, promo)
}

// UploadImage godoc
//
//	@Summary		Upload promo image
//	@Description	Allows partners to upload an image for a specific promo
//	@Tags			PromosForPartners
//	@Accept			mpfd
//	@Produce		json
//	@Param			promoID			path		string			true	"Promo ID"	format	uuid
//	@Param			file			formData	file			true	"Image file to upload"
//	@Param			Authorization	header		string			true	"Authorization header must be set for access"
//	@Success		200				{object}	map[string]string	"message: Файл успешно загружен"
//	@Failure		400				{object}	map[string]string	"error: Файл обязателен"
//	@Failure		401				{object}	map[string]string	"error: Unauthorized"
//	@Failure		500				{object}	map[string]string	"error: Внутренняя ошибка сервера"
//	@Router			/promos/{promoID}/img/upload [post]
func (p *PromoHandler) UploadImage(c echo.Context) error {
	partnerID, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	promoID := c.Param("promoID")
	fmt.Printf("promoID:%s", promoID)

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Файл обязателен"})
	}

	if err := p.promoService.UploadImage(promoID, partnerID.String(), file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Файл успешно загружен"})
}

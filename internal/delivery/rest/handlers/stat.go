package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	_ "backend/internal/models" // Импорт модели Stat
	usecase "backend/internal/use_case"
	"backend/internal/utils"
)

type StatHandler struct {
	statService *usecase.StatService
	validate    *validator.Validate
}

func NewStatsrHandler(statService *usecase.StatService) *StatHandler {
	return &StatHandler{
		statService: statService,
		validate:    validator.New(),
	}
}

// GetStatList возвращает бизнес-статистику по промо-акциям партнера.
//
//	@Summary		Получение статистики по промо-акциям
//	@Description	Возвращает метрики использования промо-кодов, такие как активные пользователи, количество покупок, retention и т. д.
//	@Tags			Statistics
//	@Accept			json
//	@Produce		json
//	@Param			partnerID	path		string					true	"UUID партнера"
//	@Success		200			{object}	models.Stat				"Успешный ответ со статистикой"
//	@Success		200			{object}	models.Stat				"Успешный ответ со статистикой"
//	@Failure		400			{object}	models.ErrorResponse	"Ошибка валидации данных. Пример: {\"error\": \"Некорректный	UUID	партнера\"}"
//	@Failure		401			{object}	models.ErrorResponse	"Ошибка авторизации. Пример: {\"error\": \"Некорректный			токен	или			пользователь	не	аутентифицирован\"}"
//	@Failure		403			{object}	models.ErrorResponse	"Ошибка доступа. Пример: {\"error\": \"Недостаточно				прав	для			просмотра		статистики\"}"
//	@Failure		500			{object}	models.ErrorResponse	"Внутренняя ошибка сервера. Пример: {\"error\": \"Ошибка		при		получении	статистики\"}"
//	@Router			/partners/{partnerID}/stat [get]
func (h *StatHandler) GetStatList(c echo.Context) error {
	id, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.Err{Message: err.Error()})
	}
	partnerID, err := uuid.Parse(c.Param("partnerID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Err{Message: "Invalid UUID"})
	}
	if partnerID != id {
		return c.JSON(http.StatusForbidden, utils.Err{Message: "Forbidden"})
	}
	stat, err := h.statService.GetStatList(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, stat)
}

// Stats godoc
//
//	@Summary		Stats endpoint
//	@Description	Stats for partner for specific promo
//	@Tags			Stats
//	@Produce		json
//	@Param			partnerID	path		string	true	"Partner ID"
//	@Param			promoID		path		string	true	"Promo ID"
//	@Success		200			{object}	models.StatPromo
//	@Failure		400			{object}	utils.Err
//	@Failure		401			{object}	utils.Err
//	@Failure		500			{object}	utils.Err
//	@Router			/partners/{partnerID}/stat/{promoID} [get]
func (h *StatHandler) GetStatByPromo(c echo.Context) error {
	// Получаем ID пользователя из токена
	id, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Получаем partnerID из URL
	partnerID, err := uuid.Parse(c.Param("partnerID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Err{Message: "Invalid partner UUID"})
	}

	// Проверяем, что partnerID соответствует ID пользователя
	if partnerID != id {
		return c.JSON(http.StatusForbidden, utils.Err{Message: "Forbidden"})
	}

	// Получаем promoID из URL
	promoID, err := uuid.Parse(c.Param("promoID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Err{Message: "Invalid promo UUID"})
	}

	// Получаем статистику по промокоду
	stat, err := h.statService.GetStatByPromo(partnerID, promoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, stat)
}

// StatExport экспортирует статистику партнера в формате Excel
// @Summary Экспорт статистики партнера
// @Description Генерирует Excel-файл с активациями промокодов партнера
// @Tags Statistics
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param partnerID path string true "UUID партнера"
// @Security BearerAuth
// @Success 200 {file} application/vnd.openxmlformats-officedocument.spreadsheetml.sheet "Файл Excel"
// @Failure 400 {object} utils.Err "Некорректный UUID"
// @Failure 401 {object} utils.Err "Неавторизован"
// @Failure 403 {object} utils.Err "Доступ запрещен"
// @Failure 500 {object} utils.Err "Ошибка сервера"
// @Router /partners/{partnerID}/stat/export [get]
func (h *StatHandler) StatExport(c echo.Context) error {
	// Получаем ID партнера из токена
	id, err := getIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.Err{Message: err.Error()})
	}

	// Получаем partnerID из параметра запроса
	partnerID, err := uuid.Parse(c.Param("partnerID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Err{Message: "Invalid UUID"})
	}

	// Вызываем сервис для экспорта
	filePath, err := h.statService.ExportPartnerStat(partnerID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Err{Message: err.Error()})
	}

	// Отправляем Excel-файл пользователю
	return c.Attachment(filePath, "partner_stats.xlsx")
}

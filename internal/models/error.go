package models

// ErrorResponse представляет формат ответа в случае ошибки.
//
// @Description Структура ошибки, возвращаемая API в формате { "error": "message" }
type ErrorResponse struct {
	Error string `json:"error" example:"Сообщение об ошибке"`
}

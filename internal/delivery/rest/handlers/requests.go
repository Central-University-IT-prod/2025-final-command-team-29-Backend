package handlers

import (
	"github.com/labstack/echo/v4"
)

type signInRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Password    string `json:"password" validate:"required"`
}

func (req *signInRequest) validate(c echo.Context) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	return nil
}

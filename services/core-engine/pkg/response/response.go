package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

type DefaultResponse struct {
	Code   int         `json:"code"`
	Status interface{} `json:"msg"`
	Data   interface{} `json:"payload"`
}

func ResponseSuccess(c echo.Context, data interface{}) error {
	response := DefaultResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   data,
	}

	return c.JSON(http.StatusOK, response)
}

func ResponseFailed(c echo.Context, err interface{}, code int) error {
	response := DefaultResponse{
		Code:   code,
		Status: cast.ToString(err),
		Data:   nil,
	}

	return c.JSON(code, response)
}

package response

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

type DefaultResponse struct {
	Code   int `json:"code"`
	Status any `json:"status"`
	Data   any `json:"data"`
	Meta   any `json:"meta,omitempty"`
}

type meta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	MaxPage   int `json:"maxPage"`
	TotalItem int `json:"totalItem"`
}

func Success(c echo.Context, data any) error {
	response := DefaultResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   data,
	}

	return c.JSON(http.StatusOK, response)
}

func SuccessList(c echo.Context, data any, page, limit, totalItem int) error {
	response := DefaultResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   data,
		Meta: meta{
			Page:      page,
			Limit:     limit,
			MaxPage:   int(math.Ceil(float64(totalItem) / float64(limit))),
			TotalItem: totalItem,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func Failed(c echo.Context, err any, code int) error {
	response := DefaultResponse{
		Code:   code,
		Status: cast.ToString(err),
		Data:   nil,
	}

	return c.JSON(code, response)
}

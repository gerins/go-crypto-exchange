package response

import (
	"context"
	"math"
	"net/http"

	"github.com/gerins/log"
	"github.com/labstack/echo/v4"

	serverError "core-engine/pkg/error"
)

type DefaultResponse struct {
	Code    int `json:"code"`
	Message any `json:"message"`
	Data    any `json:"data"`
	Meta    any `json:"meta,omitempty"`
}

type meta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	MaxPage   int `json:"maxPage"`
	TotalItem int `json:"totalItem"`
}

func Success(c echo.Context, data any) error {
	response := DefaultResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	}

	return c.JSON(http.StatusOK, response)
}

func SuccessList(c echo.Context, data any, page, limit, totalItem int) error {
	response := DefaultResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
		Meta: meta{
			Page:      page,
			Limit:     limit,
			MaxPage:   int(math.Ceil(float64(totalItem) / float64(limit))),
			TotalItem: totalItem,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func Failed(c echo.Context, err error) error {
	var (
		generalError = serverError.ErrGeneralError(nil)
		httpRespCode = generalError.HTTPCode
		rawError     = err
	)

	response := DefaultResponse{
		Code:    generalError.Code,
		Message: generalError.Message,
		Data:    nil,
	}

	// Check for wrapped server error
	if serverErr, ok := err.(serverError.ServerError); ok {
		rawError = serverErr.RawError
		httpRespCode = serverErr.HTTPCode
		response.Code = serverErr.Code
		response.Message = serverErr.Message
	}

	log.Context(c.Get("ctx").(context.Context)).ExtraData["error"] = rawError.Error()
	return c.JSON(httpRespCode, response)
}

package response

import (
	"errors"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

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
	)

	response := DefaultResponse{
		Code:    generalError.Code,
		Message: generalError.Message,
		Data:    nil,
	}

	// Default response for data not found in database
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = serverError.ErrDataNotFound(err)
	}

	// Check for wrapped server error
	if serverErr, ok := err.(serverError.ServerError); ok {
		httpRespCode = serverErr.HTTPCode
		response.Code = serverErr.Code
		response.Message = serverErr.Message
	}

	return c.JSON(httpRespCode, response)
}

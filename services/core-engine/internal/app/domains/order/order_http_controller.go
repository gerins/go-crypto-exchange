package order

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"core-engine/config"
	"core-engine/internal/app/domains/order/model"
	httpMiddleware "core-engine/internal/app/middleware/http"
	"core-engine/pkg/response"
)

type httpHandler struct {
	timeout      time.Duration
	orderUsecase model.Usecase
}

func NewHTTPHandler(orderUsecase model.Usecase, timeout time.Duration) interface {
	InitRoutes(e *echo.Echo, securityConfig config.Security)
} {
	return &httpHandler{
		timeout:      timeout,
		orderUsecase: orderUsecase,
	}
}

func (h *httpHandler) InitRoutes(e *echo.Echo, securityConfig config.Security) {
	v1 := e.Group("/api/v1/order")
	v1.Use(httpMiddleware.ValidateJwtToken([]byte(securityConfig.Jwt.Key)))
	{
		v1.POST("", h.OrderHandler)
	}
}

func (h *httpHandler) OrderHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Get("ctx").(context.Context), h.timeout)
	defer cancel()

	var requestPayload model.RequestOrder
	if err := c.Bind(&requestPayload); err != nil {
		return response.Failed(c, err, http.StatusBadRequest)
	}

	orderResult, err := h.orderUsecase.ProcessOrder(ctx, requestPayload)
	if err != nil {
		return response.Failed(c, err, http.StatusBadRequest)
	}

	return response.Success(c, orderResult)
}

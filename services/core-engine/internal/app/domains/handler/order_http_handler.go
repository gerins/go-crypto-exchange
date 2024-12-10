package handler

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"

	"core-engine/config"
	"core-engine/internal/app/domains/dto"
	"core-engine/internal/app/domains/model"
	httpMiddleware "core-engine/internal/app/middleware/http"
	"core-engine/pkg/response"
)

type orderHandler struct {
	timeout        time.Duration
	orderUsecase   model.OrderUsecase
	securityConfig config.Security
}

func NewOrderHTTPHandler(orderUsecase model.OrderUsecase, timeout time.Duration, securityConfig config.Security) interface {
	InitRoutes(e *echo.Echo)
} {
	return &orderHandler{
		timeout:        timeout,
		orderUsecase:   orderUsecase,
		securityConfig: securityConfig,
	}
}

func (h *orderHandler) InitRoutes(e *echo.Echo) {
	v1 := e.Group("/api/v1/order")
	v1.Use(httpMiddleware.ValidateJwtToken([]byte(h.securityConfig.Jwt.Key)))
	{
		v1.POST("", h.OrderHandler)
	}
}

func (h *orderHandler) OrderHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Get("ctx").(context.Context), h.timeout)
	defer cancel()

	var requestPayload dto.OrderRequest
	if err := c.Bind(&requestPayload); err != nil {
		return response.Failed(c, err)
	}

	orderResult, err := h.orderUsecase.ProcessOrder(ctx, requestPayload)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.Success(c, orderResult)
}

package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"matching-engine/internal/app/model"
	"matching-engine/pkg/response"
)

type httpHandler struct {
	engine  model.Engine
	timeout time.Duration
}

func NewHTTPHandler(processor model.Engine, timeout time.Duration) interface{ InitRoutes(e *echo.Echo) } {
	return &httpHandler{
		engine:  processor,
		timeout: timeout,
	}
}

func (h *httpHandler) InitRoutes(e *echo.Echo) {
	v1 := e.Group("/v1/order")
	{
		v1.POST("", h.OrderHandler)
	}
}

func (h *httpHandler) OrderHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	var requestPayload model.Order
	if err := c.Bind(&requestPayload); err != nil {
		return response.ResponseFailed(c, err, http.StatusBadRequest)
	}

	orderResult := h.engine.Execute(ctx, requestPayload)

	return response.ResponseSuccess(c, orderResult)
}

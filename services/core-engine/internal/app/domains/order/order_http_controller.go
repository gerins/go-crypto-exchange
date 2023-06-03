package order

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"core-engine/internal/app/domains/order/model"
	"core-engine/pkg/response"
)

type httpHandler struct {
	timeout      time.Duration
	orderUsecase model.Usecase
}

func NewHTTPHandler(orderUsecase model.Usecase, timeout time.Duration) interface{ InitRoutes(e *echo.Echo) } {
	return &httpHandler{
		timeout:      timeout,
		orderUsecase: orderUsecase,
	}
}

func (h *httpHandler) InitRoutes(e *echo.Echo) {
	v1 := e.Group("/api/v1/order")
	{
		v1.POST("", h.OrderHandler)
	}
}

func (h *httpHandler) OrderHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
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

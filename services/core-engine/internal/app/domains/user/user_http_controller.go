package user

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"core-engine/pkg/response"
)

type httpHandler struct {
	timeout     time.Duration
	userUsecase Usecase
}

func NewHTTPHandler(userUsecase Usecase, timeout time.Duration) interface{ InitRoutes(e *echo.Echo) } {
	return &httpHandler{
		timeout:     timeout,
		userUsecase: userUsecase,
	}
}

func (h *httpHandler) InitRoutes(e *echo.Echo) {
	v1 := e.Group("/v1/user")
	{
		v1.POST("/login", h.LoginHandler)
		v1.POST("/register", h.RegisterHandler)
	}
}

func (h *httpHandler) LoginHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	var requestPayload LoginRequest
	if err := c.Bind(&requestPayload); err != nil {
		return response.ResponseFailed(c, err, http.StatusBadRequest)
	}

	loginResult, err := h.userUsecase.Login(ctx, requestPayload)
	if err != nil {
		return response.ResponseFailed(c, err, http.StatusBadRequest)
	}

	return response.ResponseSuccess(c, loginResult)
}

func (h *httpHandler) RegisterHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	var requestPayload RegisterRequest
	if err := c.Bind(&requestPayload); err != nil {
		return response.ResponseFailed(c, err, http.StatusBadRequest)
	}

	registerResult, err := h.userUsecase.Register(ctx, requestPayload)
	if err != nil {
		return response.ResponseFailed(c, err, http.StatusBadRequest)
	}

	return response.ResponseSuccess(c, registerResult)
}

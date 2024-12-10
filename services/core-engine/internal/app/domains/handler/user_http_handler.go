package handler

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"

	"core-engine/internal/app/domains/dto"
	"core-engine/internal/app/domains/model"
	"core-engine/pkg/response"
)

type userHandler struct {
	timeout     time.Duration
	userUsecase model.UserUsecase
}

func NewUserHandler(userUsecase model.UserUsecase, timeout time.Duration) interface{ InitRoutes(e *echo.Echo) } {
	return &userHandler{
		timeout:     timeout,
		userUsecase: userUsecase,
	}
}

func (h *userHandler) InitRoutes(e *echo.Echo) {
	v1 := e.Group("/api/v1/user")
	{
		v1.POST("/login", h.LoginHandler)
		v1.POST("/register", h.RegisterHandler)
	}
}

func (h *userHandler) LoginHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Get("ctx").(context.Context), h.timeout)
	defer cancel()

	var requestPayload dto.LoginRequest
	if err := c.Bind(&requestPayload); err != nil {
		return response.Failed(c, err)
	}

	loginResult, err := h.userUsecase.Login(ctx, requestPayload)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.Success(c, loginResult)
}

func (h *userHandler) RegisterHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Get("ctx").(context.Context), h.timeout)
	defer cancel()

	var requestPayload dto.RegisterRequest
	if err := c.Bind(&requestPayload); err != nil {
		return response.Failed(c, err)
	}

	registerResult, err := h.userUsecase.Register(ctx, requestPayload)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.Success(c, registerResult)
}

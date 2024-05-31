package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"core-engine/pkg/jwt"
	"core-engine/pkg/response"
)

func ValidateRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get parent context from Echo Locals
			ctx, ok := c.Get("ctx").(context.Context)
			if !ok {
				ctx = context.Background()
			}

			jwtPayload := jwt.GetPayloadFromContext(ctx)

			for _, role := range roles {
				if role == jwtPayload.Role {
					return next(c)
				}
			}

			return response.Failed(c, "unauthorized user", http.StatusUnauthorized)
		}
	}
}

package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	uuid "github.com/satori/go.uuid"
)

func RequestUUID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.NewV4().String()
		},
		RequestIDHandler: func(ctx echo.Context, s string) {
			ctx.Request().Header.Set(echo.HeaderXRequestID, s)
		},
	})
}

package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoute(e *echo.Echo) {
	e.Any("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello 妈咪")
	})
}

package welcome

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Vault API")
	})
	e.GET("/v1", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Version 1 Vault API")
	})
}

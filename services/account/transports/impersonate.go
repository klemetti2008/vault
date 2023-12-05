package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) impersonate(ctx echo.Context) error {
	id := ctx.Param("id")
	accessToken := ctx.Request().Header.Get("Authorization")
	accessToken = accessToken[len("Bearer "):]

	loginResponse, err := r.service.Impersonate(ctx.Request().Context(), id, accessToken)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(loginResponse))
}

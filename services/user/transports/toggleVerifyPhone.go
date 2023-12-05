package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) toggleVerifyPhone(ctx echo.Context) error {
	id := ctx.Param("id")
	user, err := r.service.ToggleVerifyPhone(ctx.Request().Context(), id)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(user))
}

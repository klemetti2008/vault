package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r *resource) get(ctx echo.Context) error {
	id := ctx.Param("id")
	category, err := r.service.Get(ctx.Request().Context(), id)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(category))
}

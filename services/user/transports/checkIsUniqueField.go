package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) checkIsUniqueField(ctx echo.Context) error {
	field := ctx.Param("field")
	value := ctx.Param("value")
	exists, _, err := r.service.CheckIsUniqueField(ctx.Request().Context(), field, value)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusOK, response.Success(exists))
}

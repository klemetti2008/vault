package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) approvePhone(ctx echo.Context) error {
	code := ctx.Param("code")

	ok, err := r.service.ApprovePhone(ctx.Request().Context(), code)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(ok))
}

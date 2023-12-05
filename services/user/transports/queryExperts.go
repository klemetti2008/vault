package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) queryExperts(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	context := ctx.Request().Context()
	users, er := r.service.QueryExperts(context, query)
	if er.StatusCode != 0 {
		return ctx.JSON(er.StatusCode, er)
	}
	return ctx.JSON(http.StatusOK, response.Success(users))
}

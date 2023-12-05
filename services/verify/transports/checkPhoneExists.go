package transports

import (
	"net/http"

	"gitag.ir/cookthepot/services/vault/services/verify/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) checkPhoneExists(ctx echo.Context) error {
	var input = &endpoints.CheckPhoneExistsRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	exists, err := r.service.CheckPhoneExists(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(exists))
}

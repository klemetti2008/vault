package transports

import (
	"net/http"

	"gitag.ir/cookthepot/services/vault/services/account/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) resetPassword(ctx echo.Context) error {
	var input = &endpoints.ResetPasswordRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	loginResponse, err := r.service.ResetPassword(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(loginResponse))
}

package transports

import (
	"net/http"

	"gitag.ir/cookthepot/services/vault/services/account/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r resource) refreshToken(ctx echo.Context) error {
	// get refresh token from header
	var input = &endpoints.RefreshTokenRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	loginResponse, err := r.service.RefreshToken(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(loginResponse))
}

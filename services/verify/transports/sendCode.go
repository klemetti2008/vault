package transports

import (
	"net/http"

	"gitag.ir/cookthepot/services/vault/services/verify/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

type SendCodeResponse struct {
	Code string `json:"code"`
}

func (r resource) sendCode(ctx echo.Context) error {
	var input = &endpoints.SendCodeRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	code, err := r.service.SendCode(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	result := &SendCodeResponse{
		Code: code,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}

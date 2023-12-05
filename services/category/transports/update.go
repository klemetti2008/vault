package transports

import (
	"net/http"

	"gitag.ir/cookthepot/services/vault/services/category/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r *resource) update(ctx echo.Context) error {
	id := ctx.Param("id")
	var input = &endpoints.UpdateCategoryRequest{}
	errors := input.Validate(ctx.Request())

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	category, err := r.service.Update(ctx.Request().Context(), id, *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusCreated, response.Success(category))
}

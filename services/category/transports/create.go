package transports

import (
	"fmt"
	"net/http"

	"gitag.ir/cookthepot/services/vault/services/category/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

func (r *resource) create(ctx echo.Context) error {
	var input = &endpoints.CreateCategoryRequest{}

	fmt.Printf("input: %+v", *input)
	errors := input.Validate(ctx.Request())
	fmt.Printf("input: %+v", *input)

	if errors != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	category, err := r.service.Create(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	return ctx.JSON(http.StatusCreated, response.Created(category))
}

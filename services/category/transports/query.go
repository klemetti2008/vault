package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/pagination"
	"github.com/mhosseintaher/kit/response"
	"github.com/mhosseintaher/kit/restypes"
)

func (r *resource) query(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	order := ctx.QueryParam("order")
	orderBy := ctx.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "id"
	}
	if order == "" {
		order = "desc"
	}
	count, err := r.service.Count(ctx.Request().Context(), query)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	pages := pagination.NewFromRequest(ctx.Request(), int(count))
	categories, er := r.service.Query(
		ctx.Request().Context(),
		pages.Offset(), pages.Limit(), orderBy, order, query,
	)
	if er.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	pages.Items = categories

	result := restypes.QueryResponse{
		Limit:      pages.Limit(),
		Offset:     pages.Offset(),
		Page:       pages.Page,
		TotalRows:  int64(pages.TotalCount),
		TotalPages: pages.PageCount,
		Items:      categories,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}

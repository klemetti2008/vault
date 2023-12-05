package transports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/pagination"
	"github.com/mhosseintaher/kit/response"
	"github.com/mhosseintaher/kit/restypes"
)

func (r resource) query(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	order := ctx.QueryParam("order")
	orderBy := ctx.QueryParam("order_by")
	suspendedAt := ctx.QueryParam("suspended_at")
	isOfficial := ctx.QueryParam("is_official")
	isProfileCompleted := ctx.QueryParam("is_profile_completed")
	if orderBy == "" {
		orderBy = "id"
	}
	if order == "" {
		order = "desc"
	}

	c := ctx.Request().Context()

	count, err := r.service.Count(c, query, suspendedAt, isOfficial, isProfileCompleted)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	pages := pagination.NewFromRequest(ctx.Request(), int(count))
	users, er := r.service.Query(
		c, pages.Offset(), pages.Limit(),
		orderBy, order, query, suspendedAt, isOfficial, isProfileCompleted,
	)
	if er.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}
	result := restypes.QueryResponse{
		Limit:      pages.Limit(),
		Offset:     pages.Offset(),
		Page:       pages.Page,
		TotalRows:  int64(pages.TotalCount),
		TotalPages: pages.PageCount,
		Items:      users,
	}
	return ctx.JSON(http.StatusOK, response.Success(result))
}

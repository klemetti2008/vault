package role

import (
	"net/http"
	"path/filepath"

	"gitag.ir/cookthepot/services/vault/services/mid"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/log"
	"github.com/mhosseintaher/kit/response"
)

func RegisterHandlers(r *echo.Echo, service Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(filepath.Join("/api", prefix))
	g.GET("/roles/list", res.list)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) list(ctx echo.Context) error {
	c := ctx.Request().Context()
	roles, er := r.service.Query(c)
	if er != nil {
		return er
	}
	return ctx.JSON(http.StatusOK, response.Success(roles))
}

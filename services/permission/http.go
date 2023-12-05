package permission

import (
	"net/http"

	"gitag.ir/cookthepot/services/vault/services/mid"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/log"
	"github.com/mhosseintaher/kit/response"
)

func RegisterHandlers(r *echo.Echo, service Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)
	rg.GET("/permissions/domains-access", res.accessList)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) accessList(ctx echo.Context) error {
	c := ctx.Request().Context()
	acl, er := r.service.AccessList(c)
	if er != nil {
		return er
	}
	return ctx.JSON(http.StatusOK, response.Success(acl))
}

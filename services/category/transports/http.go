package transports

import (
	"gitag.ir/cookthepot/services/vault/services/category/endpoints"
	"gitag.ir/cookthepot/services/vault/services/mid"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/log"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/categories/:id", res.get)
	rg.GET("/categories", res.query)
	rg.POST("/categories", res.create)
	rg.Match([]string{"PUT", "PATCH"}, "/categories/:id", res.update)
	rg.DELETE("/categories", res.delete)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}

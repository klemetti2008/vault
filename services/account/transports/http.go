package transports

import (
	"gitag.ir/cookthepot/services/vault/services/account/endpoints"
	"gitag.ir/cookthepot/services/vault/services/mid"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/log"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)
	g.POST("/login", res.login)
	g.POST("/register", res.register)
	g.POST("/resetPassword", res.resetPassword)

	g.DELETE("/tokens", res.deleteTokens)
	g.POST("/tokens/refresh", res.refreshToken)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.POST("/logout", res.logout)
	rg.GET("/userinfo", res.userInfo)
	rg.PATCH("/change-password", res.changePassword)
	rg.POST("/impersonate/:id", res.impersonate)

	rg.GET("/tokens/:user_id", res.getAllTokensByUserId)

	rg.POST("/accounts/approve-phone/:code", res.approvePhone)
	rg.POST("/accounts/approve-email/:token", res.approveEmail)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}

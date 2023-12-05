package account

import (
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/log"

	"gitag.ir/cookthepot/services/vault/notification"
	"gitag.ir/cookthepot/services/vault/services/account/endpoints"
	"gitag.ir/cookthepot/services/vault/services/account/transports"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, logger log.Logger, notifier notification.Notifier, AccessTokenSigningKey string, RefreshTokenSigningKey string, AccessTokenTokenExpiration int, RefreshTokenExpiration int, prefix string) {

	service := endpoints.MakeService(db, logger, notifier, AccessTokenSigningKey, RefreshTokenSigningKey, AccessTokenTokenExpiration, RefreshTokenExpiration)
	transports.RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))

}

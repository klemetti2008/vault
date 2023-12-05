package verify

import (
	"path/filepath"

	"gitag.ir/cookthepot/services/vault/notification"
	"gitag.ir/cookthepot/services/vault/services/verify/endpoints"
	"gitag.ir/cookthepot/services/vault/services/verify/transports"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/log"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, notifier notification.Notifier, logger log.Logger, prefix string) {
	service := endpoints.MakeService(db, logger, notifier)
	transports.RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))
}

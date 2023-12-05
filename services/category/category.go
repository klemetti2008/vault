package category

import (
	"path/filepath"

	"gitag.ir/cookthepot/services/vault/services/category/endpoints"
	"gitag.ir/cookthepot/services/vault/services/category/transports"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/log"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, logger log.Logger, prefix string) {
	service := endpoints.MakeService(db, logger)
	transports.RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))
}

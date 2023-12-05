package role

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/log"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context) (roles []models.Role, err error)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{db, logger}
}

func (s *service) Query(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := s.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}

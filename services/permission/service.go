package permission

import (
	"context"

	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/log"
	"gorm.io/gorm"
)

type Service interface {
	AccessList(ctx context.Context) (acl map[string]bool, err error)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{db, logger}
}

func (s *service) AccessList(ctx context.Context) (map[string]bool, error) {
	acl := map[string]bool{
		"CanAccessUser": policy.CanAccessUser(ctx),
	}
	return acl, nil
}

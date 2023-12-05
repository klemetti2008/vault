package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/log"
	"github.com/mhosseintaher/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Get(ctx context.Context, id string) (category models.Category, err response.ErrorResponse)
	Query(ctx context.Context, offset, limit int, orderBy, order, query string) (
		categories []models.Category, err response.ErrorResponse,
	)
	Count(ctx context.Context, query string) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreateCategoryRequest) (category models.Category, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdateCategoryRequest) (category models.Category, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{
		db:     db,
		logger: logger,
	}
}

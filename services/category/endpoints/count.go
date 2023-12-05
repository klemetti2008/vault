package endpoints

import (
	"context"
	"fmt"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Count(ctx context.Context, query string) (int64, response.ErrorResponse) {
	var count int64
	query = fmt.Sprintf("%%%s%%", query)
	err := s.db.WithContext(ctx).Model(&models.Category{}).Where("title LIKE ?", query).Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return count, response.GormErrorResponse(err, "An error occurred while calculating the data")
	}
	return count, response.ErrorResponse{}
}

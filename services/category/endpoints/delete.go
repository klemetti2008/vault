package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	if !policy.CanDeleteCategory(ctx) {
		s.logger.With(ctx).Error("You do not have permission to delete a category")
		return []int{}, response.ErrorForbidden("You do not have permission to delete a category")
	}

	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.Category{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "Error in deleting category")
	}

	return ids, response.ErrorResponse{}
}

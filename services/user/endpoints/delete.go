package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	if !policy.CanDeleteUser(ctx) {
		return []int{}, response.ErrorForbidden("You do not have permission to access this user")
	}
	err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.User{}).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "An error occurred while deleting users")
	}
	return ids, response.ErrorResponse{}
}

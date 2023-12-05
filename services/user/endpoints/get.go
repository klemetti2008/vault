package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Get(ctx context.Context, id string) (models.User, response.ErrorResponse) {
	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Roles").
		First(&user, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "Error in finding the user")
	}

	if !policy.CanGetUser(ctx, user) {
		s.logger.With(ctx).Error("You do not have permission to access this user")
		return models.User{}, response.ErrorForbidden("You do not have permission to access this user")
	}
	return user, response.ErrorResponse{}
}

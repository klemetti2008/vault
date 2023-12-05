package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/faker"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) ToggleVerifyEmail(ctx context.Context, id string) (models.User, response.ErrorResponse) {
	var user models.User
	if !policy.CanToggleVerifyEmail(ctx) {
		s.logger.With(ctx).Error("You do not have permission to access this user")
		err := response.ErrorForbidden("You do not have permission to access this user")
		return user, err
	}

	err := s.db.WithContext(ctx).First(&user, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "Error in finding the user")
	}

	if user.EmailVerifiedAt.Valid {
		user.EmailVerifiedAt = dtp.NullTime{}
	} else {
		user.EmailVerifiedAt = faker.SQLNow()
	}

	err = s.db.WithContext(ctx).Save(&user).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "Error in saving the user")
	}

	return user, response.ErrorResponse{}
}

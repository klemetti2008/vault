package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) UserInfo(ctx context.Context, accessToken string) (models.User, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	var user models.User
	var token models.Token
	err := s.db.WithContext(ctx).First(&token, "access_token", accessToken).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.ErrorUnAuthorized(nil, "the sent token is not valid")
	}

	err = s.db.WithContext(ctx).
		Preload("Roles").
		First(&user, "id", Id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.ErrorUnAuthorized(nil, "cannot find a user with these credentials")
	}

	if user.SuspendedAt.Valid {
		s.logger.With(ctx).Error(err)
		return user, response.ErrorUnAuthorized(nil, "your account has been suspended")
	}

	return user, response.ErrorResponse{}
}

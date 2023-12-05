package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Impersonate(ctx context.Context, id string, currentAccessToken string) (
	LoginResponse,
	response.ErrorResponse,
) {
	var loginResponse LoginResponse
	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Roles").
		First(&user, "id", id).Error
	if !policy.CanImpersonate(ctx) {
		s.logger.With(ctx).Error("You do not have permission to access this section")
		return loginResponse, response.ErrorForbidden(err, "You do not have permission to access this section")
	}

	if err != nil {
		s.logger.With(ctx).Error(err)
		return loginResponse, response.GormErrorResponse(err, "Error in finding the user")
	}

	accessToken, refreshToken, responseError := s.generateTokens(ctx, user, currentAccessToken)
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(responseError)
		return loginResponse, responseError
	}
	loginResponse = LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return loginResponse, response.ErrorResponse{}
}

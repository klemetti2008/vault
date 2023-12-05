package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Logout(ctx context.Context, accessTokens string) (string, response.ErrorResponse) {
	var token models.Token
	tx := s.db.WithContext(ctx).
		Where("access_token = ?", accessTokens).First(&token)
	if tx.Error != nil {
		s.logger.With(ctx).Error(tx.Error)
		return "", response.GormErrorResponse(tx.Error, "Error in logout")
	}

	err := tx.Delete(&models.Token{}).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "Error in logout")
	}

	return token.AccessToken, response.ErrorResponse{}
}

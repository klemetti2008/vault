package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) GetAllTokensByUserId(ctx context.Context, userId string) ([]models.Token, response.ErrorResponse) {
	var tokens []models.Token
	err := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&tokens).Error

	if !policy.CanGetAllTokensByUserId(ctx, tokens) {
		s.logger.With(ctx).Error("You do not have permission to access this section")
		return tokens, response.ErrorForbidden(err, "You do not have permission to access this section")
	}

	return tokens, response.ErrorResponse{}
}

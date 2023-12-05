package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
)

func (s *service) DeleteTokens(ctx context.Context, accessTokens []string) ([]string, error) {
	var res []string
	err := s.db.WithContext(ctx).
		Where("access_token IN ?", accessTokens).
		Delete(&models.Token{}).Error
	if err != nil {
		return res, err
	}

	res = accessTokens
	return res, nil
}

package endpoints

import (
	"context"
	"fmt"

	"gitag.ir/cookthepot/services/vault/database"
	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Count(
	ctx context.Context, query,
	suspendedAt, isOfficial, isProfileComplete string,
) (int64, response.ErrorResponse) {
	var count int64

	query = fmt.Sprintf("%%%s%%", query)

	tx := s.db.WithContext(ctx).
		Where("name LIKE ?", query).
		Where("last_name LIKE ?", query).
		Where("phone LIKE ?", query)
	if query != "%%" {
		tx.Where("email LIKE ?", query)
		tx.Where("id_code LIKE ?", query)
		tx.Where("username LIKE ?", query)
	}
	if suspendedAt != "" && suspendedAt != "-1" {
		tx.Where("suspended_at" + database.NullClause(suspendedAt))
	}
	if isOfficial != "" && isOfficial != "-1" {
		tx.Where("made_official_at" + database.NullClause(isOfficial))
	}
	if isProfileComplete != "" && isProfileComplete != "-1" {
		tx.Where("profile_completed_at" + database.NullClause(isProfileComplete))
	}
	err := tx.Model(&models.User{}).Count(&count).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return count, response.GormErrorResponse(err, "Error in calculating the data")
	}

	return count, response.ErrorResponse{}
}

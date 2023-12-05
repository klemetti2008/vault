package endpoints

import (
	"context"
	"fmt"

	"gitag.ir/cookthepot/services/vault/database"
	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"gitag.ir/cookthepot/services/vault/services/role"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) QueryExperts(
	ctx context.Context, query string,
) ([]models.User, response.ErrorResponse) {
	var users []models.User
	query = fmt.Sprintf("%%%s%%", query)
	if !policy.CanQueryUsers(ctx) {
		s.logger.With(ctx).Error("You do not have permission to access this user")
		return users, response.ErrorForbidden("You do not have permission to access this user")
	}
	roleTitles := []string{role.Admin, role.Chef, role.Expert}

	subQuery := s.db.Table("user_role").
		Joins("join roles on roles.id = user_role.role_id").
		Where("roles.title IN ?", roleTitles).
		Select("user_id")

	tx := s.db.WithContext(ctx).
		Where("id IN (?)", subQuery).
		Preload("Roles")
	if query != "%%" {
		tx.Where(
			"LOWER(name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(email) LIKE ? OR LOWER("+
				"id_code) LIKE ? OR LOWER(username) LIKE ?", query, query, query, query, query, query,
		)
	}
	tx.Where("suspended_at" + database.NullClause("0"))
	tx.Where("profile_completed_at" + database.NullClause("1"))
	err := tx.Find(&users).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return users, response.GormErrorResponse(err, "An error occurred while finding users")
	}

	return users, response.ErrorResponse{}
}

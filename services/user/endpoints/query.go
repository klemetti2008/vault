package endpoints

import (
	"context"
	"fmt"
	"strings"

	"gitag.ir/cookthepot/services/vault/database"
	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Query(
	ctx context.Context, offset, limit int, orderBy, order, query,
	suspendedAt, isOfficial, isProfileComplete string,
) ([]models.User, response.ErrorResponse) {
	var users []models.User
	query = fmt.Sprintf("%%%s%%", query)
	query = strings.ToLower(query)

	if !policy.CanQueryUsers(ctx) {
		s.logger.With(ctx).Error("You do not have permission to access this user")
		return users, response.ErrorForbidden("You do not have permission to access this user")
	}
	fmt.Println("query", query)
	tx := s.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", orderBy, order)).
		Offset(offset).Limit(limit).
		Preload("Roles")
	if query != "%%" {
		tx.Where(
			"LOWER(name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(email) LIKE ? OR LOWER("+
				"id_code) LIKE ? OR LOWER(username) LIKE ?", query, query, query, query, query, query,
		)
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
	err := tx.Find(&users).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return users, response.GormErrorResponse(err, "Error in finding users")
	}

	return users, response.ErrorResponse{}
}

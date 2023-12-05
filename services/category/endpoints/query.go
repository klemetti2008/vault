package endpoints

import (
	"context"
	"fmt"
	"strings"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Query(
	ctx context.Context, offset, limit int, orderBy, order, query string,
) ([]models.Category, response.ErrorResponse) {
	var categories []models.Category
	query = fmt.Sprintf("%%%s%%", query)
	query = strings.ToLower(query)
	err := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Where("LOWER(title) LIKE ?", query).
		Order(fmt.Sprintf("%s %s", orderBy, order)).
		Preload("User").
		Preload("Parent.Parent.Parent").
		Preload("Children.Children.Children").
		Find(&categories).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Category{}, response.GormErrorResponse(err, "An error occurred while finding categories")
	}

	return categories, response.ErrorResponse{}
}

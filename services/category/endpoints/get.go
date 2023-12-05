package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) Get(ctx context.Context, id string) (models.Category, response.ErrorResponse) {
	var category models.Category
	err := s.db.WithContext(ctx).
		Preload("User").
		Preload("Parent.Parent.Parent").
		Preload("Children.Children.Children").
		First(&category, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Category{}, response.GormErrorResponse(err, "Error in fetching difficulty level")
	}

	return category, response.ErrorResponse{}
}

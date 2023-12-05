package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/response"
)

type UpdateCategoryRequest struct {
	ParentID int    `json:"parent_id"`
	Title    string `json:"title"`
}

func (c *UpdateCategoryRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"parent_id": govalidity.New("parent_id").Optional(),
		"title":     govalidity.New("title").Required().MinMaxLength(2, 200),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"parent_id": "Related Category",
			"title":     "Title",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateCategoryRequest) (
	models.Category, response.ErrorResponse,
) {
	var category models.Category

	if !policy.CanUpdateCategory(ctx) {
		s.logger.With(ctx).Error("You do not have permission to update a category")
		return models.Category{}, response.ErrorForbidden(nil, "You do not have permission to update a category")
	}

	err := s.db.WithContext(ctx).First(&category, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Category{}, response.GormErrorResponse(err, "An error occurred while finding the category")
	}
	category.ParentID = dtp.NullInt64{
		Int64: int64(input.ParentID),
		Valid: input.ParentID != 0,
	}
	category.Title = input.Title

	err = s.db.WithContext(ctx).Save(&category).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Category{}, response.GormErrorResponse(err, "An error occurred while updating the category")
	}
	return category, response.ErrorResponse{}
}

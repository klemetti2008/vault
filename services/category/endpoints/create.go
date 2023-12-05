package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/response"
)

type CreateCategoryRequest struct {
	ParentID int    `json:"parent_id"`
	Title    string `json:"title"`
}

func (c *CreateCategoryRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
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

func (s *service) Create(ctx context.Context, input CreateCategoryRequest) (models.Category, response.ErrorResponse) {
	if !policy.CanCreateCategory(ctx) {
		s.logger.With(ctx).Error("You do not have permission to create a category")
		return models.Category{}, response.ErrorForbidden("You do not have permission to create a category")
	}

	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	category := models.Category{
		ParentID: dtp.NullInt64{
			Int64: int64(input.ParentID),
			Valid: input.ParentID != 0,
		},
		UserID: id,
		Title:  input.Title,
	}
	err := s.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return category, response.GormErrorResponse(err, "An error occurred while creating the category")
	}
	return category, response.ErrorResponse{}
}

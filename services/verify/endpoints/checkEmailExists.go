package endpoints

import (
	"context"
	"errors"
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
	"gorm.io/gorm"
)

type CheckEmailExistsRequest struct {
	Email string `json:"email"`
}

func (r *CheckEmailExistsRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"email": govalidity.New("email").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"email": "Email",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) CheckEmailExists(ctx context.Context, input CheckEmailExistsRequest) (bool, response.ErrorResponse) {
	var user models.User
	var count int64
	err := s.db.WithContext(ctx).First(&user, "email", input.Email).Count(&count).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.With(ctx).Error(err)
		return false, response.GormErrorResponse(err, "An error occurred while finding the email")
	}
	exists := count > 0
	return exists, response.ErrorResponse{}
}

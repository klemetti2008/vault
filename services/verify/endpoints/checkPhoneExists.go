package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
)

type CheckPhoneExistsRequest struct {
	Phone string `json:"phone,omitempty"`
}

func (r *CheckPhoneExistsRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone": govalidity.New("phone").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"phone": "Phone",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) CheckPhoneExists(ctx context.Context, input CheckPhoneExistsRequest) (bool, response.ErrorResponse) {
	var user models.User
	var count int64
	err := s.db.WithContext(ctx).First(&user, "phone", input.Phone).Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return false, response.GormErrorResponse(err, "An error occurred while finding the phone number")
	}
	exists := count > 0
	return exists, response.ErrorResponse{}
}

package endpoints

import (
	"context"
	"net/http"
	"time"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
)

type ExchangeRequest struct {
	Code string `json:"code,omitempty"`
}

func (r *ExchangeRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"code": govalidity.New("code").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"code": "Code",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) Exchange(ctx context.Context, input ExchangeRequest) (string, response.ErrorResponse) {
	var verification models.Verification
	err := s.db.WithContext(ctx).Find(&verification, "code", input.Code).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "An error occurred while finding the code")
	}
	if verification.Expired() {
		return "", response.GormErrorResponse(err, "The code has expired.")
	}

	verification.ExpiresAt = time.Now()
	if err = s.db.WithContext(ctx).Save(&verification).Error; err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "An error occurred while saving the code")
	}

	sessionCode := verification.SessionCode
	return sessionCode, response.ErrorResponse{}
}

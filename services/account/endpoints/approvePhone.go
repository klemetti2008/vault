package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/faker"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) ApprovePhone(ctx context.Context, code string) (string, response.ErrorResponse) {
	var user models.User
	var phone string
	var ok string

	phone, responseError := s.checkAndDeleteVerificationByCode(ctx, code)
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(responseError)
		return ok, responseError
	}

	_, user, err := s.findUser(ctx, phone)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return ok, response.GormErrorResponse(err, "An error occurred in finding the mobile number")
	}

	user.PhoneVerifiedAt = faker.SQLNow()
	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return ok, response.GormErrorResponse(err, "An error occurred in saving the mobile number")
	}

	ok = "phone number verified successfully"

	return ok, response.ErrorResponse{}
}

package endpoints

import (
	"context"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/response"
)

// rethink about getting the phone from the user of from the verification field it's self
func (s *service) checkAndDeleteVerificationBySessionCodeAndPhone(
	ctx context.Context, sessionCode string, phone string,
) response.ErrorResponse {
	var verification models.Verification
	err := s.db.WithContext(ctx).First(&verification, "session_code", sessionCode).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return response.GormErrorResponse(err, "validation error")
	}
	if verification.Expired() {
		s.logger.Info(verification.NotValidPhone(phone), verification.Expired())
		return response.ErrorBadRequest(err, "the verification code has expired. please try again")
	}
	err = s.db.WithContext(ctx).Delete(&models.Verification{}, verification.ID).Error
	return response.GormErrorResponse(err, "error in validation")
}

func (s *service) checkAndDeleteVerificationBySessionCodeAndEmail(
	ctx context.Context, sessionCode string, email string,
) response.ErrorResponse {
	var verification models.Verification
	err := s.db.WithContext(ctx).First(&verification, "session_code", sessionCode).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return response.GormErrorResponse(err, "validation error")
	}
	if verification.Expired() {
		// s.logger.Info(verification.NotValidPhone(phone), verification.Expired())
		s.logger.Info(verification.Expired())
		return response.ErrorBadRequest(err, "the verification code has expired. please try again")
	}
	err = s.db.WithContext(ctx).Delete(&models.Verification{}, verification.ID).Error
	return response.GormErrorResponse(err, "Error in validation")

}

func (s *service) checkAndDeleteVerificationByCode(ctx context.Context, code string) (string, response.ErrorResponse) {
	var phone string
	var verification models.Verification
	err := s.db.WithContext(ctx).First(&verification, "code", code).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return phone, response.GormErrorResponse(err, "validation error. code not found")
	}
	if verification.Expired() {
		s.logger.Info(verification.Expired())
		return phone, response.ErrorBadRequest(err, "the verification code has expired. please try again")
	}

	phone = verification.Phone

	err = s.db.WithContext(ctx).Delete(&models.Verification{}, verification.ID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return phone, response.GormErrorResponse(err, "Validation error")
	}
	return phone, response.ErrorResponse{}
}

package endpoints

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitag.ir/cookthepot/services/vault/config"
	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/modules/codegen"
	"github.com/gofrs/uuid"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/log"
	"github.com/mhosseintaher/kit/response"
	"gorm.io/gorm"
)

type SendCodeRequest struct {
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
}

func (r *SendCodeRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone": govalidity.New("phone").Optional(),
		"email": govalidity.New("email").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"phone": "Phone",
			"email": "Email",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func checkError(ctx context.Context, logger log.Logger, err error) bool {
	if err != nil {
		logger.With(ctx).Error(err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
	}
	return false
}

func (s *service) SendCode(ctx context.Context, input SendCodeRequest) (string, response.ErrorResponse) {
	if input.Phone == "" && input.Email == "" {
		return "", response.ErrorBadRequest(nil, "At least one of the phone or email fields must have a value")
	}

	if config.AppConfig.Environment != "development" {
		if input.Phone != "" {
			var lastPhoneCode models.Verification
			err := s.db.WithContext(ctx).Last(&lastPhoneCode, "phone", input.Phone).Error
			if checkError(ctx, s.logger, err) {
				return "", response.GormErrorResponse(err, "An error occurred while finding the code")
			}

			if !lastPhoneCode.Expired() {
				seconds := int(lastPhoneCode.ExpiresAt.Sub(time.Now()).Seconds())
				er := response.GormErrorResponse(err,
					"Sending the code is not possible. Please wait for "+
						fmt.Sprintf("%d", seconds)+
						" seconds ",
				)
				return "", er
			}
		}

		if input.Email != "" {
			var lastEmailCode models.Verification
			err := s.db.WithContext(ctx).Last(&lastEmailCode, "email", input.Email).Error
			if checkError(ctx, s.logger, err) {
				return "", response.GormErrorResponse(err, "An error occurred while finding the code")
			}

			if !lastEmailCode.Expired() {
				seconds := int(lastEmailCode.ExpiresAt.Sub(time.Now()).Seconds())
				er := response.GormErrorResponse(err,
					"Sending the code is not possible. Please wait for "+
						fmt.Sprintf("%d", seconds)+
						" seconds ",
				)
				return "", er
			}
		}

	}

	uid, e := uuid.NewV4()
	if e != nil {
		return "", response.ErrorInternalServerError(e, "An error occurred while creating the code")
	}

	verification := models.Verification{
		Code:  codegen.GenerateRandomNumber(),
		Phone: input.Phone,
		Email: input.Email,
		// FIXME: make this configurable
		ExpiresAt:   time.Now().Add(time.Minute * 5),
		SessionCode: uid.String(),
	}
	err := s.db.WithContext(ctx).Create(&verification).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "An error occurred while creating the code")
	}
	code := verification.Code

	if config.AppConfig.SendSMS {
		if input.Phone != "" {
			err = s.notifier.SendPhoneVerifyCode(input.Phone, code)
		}
	}

	if config.AppConfig.SendEmail {
		message := "Your verification code is: " + code
		if input.Email != "" {
			err = s.notifier.SendEmailVerifyCode(input.Email, message)
		}
	}

	if config.AppConfig.Environment != "development" {
		code = "successfully sent"
	}
	return code, response.ErrorResponse{}

}

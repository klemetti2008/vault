package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/modules/encrypt"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
)

type ResetPasswordRequest struct {
	Phone       string `json:"phone,omitempty"`
	Password    string `json:"password,omitempty"`
	SessionCode string `json:"session_code,omitempty"`
}

func (r *ResetPasswordRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone":        PhoneNumberValidity,
		"password":     PasswordValidity,
		"session_code": govalidity.New("session_code").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"phone":        "Phone Number",
			"password":     "Password",
			"session_code": "Session Code",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, input ResetPasswordRequest) (LoginResponse, response.ErrorResponse) {
	var user models.User
	var exists bool
	var res LoginResponse
	exists, user, err := s.findUser(ctx, input.Phone)
	if !exists {
		return res, response.ErrorBadRequest(nil, "No user with this username was found")
	}

	responseError := s.checkAndDeleteVerificationBySessionCodeAndPhone(ctx, input.SessionCode, input.Phone)
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(err)
		return res, responseError
	}

	var hashedPassword string
	hashedPassword, err = encrypt.HashPassword(input.Password)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.ErrorInternalServerError(err, "An error occurred on the server")
	}

	user.Password = hashedPassword
	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.GormErrorResponse(err, "An error occurred while generating the password")
	}

	userIDString := strconv.Itoa(int(user.ID))
	_, err = s.DeleteTokens(ctx, []string{userIDString})
	if err != nil {
		s.logger.Error(err)
		return LoginResponse{}, response.ErrorInternalServerError(nil, "An error occurred on the server")

	}

	accessToken, refreshToken, responseError := s.generateTokens(ctx, user, "")
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(err)
		return res, responseError
	}
	res = LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return res, response.ErrorResponse{}
}

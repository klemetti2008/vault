package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/modules/encrypt"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
)

type ChangePasswordRequest struct {
	Password    string `json:"password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

func (r *ChangePasswordRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"password":     govalidity.New("password").Required(),
		"new_password": PasswordValidity,
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"password":     "Password",
			"new_password": "New Password",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) ChangePassword(
	ctx context.Context, input ChangePasswordRequest,
) (
	LoginResponse,
	response.ErrorResponse,
) {
	var user models.User

	var loginResponse LoginResponse

	Id := policy.ExtractIdClaim(ctx)
	err := s.db.WithContext(ctx).First(&user, "id", Id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return loginResponse, response.GormErrorResponse(err, "An error occurred in the database")
	}

	if !user.ValidatePassword(input.Password) {
		return LoginResponse{}, response.ErrorBadRequest(err, "Incorrect password")
	}

	var hashedPassword string
	hashedPassword, err = encrypt.HashPassword(input.NewPassword)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return loginResponse, response.ErrorInternalServerError(err, "An error occurred on the server.")
	}

	user.Password = hashedPassword
	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		err := response.GormErrorResponse(err, "An error occurred while saving.")
		return loginResponse, err
	}

	userIDString := strconv.Itoa(int(user.ID))
	_, err = s.DeleteTokens(ctx, []string{userIDString})
	if err != nil {
		s.logger.Error(err)
		return LoginResponse{}, response.ErrorInternalServerError(nil, "An error occurred while deleting tokens")
	}

	accessTkn, refreshTkn, responseError := s.generateTokens(ctx, user, "")
	if responseError.StatusCode != 0 {
		s.logger.With(ctx).Error(responseError)
		return loginResponse, responseError
	}
	loginResponse = LoginResponse{
		User:         user,
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}
	return loginResponse, response.ErrorResponse{}
}

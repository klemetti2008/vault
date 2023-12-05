package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
)

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"refresh_token": govalidity.New("refresh_token").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"refresh_token": "ًRefresh Token",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) RefreshToken(ctx context.Context, input RefreshTokenRequest) (RefreshTokenResponse, response.ErrorResponse) {
	var token RefreshTokenResponse
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return token, response.GormErrorResponse(err, "An error occurred while creating the transaction")
	}

	var oldToken models.Token
	if err = tx.WithContext(ctx).Preload("User").First(
		&oldToken, "refresh_token",
		input.RefreshToken,
	).Error; err != nil {
		tx.Rollback()
		return token, response.GormErrorResponse(err, "An error occurred while getting the refresh token")
	}

	if !oldToken.IsValidRefreshToken() {
		tx.Rollback()
		s.logger.With(ctx).Error("Token validation error")
		return token, response.GormErrorResponse(err, "Token validation error")
	}

	_, _user, errr := s.findUser(ctx, oldToken.User.Phone)
	if errr != nil {
		tx.Rollback()
		s.logger.With(ctx).Error(errr)
		return token, response.GormErrorResponse(errr, "An error occurred while fetching the user")
	}

	newAccessToken, newRefreshToken, er := s.generateTokens(ctx, _user, oldToken.AccessToken)

	if er.StatusCode != 0 {
		tx.Rollback()
		return token, er
	}

	if err = tx.Commit().Error; err != nil {
		s.logger.With(ctx).Error(err)
		return token, response.GormErrorResponse(err, "An error occurred while saving the transaction")
	}

	return RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, response.ErrorResponse{}
}

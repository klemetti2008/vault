package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
)

type LoginResponse struct {
	User         models.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *LoginRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"username": govalidity.New("username").MinLength(3).MaxLength(128).Required(),
		"password": govalidity.New("password").MinLength(8).MaxLength(128).Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"username": "Username",
			"password": "Password",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)
	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) authenticate(ctx context.Context, username, password string) Identity {
	logger := s.logger.With(ctx, "user", username)
	_, user, err := s.findUser(ctx, username)
	if err != nil {
		logger.Error("authenticate the user has problem in findUser() method", err)
	}
	if user.ValidatePassword(password) {
		logger.Infof("authentication successful")
		return user
	}
	logger.Infof("authentication failed")

	return nil
}

func (s *service) Login(ctx context.Context, input LoginRequest) (LoginResponse, response.ErrorResponse) {
	identity := s.authenticate(ctx, input.Username, input.Password)
	if identity != nil {
		token, refreshToken, errr := s.generateTokens(ctx, identity, "")
		if errr.StatusCode != 0 {
			s.logger.With(ctx).Error(errr)

			return LoginResponse{}, errr
		}

		var user models.User
		_, user, err := s.findUser(ctx, input.Username)
		if err != nil {
			s.logger.With(ctx).Error(err)
			return LoginResponse{}, response.GormErrorResponse(err, "Incorrect username or password")
		}

		loginResponse := LoginResponse{
			User:         user,
			AccessToken:  token,
			RefreshToken: refreshToken,
		}
		return loginResponse, response.ErrorResponse{}
	}
	return LoginResponse{}, response.ErrorBadRequest(nil, "Incorrect username or password")
}

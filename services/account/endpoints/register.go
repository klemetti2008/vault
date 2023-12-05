package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/modules/encrypt"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/faker"
	"github.com/mhosseintaher/kit/response"
)

type RegisterRequest struct {
	Name        string `json:"name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	SessionCode string `json:"session_code,omitempty"`
}

func (r *RegisterRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"email":        govalidity.New("email").Email().Required(),
		"password":     PasswordValidity,
		"name":         govalidity.New("name").MinLength(2).MaxLength(200).Required(),
		"last_name":    govalidity.New("last_name").MinLength(2).MaxLength(200).Required(),
		"session_code": govalidity.New("session_code").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"name":         "Name",
			"last_name":    "Lastname",
			"email":        "Email",
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

func (s *service) Register(ctx context.Context, input RegisterRequest) (LoginResponse, response.ErrorResponse) {
	var user models.User
	var exists bool
	var res LoginResponse

	exists, user, err := s.findUser(ctx, input.Email)
	if exists {
		return res, response.ErrorBadRequest(nil, "A user with this username was found")
	}

	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.GormErrorResponse(err, "An error occurred")
	}

	eerr := s.checkAndDeleteVerificationBySessionCodeAndEmail(ctx, input.SessionCode, input.Email)
	if eerr.StatusCode != 0 {
		s.logger.With(ctx).Error(eerr)
		return res, eerr
	}

	var hashedPassword string
	hashedPassword, err = encrypt.HashPassword(input.Password)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.ErrorInternalServerError(nil, "An error occurred while creating the password")
	}

	user = models.User{
		Name:     input.Name,
		LastName: input.LastName,
		Email: dtp.NullString{
			String: input.Email,
			Valid:  true,
		},
		Password:        hashedPassword,
		PhoneVerifiedAt: faker.SQLNow(),
	}
	err = s.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return res, response.GormErrorResponse(err, "An error occurred while registering the user")
	}

	accessToken, refreshToken, errr := s.generateTokens(ctx, user, "")
	if errr.StatusCode != 0 {
		return res, errr
	}

	res = LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return res, response.ErrorResponse{}
}

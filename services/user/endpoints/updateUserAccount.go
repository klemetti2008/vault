package endpoints

import (
	"context"
	"net/http"
	"time"

	accountentpoints "gitag.ir/cookthepot/services/vault/services/account/endpoints"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/response"
)

type UpdateUserAccountRequest struct {
	Bio                 string    `json:"biography,omitempty"`
	Name                string    `json:"name,omitempty"`
	Email               string    `json:"email,omitempty"`
	IDCode              string    `json:"id_code,omitempty"`
	Phone               string    `json:"phone,omitempty"`
	Username            string    `json:"username,omitempty"`
	LastName            string    `json:"last_name,omitempty"`
	Nickname            string    `json:"nickname,omitempty"`
	AvatarUrl           string    `json:"avatar_url,omitempty"`
	CountryCode         string    `json:"country_code,omitempty"`
	City                string    `json:"city,omitempty"`
	Gender              string    `json:"gender,omitempty"`
	DateOfBirth         time.Time `json:"date_of_birth,omitempty"`
	MadeProfilePublicAt bool      `json:"made_profile_public_at,omitempty"`
}

func (s *service) isEnoughUpdateAccountData(user models.User, input UpdateUserAccountRequest) bool {
	if input.Name != "" && input.LastName != "" && input.Username != "" && input.Phone != "" && input.Email != "" && input.IDCode != "" &&
		input.CountryCode != "" &&
		input.Gender != "" && (input.DateOfBirth != time.Time{}) && user.PhoneVerifiedAt.Valid {
		return true
	}
	return false
}

func (req *UpdateUserAccountRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone":                  accountentpoints.PhoneNumberValidity,
		"biography":              govalidity.New("biography").Optional().MinMaxLength(2, 200),
		"name":                   govalidity.New("name").Required().MinMaxLength(2, 200),
		"email":                  govalidity.New("email").Required(),
		"id_code":                govalidity.New("id_code").Optional(),
		"username":               govalidity.New("username").Required().MinMaxLength(2, 200),
		"last_name":              govalidity.New("last_name").Required().MinMaxLength(2, 200),
		"nickname":               govalidity.New("nickname").MinMaxLength(2, 200).Optional(),
		"avatar_url":             govalidity.New("avatar_url").Optional(),
		"country_code":           govalidity.New("country_code").Required(),
		"city":                   govalidity.New("city").Optional(),
		"gender":                 govalidity.New("gender").Required(),
		"date_of_birth":          govalidity.New("date_of_birth").Required(),
		"made_profile_public_at": govalidity.New("made_profile_public_at").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"biography":              "Biography",
			"name":                   "Name",
			"email":                  "Email",
			"id_code":                "National ID",
			"phone":                  "Phone",
			"username":               "Username",
			"last_name":              "Lastname",
			"nickname":               "Nickname",
			"avatar_url":             "Avatar url",
			"country_code":           "Country Code",
			"city":                   "City",
			"gender":                 "Gender",
			"date_of_birth":          "Birth Date",
			"made_profile_public_at": "Made profile public at",
		},
	)

	errr := govalidity.ValidateBody(r, schema, req)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) UpdateAccount(ctx context.Context, input UpdateUserAccountRequest) (models.User, response.ErrorResponse) {
	var user models.User
	Id := policy.ExtractIdClaim(ctx)
	err := s.db.WithContext(ctx).First(&user, "id", Id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "Error in finding the user")
	}

	if !policy.CanUpdateAccount(ctx, user) {
		return user, response.ErrorForbidden("You do not have permission to access this user")
	}
	phoneVerifiedAt := user.PhoneVerifiedAt
	if input.Phone != user.Phone {
		phoneVerifiedAt = dtp.NullTime{Valid: false, Time: time.Now()}
	}

	emailVerifiedAt := user.EmailVerifiedAt
	if input.Email != user.Email.String {
		emailVerifiedAt = dtp.NullTime{Valid: false, Time: time.Now()}
	}

	profileCompletedAt := dtp.NullTime{Valid: s.isEnoughUpdateAccountData(user, input), Time: time.Now()}

	user.Bio = input.Bio
	user.Name = input.Name
	user.Email = dtp.NullString{String: input.Email, Valid: input.Email != ""}
	user.EmailVerifiedAt = emailVerifiedAt
	user.IDCode = dtp.NullString{String: input.IDCode, Valid: input.IDCode != ""}
	user.Phone = input.Phone
	user.PhoneVerifiedAt = phoneVerifiedAt
	user.Username = dtp.NullString{String: input.Username, Valid: input.Username != ""}
	user.LastName = input.LastName
	user.Nickname = input.Nickname
	user.AvatarUrl = input.AvatarUrl
	user.CountryCode = input.CountryCode
	user.City = input.City
	user.Gender = input.Gender
	user.DateOfBirth = dtp.NullTime{
		Time:  input.DateOfBirth,
		Valid: input.DateOfBirth != time.Time{},
	}
	user.MadeProfilePublicAt = dtp.NullTime{
		Valid: input.MadeProfilePublicAt,
		Time:  time.Now(),
	}
	user.ProfileCompletedAt = profileCompletedAt

	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "Error in saving the user")
	}

	return user, response.ErrorResponse{}
}

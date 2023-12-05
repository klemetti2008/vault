package endpoints

import (
	"net/http"
	"time"

	accountentpoints "gitag.ir/cookthepot/services/vault/services/account/endpoints"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdateUserRequest struct {
	Bio                 string    `json:"biography,omitempty" `
	Name                string    `json:"name,omitempty"`
	Phone               string    `json:"phone,omitempty"`
	Email               string    `json:"email,omitempty"`
	Title               string    `json:"title,omitempty"`
	Grade               int       `json:"grade,omitempty"`
	IDCode              string    `json:"id_code,omitempty"`
	LastName            string    `json:"last_name,omitempty"`
	Username            string    `json:"username,omitempty"`
	Nickname            string    `json:"nickname,omitempty"`
	Password            string    `json:"password,omitempty"`
	AvatarUrl           string    `json:"avatar_url,omitempty"`
	CountryCode         string    `json:"country_code,omitempty"`
	City                string    `json:"city,omitempty"`
	Gender              string    `json:"gender,omitempty"`
	DateOfBirth         time.Time `json:"date_of_birth,omitempty"`
	SuspendedAt         bool      `json:"suspended_at,omitempty"`
	MadeOfficialAt      bool      `json:"made_official_at,omitempty"`
	EmailVerifiedAt     bool      `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt     bool      `json:"phone_verified_at,omitempty"`
	Roles               []int     `json:"roles,omitempty"`
	MadeProfilePublicAt bool      `json:"made_profile_public_at,omitempty"`
}

func (r *UpdateUserRequest) Validate(req *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"phone":     accountentpoints.PhoneNumberValidity,
		"biography": govalidity.New("biography").MinMaxLength(2, 200).Optional(),
		"name":      govalidity.New("name").Required().MinLength(2).MaxLength(200),
		"email":     govalidity.New("email").Email().Required(),
		"title":     govalidity.New("title").Required().MinMaxLength(2, 200),
		"grade":     govalidity.New("grade"),
		"id_code":   govalidity.New("id_code"),
		"last_name": govalidity.New("last_name").Required().MinLength(2).MaxLength(100),
		// TODO: username should include not only number to prevent reserved phone hack
		"username":               govalidity.New("username").MinMaxLength(2, 200).Required(),
		"nickname":               govalidity.New("nickname").MinMaxLength(2, 200).Optional(),
		"password":               govalidity.New("password").Optional(),
		"avatar_url":             govalidity.New("avatar_url").Optional(),
		"country_code":           govalidity.New("country_code").Required(),
		"city":                   govalidity.New("city").Optional(),
		"roles":                  govalidity.New("roles").Optional(),
		"gender":                 govalidity.New("gender").Required(),
		"date_of_birth":          govalidity.New("date_of_birth").Optional(),
		"suspended_at":           govalidity.New("suspended_at").Optional(),
		"made_official_at":       govalidity.New("made_official_at").Optional(),
		"email_verified_at":      govalidity.New("email_verified_at").Optional(),
		"phone_verified_at":      govalidity.New("phone_verified_at").Optional(),
		"made_profile_public_at": govalidity.New("made_profile_public_at").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"biography":              "Biography",
			"name":                   "Name",
			"roles":                  "Roles",
			"phone":                  "Phone",
			"email":                  "Email",
			"title":                  "Title",
			"grade":                  "Grade",
			"id_code":                "ID Code",
			"last_name":              "Lastname",
			"username":               "Username",
			"nickname":               "Nickname",
			"password":               "Password",
			"avatar_url":             "Avatar url",
			"country_code":           "Country Code",
			"city":                   "City",
			"gender":                 "Gender",
			"date_of_birth":          "Birth Date",
			"suspended_at":           "Suspended at",
			"made_official_at":       "Made Official at",
			"email_verified_at":      "Email Verified at",
			"phone_verified_at":      "Phone Verified at",
			"made_profile_public_at": "Made Profile Public at",
		},
	)

	errr := govalidity.ValidateBody(req, schema, r)

	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

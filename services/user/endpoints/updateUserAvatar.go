package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
	"github.com/mhosseintaher/kit/response"
)

type UpdateUserAvatarRequest struct {
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (req *UpdateUserAvatarRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"avatar_url": govalidity.New("avatar_url").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"avatar_url": "Avatar url",
		},
	)

	errr := govalidity.ValidateBody(r, schema, req)
	if len(errr) > 0 {
		return govalidity.DumpErrors(errr)
	}

	return nil
}

func (s *service) UpdateAvatar(ctx context.Context, input UpdateUserAvatarRequest) (models.User, response.ErrorResponse) {
	var user models.User
	Id := policy.ExtractIdClaim(ctx)
	err := s.db.WithContext(ctx).First(&user, "id", Id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "Error in finding the user")
	}

	if !policy.CanUpdateAvatar(ctx, user) {
		s.logger.With(ctx).Error("You do not have permission to access this user")
		return user, response.ErrorForbidden("You do not have permission to access this user")
	}
	user.AvatarUrl = input.AvatarUrl

	err = s.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "Error in saving the user")
	}

	return user, response.ErrorResponse{}
}

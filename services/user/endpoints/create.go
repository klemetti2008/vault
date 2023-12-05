package endpoints

import (
	"context"
	"time"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/modules/encrypt"
	"gitag.ir/cookthepot/services/vault/policy"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) isEnoughCreateData(input CreateUserRequest) bool {
	if input.Name != "" && input.LastName != "" && input.Username != "" && input.Phone != "" && input.Email != "" && input.IDCode != "" &&
		input.CountryCode != "" &&
		input.Gender != "" && (input.DateOfBirth != time.Time{}) && input.PhoneVerifiedAt {
		return true
	}
	return false
}

func (s *service) Create(ctx context.Context, input CreateUserRequest) (models.User, response.ErrorResponse) {

	if !policy.CanCreateUser(ctx) {
		s.logger.With(ctx).Error("You do not have permission to access this user")
		return models.User{}, response.ErrorForbidden(nil, "You do not have permission to access this user")
	}

	profileCompletedAt := dtp.NullTime{Valid: s.isEnoughCreateData(input), Time: time.Now()}

	var hashedPassword string
	hashedPassword, err := encrypt.HashPassword(input.Password)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.ErrorInternalServerError(nil, "An error occurred on the server")
	}

	var roles []*models.Role
	if len(input.Roles) != 0 {
		err := s.db.Find(&roles, input.Roles).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.ErrorBadRequest(nil, "Invalid roles sent")
		}
	}

	user := models.User{
		Name:        input.Name,
		Phone:       input.Phone,
		Title:       input.Title,
		LastName:    input.LastName,
		Password:    hashedPassword,
		CountryCode: input.CountryCode,
		City:        input.City,
		Roles:       roles,
		Gender:      input.Gender,
		DateOfBirth: dtp.NullTime{
			Time:  input.DateOfBirth,
			Valid: input.DateOfBirth != time.Time{},
		},
		Bio:       input.Bio,
		Grade:     input.Grade,
		Nickname:  input.Nickname,
		AvatarUrl: input.AvatarUrl,
		Username: dtp.NullString{
			String: input.Username,
			Valid:  input.Username != "",
		},
		Email: dtp.NullString{
			String: input.Email,
			Valid:  input.Email != "",
		},
		IDCode: dtp.NullString{
			String: input.IDCode,
			Valid:  input.IDCode != "",
		},
		SuspendedAt: dtp.NullTime{
			Valid: input.SuspendedAt,
			Time:  time.Now(),
		},
		MadeOfficialAt: dtp.NullTime{
			Valid: input.MadeOfficialAt,
			Time:  time.Now(),
		},
		EmailVerifiedAt: dtp.NullTime{
			Valid: input.EmailVerifiedAt,
			Time:  time.Now(),
		},
		PhoneVerifiedAt: dtp.NullTime{
			Valid: input.PhoneVerifiedAt,
			Time:  time.Now(),
		},
		MadeProfilePublicAt: dtp.NullTime{
			Valid: input.MadeProfilePublicAt,
			Time:  time.Now(),
		},
		ProfileCompletedAt: profileCompletedAt,
	}
	err = s.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "An error occurred while registering the user")
	}
	return user, response.ErrorResponse{}
}

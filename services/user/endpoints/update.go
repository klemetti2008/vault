package endpoints

import (
	"context"
	"fmt"
	"time"

	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/modules/encrypt"
	"gitag.ir/cookthepot/services/vault/policy"
	"gitag.ir/cookthepot/services/vault/services/role"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/response"
	"gorm.io/gorm"
)

func (s *service) isEnoughUpdateData(input UpdateUserRequest) bool {
	if input.Name != "" && input.LastName != "" && input.Username != "" && input.Phone != "" && input.Email != "" && input.IDCode != "" &&
		input.CountryCode != "" &&
		input.Gender != "" && (input.DateOfBirth != time.Time{}) && input.PhoneVerifiedAt {
		return true
	}
	return false
}

func (s *service) Update(ctx context.Context, id string, input UpdateUserRequest) (models.User, response.ErrorResponse) {
	var user models.User
	tx := s.db.WithContext(ctx).Begin() // Begin a new transaction
	if tx.Error != nil {
		return models.User{}, response.GormErrorResponse(tx.Error, "An error occurred on the server")
	}

	if !policy.CanUpdateUser(ctx) {
		s.logger.With(ctx).Error("You do not have permission to access this user")
		return models.User{}, response.ErrorForbidden("You do not have permission to access this user")
	}

	err := tx.Preload("Roles").First(&user, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "Error in finding the user")
	}

	profileCompletedAt := dtp.NullTime{Valid: s.isEnoughUpdateData(input), Time: time.Now()}

	if input.Password != "" {
		var hashedPassword string
		hashedPassword, err = encrypt.HashPassword(input.Password)
		if err != nil {
			return models.User{}, response.ErrorInternalServerError("An error occurred on the server")
		}
		user.Password = hashedPassword
	}

	//check the unique fields not to be duplicated and if user object found and the user id is not the same with the input id return/proper/error
	if input.Email != "" {
		var tmpUser models.User
		err = tx.First(&tmpUser, "email = ? AND id != ?", input.Email, id).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "Error in finding the user")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "The entered email is duplicated")
		}
	}

	if input.Username != "" {
		var tmpUser models.User
		err = tx.First(
			&tmpUser, "username = ? AND id != ?", input.Username, id,
		).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "Error in finding the user")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "The entered username is duplicated")
		}

	}

	if input.Phone != "" {
		var tmpUser models.User
		err = tx.First(&tmpUser, "phone = ? AND id != ?", input.Phone, id).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "Error in finding the user")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "The entered phone number is duplicated")
		}
	}

	if input.IDCode != "" {
		var tmpUser models.User
		err = tx.First(&tmpUser, "id_code = ? AND id != ?", input.IDCode, id).Error
		if err != gorm.ErrRecordNotFound && err != nil {
			s.logger.With(ctx).Error(err)
			return models.User{}, response.GormErrorResponse(err, "Error in finding the user")
		}
		if err == nil && tmpUser.ID != user.ID {
			return models.User{}, response.GormErrorResponse(err, "The entered national ID code is duplicated")
		}
	}

	user.Bio = input.Bio
	user.Name = input.Name
	user.Phone = input.Phone
	user.Grade = input.Grade
	user.Title = input.Title
	user.Email = dtp.NullString{
		String: input.Email,
		Valid:  input.Email != "",
	}
	user.IDCode = dtp.NullString{
		String: input.IDCode,
		Valid:  input.IDCode != "",
	}
	user.LastName = input.LastName
	user.Nickname = input.Nickname
	user.Username = dtp.NullString{
		String: input.Username,
		Valid:  input.Username != "",
	}
	user.AvatarUrl = input.AvatarUrl
	user.CountryCode = input.CountryCode
	user.City = input.City
	user.Gender = input.Gender
	user.DateOfBirth = dtp.NullTime{
		Time:  input.DateOfBirth,
		Valid: input.DateOfBirth != time.Time{},
	}
	user.SuspendedAt = dtp.NullTime{
		Valid: input.SuspendedAt,
		Time:  time.Now(),
	}
	user.MadeOfficialAt = dtp.NullTime{
		Valid: input.MadeOfficialAt,
		Time:  time.Now(),
	}
	user.EmailVerifiedAt = dtp.NullTime{
		Valid: input.EmailVerifiedAt,
		Time:  time.Now(),
	}
	user.PhoneVerifiedAt = dtp.NullTime{
		Valid: input.PhoneVerifiedAt,
		Time:  time.Now(),
	}
	user.MadeProfilePublicAt = dtp.NullTime{
		Valid: input.MadeProfilePublicAt,
		Time:  time.Now(),
	}
	user.ProfileCompletedAt = profileCompletedAt

	var roles []*models.Role
	if len(input.Roles) != 0 {
		err = s.db.Find(&roles, input.Roles).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return user, response.GormErrorResponse(err, "Error in finding roles")
		}
	}

	Id := policy.ExtractIdClaim(ctx)
	userID := fmt.Sprintf("%v", user.ID)
	// if user is admin and wants to delete admin role, return error
	// first find admin role in roles
	var adminRole *models.Role
	for _, r := range roles {
		if r.Title == role.Admin {
			adminRole = r
			break
		}
	}

	for _, r := range user.Roles {
		if r.Title == role.Admin && adminRole == nil && userID == Id {
			err := response.ErrorForbidden("You do not have permission to remove the admin role")
			return models.User{}, err
		}
	}
	err = tx.Model(&user).Association("Roles").Replace(roles)
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "Error in saving the user")
	}
	err = tx.Save(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.User{}, response.GormErrorResponse(err, "Error in saving the user")
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback() // Rollback the transaction if there's any error or panic
		} else {
			err = tx.Commit().Error // Commit the transaction if there are no errors
		}
	}()
	return user, response.ErrorResponse{}
}

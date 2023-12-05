package endpoints

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitag.ir/cookthepot/services/vault/config"
	"gitag.ir/cookthepot/services/vault/models"
	"github.com/golang-jwt/jwt"
	"github.com/mhosseintaher/kit/response"
)

func (s *service) generateAccessTokenJWT(identity Identity) (string, error) {
	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"Audience":  config.AppConfig.AppUrl,
			"ExpiresAt": time.Now().Add(time.Duration(s.tokenExpiration) * time.Minute).Unix(),
			"Id":        identity.GetID(),
			"IssuedAt":  time.Now().Unix(),
			"Issuer":    identity.GetFullName(),
			"NotBefore": time.Now().Unix(),
			"Subject":   identity.GetPhone(),
			"Roles":     identity.GetRoles(),
		},
	).SignedString([]byte(s.accessTokenSigningKey))
	if err != nil {
		s.logger.Error("sign access token", err)
		return "", err
	}
	return token, err
}

func (s *service) generateRefreshTokenJWT(identity Identity) (string, error) {
	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"Audience":  config.AppConfig.AppUrl,
			"ExpiresAt": time.Now().Add(time.Duration(s.refreshTokenExpiration) * time.Minute).Unix(),
			"Id":        identity.GetID(),
			"IssuedAt":  time.Now().Unix(),
			"Issuer":    identity.GetFullName(),
			"NotBefore": time.Now().Unix(),
			"Subject":   identity.GetPhone(),
			"Roles":     identity.GetRoles(),
		},
	).SignedString([]byte(s.refreshTokenSigningKey))
	if err != nil {
		s.logger.Error("sign refresh token", err)
		return "", err
	}
	return refreshToken, err
}

func (s *service) canCreateToken(ctx context.Context, userID uint, offset int) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&models.Token{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		s.logger.Error("count tokens", err)
		return false, err
	}
	return count < int64(config.AppConfig.MaxLoginDeviceCount)+int64(offset), err
}

func (s *service) generateTokens(ctx context.Context, identity Identity, removeCurrentAccessToken string) (
	string, string, response.ErrorResponse,
) {
	var accessToken string
	var refreshToken string
	var err error
	var user models.User
	_, user, err = s.findUser(ctx, identity.GetPhone())
	if err != nil {
		s.logger.With(ctx).Error(err)
		return accessToken, refreshToken, response.GormErrorResponse(err, "An error occurred on the server")
	}
	if user.SuspendedAt.Valid {
		return accessToken, refreshToken, response.ErrorBadRequest(nil, "Your user account has been suspended")
	}

	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		s.logger.With(ctx).Error(tx.Error)
		return accessToken, refreshToken, response.ErrorInternalServerError(nil, "An error occurred on the server")
	}

	// TODO: later we can get device information to store login history
	userID, er := strconv.Atoi(identity.GetID())
	if er != nil {
		s.logger.Error("convert string to int", er)
		tx.Rollback()
		return "", "", response.ErrorInternalServerError(err, "an error occurred on the server")
	}

	offset := 0
	if removeCurrentAccessToken != "" {
		offset = 1
	}

	canCreateToken, er := s.canCreateToken(ctx, uint(userID), offset)
	if er != nil {
		s.logger.Error("can create token", er)
		tx.Rollback()
		return "", "", response.GormErrorResponse(er, "An error occurred on the server")
	}

	if !canCreateToken {
		if config.AppConfig.AutoDeleteDevice {
			var token models.Token
			err = tx.WithContext(ctx).Where("user_id = ?", userID).Order("created_at asc").First(&token).Error
			if err != nil {
				s.logger.Error("find token", err)
				tx.Rollback()
				return "", "", response.GormErrorResponse(err, "An error occurred on the server")
			}

			err = tx.WithContext(ctx).Delete(&token).Error
			if err != nil {
				s.logger.Error("delete token", err)
				tx.Rollback()
				return "", "", response.GormErrorResponse(err, "An error occurred on the server")
			}
		} else {
			tx.Rollback()
			//	get all tokens from db
			var tokens []models.Token
			s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens)
			return "", "", response.ErrorBadRequest(tokens, "The number of authorized devices for you has reached its limit. Please delete one of the previous devices.")
		}
	}

	accessToken, err = s.generateAccessTokenJWT(identity)
	if err != nil {
		s.logger.Error("generate access token", err)
		tx.Rollback()
		return "", "", response.ErrorInternalServerError(err, "an error occurred on the server")
	}

	refreshToken, err = s.generateRefreshTokenJWT(identity)
	if err != nil {
		s.logger.Error("generate refresh token", err)
		tx.Rollback()
		return "", "", response.ErrorInternalServerError(err, "an error occurred on the server")
	}

	token := models.Token{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err = tx.WithContext(ctx).Create(&token).Error; err != nil {
		s.logger.Error("create token", err)
		tx.Rollback()
		return "", "", response.GormErrorResponse(err, "An error occurred while creating the token")

	}

	var currentToken models.Token
	if removeCurrentAccessToken != "" {
		if err = tx.WithContext(ctx).Where("access_token = ?", removeCurrentAccessToken).First(&currentToken).
			Error; err != nil {
			s.logger.Error("find current token", err)
			tx.Rollback()
			return "", "", response.GormErrorResponse(err, "an error occurred while generating the token")
		}
		fmt.Println("current token", currentToken)

		// delete current token
		if err = tx.WithContext(ctx).Delete(&currentToken).Error; err != nil {
			s.logger.Error("delete current token", err)
			tx.Rollback()
			return "", "", response.GormErrorResponse(err, "an error occurred while generating the token")
		}
	}

	if err = tx.Commit().Error; err != nil {
		s.logger.Error("commit transaction", err)
		return "", "", response.GormErrorResponse(err, "an error occurred while completing token creation")
	}

	return accessToken, refreshToken, response.ErrorResponse{}
}

package transports

import (
	"net/http"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/labstack/echo/v4"
	"github.com/mhosseintaher/kit/response"
)

type UserInfoResponse struct {
	User models.User `json:"user"`
}

func (r resource) userInfo(ctx echo.Context) error {
	accessToken := ctx.Request().Header.Get("Authorization")
	accessToken = accessToken[len("Bearer "):]

	user, err := r.service.UserInfo(ctx.Request().Context(), accessToken)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	result := UserInfoResponse{
		User: user,
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}

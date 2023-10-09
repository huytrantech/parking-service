package middleware

import (
	"github.com/labstack/echo/v4"
	"parking-service/core"
	"parking-service/core/utils/jwt_utils"
	"parking-service/model/api_model"
)

type IExternalMiddleware interface {
	AuthTokenPublic(next echo.HandlerFunc) echo.HandlerFunc
}

type externalMiddleware struct {
}

func NewExternalMiddleware() IExternalMiddleware {
	return &externalMiddleware{}
}
func (middleware *externalMiddleware) AuthTokenPublic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		claimToken, err := jwt_utils.VerifyJWTTokenWithExpiredRequest(token)
		if err != nil {
			return api_model.FailResponse(c, err)
		}
		if claimToken != "123456789" {
			return api_model.FailResponse(c, core.NewUnAuthorization())
		}
		return next(c)
	}
}

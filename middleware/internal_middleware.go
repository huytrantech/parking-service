package middleware

import (
	"github.com/labstack/echo/v4"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/utils/jwt_utils"
	"parking-service/model/api_model"
)

type IInternalMiddleware interface {
	AuthTokenInside(next echo.HandlerFunc) echo.HandlerFunc
}

type internalMiddleware struct {
}

func NewInternalMiddleware() IInternalMiddleware {
	return &internalMiddleware{}
}

func (middleware *internalMiddleware) AuthTokenInside(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if len(token) == 0 {
			return api_model.FailResponse(c, core.NewUnAuthorization())
		}
		claims, err := jwt_utils.VerifyInternalJWTToken(token)
		if err != nil {
			return api_model.FailResponse(c, err)
		}

		c.Set(constants.HEADER_PHONE, claims.Phone)
		c.Set(constants.HEADER_USERNAME, claims.Username)
		return next(c)
	}
}

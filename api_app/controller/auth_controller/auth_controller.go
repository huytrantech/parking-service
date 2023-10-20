package auth_controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"parking-service/core"
	"parking-service/core/utils/jwt_utils"
	"parking-service/model/api_model"
	"parking-service/model/proxy_model/account"
	"parking-service/proxy/account_proxy"
)

type IAuthController interface {
	Login(c echo.Context) error
}

type authController struct {
	IAccountProxy account_proxy.IAccountProxy
}

func NewAuthController(IAccountProxy account_proxy.IAccountProxy) IAuthController {
	return &authController{IAccountProxy: IAccountProxy}
}

// Login godoc
// @Summary Login internal
// @Description  Login internal
// @Tags auth
// @Accept json
// @Produce json
// @param model body account.LoginAccountRequest true "model"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=account.LoginAccountResponse}
// @Router /api/internal/v1/auth/login [post]
func (ctr *authController) Login(c echo.Context) error {
	var request account.LoginAccountRequest
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	accountLoginResp, err := ctr.IAccountProxy.Login(context.Background(), request)
	if err != nil {
		return api_model.FailResponse(c, err)
	}
	if accountLoginResp == nil {
		return api_model.FailResponse(c, core.NewUnAuthorization())
	}
	jwtToken, err := jwt_utils.GenInternalJWTToken(accountLoginResp.Token, accountLoginResp.Phone, accountLoginResp.Username)
	if err != nil {
		return api_model.FailResponse(c, err)
	}
	accountLoginResp.Token = jwtToken
	return api_model.SuccessResponse(c, accountLoginResp)
}

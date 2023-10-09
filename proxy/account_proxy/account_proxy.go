package account_proxy

import (
	"context"
	"parking-service/core"
	"parking-service/core/utils"
	"parking-service/model/proxy_model/account"
)

type IAccountProxy interface {
	Login(ctx context.Context, request account.LoginAccountRequest) (
		*account.LoginAccountResponse , error)
}

type accountProxy struct {

}

func NewAccountProxy() IAccountProxy {
	return &accountProxy{}
}

func (proxy *accountProxy) Login(ctx context.Context, request account.LoginAccountRequest) (
	*account.LoginAccountResponse , error) {
	var accountToken string
	var username string
	if !(request.Phone == "0946515847" && request.Password == "12345678") &&
		!(request.Phone == "0906880616" && request.Password == "12345678") &&
		!(request.Phone == "0938638480" && request.Password == "12345678"){
		return nil , core.NewUnAuthorization()
	}
	switch request.Phone {
	case "0946515847":
		username = "tqhuy"
		break
	case "0906880616":
		username = "vxphuong"
		break
	case "0938638480":
		username = "dnhai"
		break
	}
	accountToken =  utils.GenInternalAccountToken(request.Phone , request.Password)
	if len(accountToken) == 0 {
		return nil , core.NewUnAuthorization()
	}
	var response account.LoginAccountResponse
	response.Token = accountToken
	response.Phone = request.Phone
	response.Username = username
	return &response , nil
}

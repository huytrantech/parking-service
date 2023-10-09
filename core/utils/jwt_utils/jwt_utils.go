package jwt_utils

import (
	"github.com/dgrijalva/jwt-go"
	"parking-service/constants"
	"parking-service/core"
	"time"
)
type Claims struct {
	Token                   string `json:"token"`
	TokenRequestExpiredTime int64  `json:"token_request_expired_time"`
	jwt.StandardClaims
}

type InternalClaims struct{
	Token string `json:"token"`
	Phone string `json:"phone"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenJWTToken(token string) (tokenJwt string, err error) {

	claim := Claims{
		Token:                   token,
		TokenRequestExpiredTime: time.Now().Add(24 * time.Hour).Unix(),
		StandardClaims:          jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 *time.Hour).Unix(),
		},
	}

	claimToken := jwt.NewWithClaims(jwt.SigningMethodHS256 , claim)

	tokenJwt , err = claimToken.SignedString([]byte(constants.KEY_JWT))
	if err != nil {
		err = core.NewBadRequestErrorMessage(err.Error())
		return
	}
	if len(tokenJwt) == 0 {
		err = core.NewBadRequestErrorMessage("Generate token fail")
	}
	return
}

func GenInternalJWTToken(token string, phone string , username string) (tokenJwt string, err error) {

	claim := InternalClaims{
		Token:                   token,
		Phone: phone,
		Username: username,
		StandardClaims:          jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 *time.Hour).Unix(),
		},
	}

	claimToken := jwt.NewWithClaims(jwt.SigningMethodHS256 , claim)

	tokenJwt , err = claimToken.SignedString([]byte(constants.KEY_JWT_INTERNAL))
	if err != nil {
		err = core.NewBadRequestErrorMessage(err.Error())
		return
	}
	if len(tokenJwt) == 0 {
		err = core.NewBadRequestErrorMessage("Generate token fail")
	}
	return
}

func VerifyJWTTokenWithExpiredRequest(jwtToken string) (token string , err error) {

	claims := Claims{}
	tkn , err := jwt.ParseWithClaims(jwtToken , &claims , func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.KEY_JWT),nil
	})

	if err != nil ||  tkn.Valid == false{
		err = core.NewUnAuthorization()
		return
	}

	if len(claims.Token) == 0 || claims.TokenRequestExpiredTime < time.Now().Unix() {
		err = core.NewUnAuthorization()
		return
	}
	token = claims.Token
	return
}

func VerifyInternalJWTToken(jwtToken string) (claims InternalClaims , err error) {

	tkn , err := jwt.ParseWithClaims(jwtToken , &claims , func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.KEY_JWT_INTERNAL),nil
	})

	if err != nil ||  tkn.Valid == false{
		err = core.NewUnAuthorization()
		return
	}

	if len(claims.Token) == 0{
		err = core.NewUnAuthorization()
		return
	}
	return
}
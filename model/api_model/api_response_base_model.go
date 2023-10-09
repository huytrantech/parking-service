package api_model

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"parking-service/constants"
	"parking-service/core"
)

type BaseResponseAPIModel struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func SuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, BaseResponseAPIModel{
		Code:    http.StatusOK,
		Data:    data,
		Message: "Success"})
}

func UnAuthorizeResponse(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, BaseResponseAPIModel{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized"})
}

type ResponseErrorModel struct {
	typErr    int
	ErrorCode int
	Error     error
}

func BadRequest(err error) ResponseErrorModel {
	return ResponseErrorModel{
		ErrorCode: constants.ERR_400,
		Error:     err,
		typErr:    1,
	}
}

func FailResponse(c echo.Context, errResponse error) error {
	if customError, ok := errResponse.(*core.ParkingError); ok {
		switch customError.TypeError {
		case core.BadRequest:
			return c.JSON(http.StatusBadRequest, BaseResponseAPIModel{
				Code:    http.StatusBadRequest,
				Data:    nil,
				Message: errResponse.Error(),
			})
		case core.InternalError:
			return c.JSON(http.StatusInternalServerError, BaseResponseAPIModel{
				Code:    http.StatusInternalServerError,
				Data:    nil,
				Message: errResponse.Error(),
			})
		case core.UnAuthorization:
			return c.JSON(http.StatusUnauthorized, BaseResponseAPIModel{
				Code:    http.StatusUnauthorized,
				Data:    nil,
				Message: customError.Msg,
			})
		case core.Forbidden:
			return c.JSON(http.StatusForbidden, BaseResponseAPIModel{
				Code:    http.StatusForbidden,
				Data:    nil,
				Message: customError.Msg,
			})
		}
	}
	return c.JSON(http.StatusInternalServerError, BaseResponseAPIModel{
		Code:    http.StatusInternalServerError,
		Data:    nil,
		Message: constants.GetErrorConstant()[0],
	})
}

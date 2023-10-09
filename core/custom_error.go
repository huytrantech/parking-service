package core

const BadRequest = 1
const InternalError = 2
const UnAuthorization = 3
const Forbidden = 4
const UnkKnow = 0

type ParkingError struct {
	TypeError int
	ErrorCode int
	Msg       string
}

func (pe ParkingError) Error() string {
	if len(pe.Msg) > 0 {
		return pe.Msg
	}
	return "<nil>"
}

func NewBadRequestErrorMessage(msg string) error {
	return &ParkingError{
		TypeError: BadRequest,
		ErrorCode: 0,
		Msg: msg,
	}
}

func NewBadRequestError(err error) error {
	return &ParkingError{
		TypeError: BadRequest,
		ErrorCode: 0,
		Msg: err.Error(),
	}
}

func NewInternalError(err error,errCode int) error {
	return &ParkingError{
		TypeError: InternalError,
		ErrorCode: errCode,
		Msg: err.Error(),
	}
}

func NewInternalErrorMessage(msg string,errCode int) error {
	return &ParkingError{
		TypeError: InternalError,
		ErrorCode: errCode,
		Msg: msg,
	}
}

func NewUnAuthorization() error {
	return &ParkingError{
		TypeError: UnAuthorization,
		Msg: "UnAuthorization",

	}
}
func NewForbidden(msg string) error {
	return &ParkingError{
		TypeError: Forbidden,
		Msg: msg,

	}
}
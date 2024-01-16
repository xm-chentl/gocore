package frame

import "fmt"

var (
	ErrParameterFailed       = NewError(601, "handle parameters is failed")
	ErrQueryParameterFailed  = NewError(602, "handle query parameters is failed")
	ErrBodyParameterFailed   = NewError(603, "handle body parameters is failed")
	ErrHeaderParameterFailed = NewError(604, "handle header parameters is failed")
	ErrInjectFailed          = NewError(605, "handle inject failed")
	ErrHandleNotExist        = NewError(606, "handle(route) does not exists")
	ErrInvalidHandle         = NewError(607, "invalid handle")
)

type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

func NewError(code int, format string, args ...interface{}) error {
	return &CustomError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

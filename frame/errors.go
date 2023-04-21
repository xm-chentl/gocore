package frame

import "fmt"

var (
	ErrParameterFailed = NewError(601, "handle parameters is failed")
	ErrInjectFailed    = NewError(602, "handle inject failed")
	ErrHandleNotExist  = NewError(604, "handle(route) does not exists")
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

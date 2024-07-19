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
	ErrDataCreateFailed      = NewError(608, "create data failed")
	ErrDataQueryFailed       = NewError(609, "query data failed")
	ErrDataNotExists         = NewError(610, "data does not exist")
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

func WithError(err error, format string, args ...interface{}) error {
	c, ok := err.(*CustomError)
	if !ok {
		return NewError(600, format, args...)
	}

	c.Message = fmt.Sprintf("%s: %s",
		c.Message,
		fmt.Sprintf(format, args...),
	)

	return c
}

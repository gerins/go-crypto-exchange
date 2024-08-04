package error

import (
	"fmt"
	"net/http"
)

type ServerError struct {
	HTTPCode int    // HTTP Status code
	Code     int    // Internal error code
	Message  string // Internal error message
	RawError error  // RAW Error
}

func (e ServerError) Error() string {
	return fmt.Sprintf("%v, %v", e.Message, e.RawError)
}

func New(rawError error, httpCode int, code int, message string) ServerError {
	return ServerError{
		HTTPCode: httpCode,
		Code:     code,
		Message:  message,
		RawError: rawError,
	}
}

var (
	ErrGeneralDatabaseError = func(err error) ServerError {
		return ServerError{http.StatusInternalServerError, 800, "failed executing query to database", err}
	}
	ErrInvalidUsernameOrPassword = func(err error) ServerError {
		return ServerError{http.StatusBadRequest, 900, "login failed invalid username or password", err}
	}
)

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

func (e ServerError) Unwrap() error {
	return e.RawError
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
	ErrGeneralError = func(err error) ServerError {
		return ServerError{http.StatusInternalServerError, 999, "general error", err}
	}
	ErrUnauthorized = func(err error) ServerError {
		return ServerError{http.StatusUnauthorized, 900, "unauthorized role access", err}
	}
	ErrInvalidUsernameOrPassword = func(err error) ServerError {
		return ServerError{http.StatusBadRequest, 901, "login failed invalid username or password", err}
	}
	ErrUserBlocked = func(err error) ServerError {
		return ServerError{http.StatusUnauthorized, 902, "login failed user blocked", err}
	}
	ErrGeneralDatabaseError = func(err error) ServerError {
		return ServerError{http.StatusInternalServerError, 800, "internal dependencies error", err}
	}
	ErrDataNotFound = func(err error) ServerError {
		return ServerError{http.StatusBadRequest, 700, "data not found", err}
	}
)

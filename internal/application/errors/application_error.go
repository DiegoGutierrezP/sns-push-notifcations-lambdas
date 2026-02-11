package errors

import (
	"errors"
	domainErrors "lmbd-digital-push-notifications/internal/domain/errors"
)

type ApplicationError struct {
	Code       domainErrors.ErrorCode
	StatusCode int    // Código HTTP (400, 404, 500, etc.)
	Message    string // Mensaje para el cliente
	Cause      error
}

func (e *ApplicationError) Error() string {
	return string(e.Code)
}

func NewApplicationError(code domainErrors.ErrorCode, statusCode int, message string, cause ...error) *ApplicationError {
	appErr := &ApplicationError{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
	}

	if len(cause) > 0 && cause[0] != nil {
		appErr.Cause = cause[0]
	}

	return appErr
}

func AsApplicationError(err error) (*ApplicationError, bool) {
	var appErr *ApplicationError
	ok := errors.As(err, &appErr)
	return appErr, ok
}

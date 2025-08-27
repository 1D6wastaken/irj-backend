package catalogs

import (
	"errors"
	"net/http"

	_http "irj/pkg/http"
)

// Common.
var (
	ErrUnexpectedError = errors.New("unexpected error")
	ErrServerTimeout   = errors.New("server timeout")
	ErrInvalidPassword = errors.New("password does not respect the password policy")
	ErrPasswordTooLong = errors.New("password size is greater than 72 bytes")
	ErrInvalidParam    = errors.New("invalid parameter")
)

// DB.
var (
	ErrDBResourceRetrieval = errors.New("unexpected error during resource retrieval")
	ErrDBResourceUpdate    = errors.New("unexpected error during resource update")
	ErrDBResourceCreation  = errors.New("unexpected error during resource creation")
	ErrDBResourceDeletion  = errors.New("unexpected error during resource deletion")
)

// Users.
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrTokenNotFound         = errors.New("token not found")
	ErrUserAlreadyRegistered = errors.New("user already registered")
	ErrEmailAlreadyTaken     = errors.New("email already taken")
	ErrUserAlreadyActive     = errors.New("user already active")
	ErrUserNotActive         = errors.New("user not active")
	ErrMailNotConfirmed      = errors.New("user has not confirmed email")
	ErrUserNotAdmin          = errors.New("user not admin")
	ErrWrongCredentials      = errors.New("invalid email/password combination")
	ErrInvalidOldPassword    = errors.New("invalid old password")
)

//nolint:lll
var errorMapping = map[error]_http.APIError{
	ErrUnexpectedError: {HTTPCode: http.StatusInternalServerError, Message: "Internal server error", Details: ErrUnexpectedError.Error()},
	ErrServerTimeout:   {HTTPCode: http.StatusInternalServerError, Message: "Timeout", Details: ErrServerTimeout.Error()},
	ErrInvalidPassword: {HTTPCode: http.StatusBadRequest, Message: "Password invalid", Details: ErrInvalidPassword.Error()},
	ErrInvalidParam:    {HTTPCode: http.StatusBadRequest, Message: "Invalid parameter", Details: ErrInvalidParam.Error()},
	ErrPasswordTooLong: {HTTPCode: http.StatusBadRequest, Message: "Password invalid", Details: ErrPasswordTooLong.Error()},

	ErrDBResourceRetrieval: {HTTPCode: http.StatusInternalServerError, Message: "Internal server error", Details: ErrDBResourceRetrieval.Error()},
	ErrDBResourceUpdate:    {HTTPCode: http.StatusInternalServerError, Message: "Internal server error", Details: ErrDBResourceUpdate.Error()},
	ErrDBResourceDeletion:  {HTTPCode: http.StatusInternalServerError, Message: "Internal server error", Details: ErrDBResourceDeletion.Error()},
	ErrDBResourceCreation:  {HTTPCode: http.StatusInternalServerError, Message: "Internal server error", Details: ErrDBResourceCreation.Error()},

	ErrUserNotFound:          {HTTPCode: http.StatusNotFound, Message: "Resource not found", Details: ErrUserNotFound.Error()},
	ErrTokenNotFound:         {HTTPCode: http.StatusNotFound, Message: "Resource not found", Details: ErrTokenNotFound.Error()},
	ErrUserAlreadyRegistered: {HTTPCode: http.StatusConflict, Message: "Conflict", Details: ErrUserAlreadyRegistered.Error()},
	ErrEmailAlreadyTaken:     {HTTPCode: http.StatusConflict, Message: "Conflict", Details: ErrEmailAlreadyTaken.Error()},
	ErrUserAlreadyActive:     {HTTPCode: http.StatusConflict, Message: "Conflict", Details: ErrUserAlreadyActive.Error()},
	ErrUserNotActive:         {HTTPCode: http.StatusForbidden, Message: "Resource not found", Details: ErrUserNotActive.Error()},
	ErrMailNotConfirmed:      {HTTPCode: http.StatusUnauthorized, Message: "Resource not found", Details: ErrUserNotActive.Error()},
	ErrUserNotAdmin:          {HTTPCode: http.StatusForbidden, Message: "Resource not found", Details: "You are not able to perform this action"},
	ErrWrongCredentials:      {HTTPCode: http.StatusNotFound, Message: "Resource not found", Details: ErrUserNotFound.Error()},
	ErrInvalidOldPassword:    {HTTPCode: http.StatusForbidden, Message: "Invalid old password", Details: ErrInvalidOldPassword.Error()},
}

func GetError(err error) (int, _http.APIError) {
	if err == nil {
		return http.StatusInternalServerError, errorMapping[ErrUnexpectedError]
	} else if value, ok := errorMapping[err]; ok {
		return value.HTTPCode, value
	}

	return http.StatusInternalServerError, _http.APIError{
		HTTPCode: http.StatusInternalServerError,
		Message:  "Internal server error",
		Details:  err.Error(),
	}
}

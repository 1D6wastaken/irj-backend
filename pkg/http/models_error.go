package http

import (
	"net/http"
)

var (
	ErrBadRequest       = APIError{HTTPCode: http.StatusBadRequest, Message: "", Details: ""}
	ErrInvalidBody      = APIError{HTTPCode: http.StatusBadRequest, Message: "Incoming body failed to be decoded", Details: ""}
	ErrNotFound         = APIError{HTTPCode: http.StatusNotFound, Message: "Resource not found", Details: ""}
	ErrMethodNotAllowed = APIError{HTTPCode: http.StatusMethodNotAllowed, Message: "Method not allowed", Details: ""}
	ErrConflict         = APIError{HTTPCode: http.StatusConflict, Message: "Resource already exists", Details: ""}
	ErrUnauthorized     = APIError{HTTPCode: http.StatusUnauthorized, Message: "Unauthorized request", Details: ""}
	ErrInternalError    = APIError{HTTPCode: http.StatusInternalServerError, Message: "Server panicked unexpectedly", Details: ""}
)

type APIError struct {
	HTTPCode int    `json:"-"`
	Message  string `json:"message"`
	Details  string `json:"details,omitempty"`
}

// Just to avoid linter issues...
func (a APIError) Error() string {
	return a.Message
}

func (a APIError) Err(err error) *APIError {
	if err != nil {
		a.Details = err.Error()
	}

	return &a
}

func (a APIError) WithDetails(details string) *APIError {
	a.Details = details

	return &a
}

func (a APIError) Msg(message string) *APIError {
	a.Message = message

	return &a
}

func (a APIError) Write(w http.ResponseWriter) {
	_ = WriteJSONResponse(w, a.HTTPCode, a)
}

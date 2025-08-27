package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-openapi/strfmt"
)

func BuildBearer(token string) string {
	return "Bearer " + token
}

func CloneRequest(r *http.Request) (*http.Request, error) {
	var (
		body []byte
		err  error
	)

	if r.Body != nil {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	clone := r.Clone(r.Context())
	clone.Body = io.NopCloser(bytes.NewBuffer(body))

	return clone, nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, jsonBody interface{}) *APIError {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if jsonBody != nil {
		if err := json.NewEncoder(w).Encode(jsonBody); err != nil {
			return ErrInternalError.Msg("failed to encode json body").Err(err)
		}
	}

	return nil
}

type GeneratedDTO interface {
	Validate(registry strfmt.Registry) error
	UnmarshalBinary(data []byte) error
}

func DecodeJSONBody[T any](r *http.Request) (T, error) {
	defer func() {
		_ = r.Body.Close()
	}()

	var dto T

	//nolint:gocritic
	return dto, json.NewDecoder(r.Body).Decode(&dto)
}

func DecodeAndValidateJSONBody[T GeneratedDTO](r *http.Request) (T, error) {
	dto, err := DecodeJSONBody[T](r)
	if err != nil {
		return dto, err
	}

	if err := dto.Validate(nil); err != nil {
		return dto, err
	}

	return dto, nil
}

type CustomHandler = func(http.ResponseWriter, *http.Request) *APIError

func HandlerMiddleware(handler CustomHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			err.Write(w)
		}
	}
}

type UnexpectedHTTPResponseError struct {
	StatusCode int
	Body       string
}

func (r UnexpectedHTTPResponseError) Error() string {
	return fmt.Sprintf("received unexpected response with statusCode %d and body %s", r.StatusCode, r.Body)
}

func NewUnexpectedHTTPResponseError(r *http.Response) error {
	body := ""

	if r.Body != nil {
		rawBody, err := io.ReadAll(r.Body)
		if err == nil && len(rawBody) > 0 {
			body = string(rawBody)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	}

	return UnexpectedHTTPResponseError{
		StatusCode: r.StatusCode,
		Body:       body,
	}
}

package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"irj/internal/catalogs"
	"irj/internal/config"
	"irj/internal/jwt"
	_http "irj/pkg/http"

	"github.com/rs/zerolog"
)

func CheckJWT(conf *config.Config) _http.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := zerolog.Ctx(r.Context())

			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader != "" {
				if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
					logger.Warn().Msg("token is not a bearer")

					_ = _http.WriteJSONResponse(w, http.StatusUnauthorized, mapJWTErrorToJSONAPIError(jwt.ErrInvalidFormatJWT))

					return
				}

				authHeader = authHeader[7:]

				token, err := parseJWTFromHeader(authHeader, conf, logger)
				if err != nil {
					_ = _http.WriteJSONResponse(w, http.StatusUnauthorized, mapJWTErrorToJSONAPIError(err))

					return
				}

				r = r.WithContext(context.WithValue(r.Context(), catalogs.AccessToken, token))

				next(w, r)

				return
			}

			_ = _http.WriteJSONResponse(w, http.StatusUnauthorized, mapJWTErrorToJSONAPIError(jwt.ErrMissingJWT))
		}
	}
}

func parseJWTFromHeader(tok string, conf *config.Config, logger *zerolog.Logger) (jwt.SessionInfo, error) {
	logger.Debug().Str("access-token", tok).Msg("jwt from Authorization header")

	token, err := jwt.CheckSignatureAndDecode(tok, conf.Session.JWTSecret)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to decode jwt from Authorization header")

		return jwt.SessionInfo{}, err
	}

	if err = token.CheckExpiration(); err != nil {
		logger.Warn().Err(err).Msg("jwt expired")

		return jwt.SessionInfo{}, err
	}

	if err = token.Validate(); err != nil {
		logger.Warn().Err(err).Msg("jwt is invalid")

		return jwt.SessionInfo{}, err
	}

	return token, nil
}

func mapJWTErrorToJSONAPIError(err error) *_http.APIError {
	switch {
	case errors.Is(err, jwt.ErrMissingJWT):
		return &_http.APIError{
			HTTPCode: http.StatusUnauthorized,
			Message:  "Invalid authorization",
			Details:  err.Error(),
		}
	case errors.Is(err, jwt.ErrInvalidFormatJWT), errors.Is(err, jwt.ErrInvalidGrade):
		return &_http.APIError{
			HTTPCode: http.StatusUnauthorized,
			Message:  "Invalid authorization",
			Details:  err.Error(),
		}
	case errors.Is(err, jwt.ErrExpiredJWT):
		return &_http.APIError{
			HTTPCode: http.StatusUnauthorized,
			Message:  "Invalid authorization",
			Details:  err.Error(),
		}
	default:
		return &_http.APIError{
			HTTPCode: http.StatusInternalServerError,
			Message:  "Internal server error",
			Details:  err.Error(),
		}
	}
}

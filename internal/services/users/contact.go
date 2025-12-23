package users

import (
	"context"
	"net/http"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/api"
	"irj/pkg/collections"
	_http "irj/pkg/http"

	"github.com/rs/zerolog"
)

func (u *UserService) ContactForm(w http.ResponseWriter, r *http.Request) *_http.APIError {
	req, err := _http.DecodeAndValidateJSONBody[*api.ContactMessage](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	_, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	logger := zerolog.Ctx(r.Context())

	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	row, err := u.postgresService.Queries.GetUsersByGrade(ctx, queries.UserGradeADMIN)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get admin user")

		return _http.ErrInternalError.Msg("unable to send the message").Err(err)
	}

	to := collections.Map(row, func(u queries.GetUsersByGradeRow) smtp.EmailPerson {
		return smtp.EmailPerson{
			Email: u.Email,
			Name:  u.Prenom + " " + u.Nom,
		}
	})

	if err := u.smtpService.SendContactEmail(ctx, *req.Subject, *req.Message, *req.Email, to); err != nil {
		logger.Error().Err(err).Msg("failed to send contact email")

		return _http.ErrInternalError.Msg("unable to send the message").Err(err)
	}

	return _http.WriteJSONResponse(w, http.StatusCreated, nil)
}

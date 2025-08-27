package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"irj/internal/catalogs"
	queries "irj/internal/postgres/_generated"
	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (u *UserService) UserEmailConfirm(w http.ResponseWriter, r *http.Request) *_http.APIError {
	id := httprouter.ParamsFromContext(r.Context()).ByName("token")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("token is required")
	}

	err := processEmailConfirmation(r.Context(), u, id)
	if err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type emailConfirmExchangeData struct {
	logger *zerolog.Logger
	err    error
	id     string
}

type emailConfirmState func(ctx context.Context, s *UserService, data *emailConfirmExchangeData) emailConfirmState

func processEmailConfirmation(ctx context.Context, s *UserService, id string) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := emailConfirmExchangeData{
		logger: logger,
		id:     id,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := confirmEmailCheckStatus; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during user login")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

func confirmEmailCheckStatus(ctx context.Context, s *UserService, exData *emailConfirmExchangeData) emailConfirmState {
	row, err := s.postgresService.Queries.GetUserByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.err = catalogs.ErrUserNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get user by id")
		exData.err = catalogs.ErrDBResourceRetrieval

		return nil
	}

	if row.Grade == queries.UserGradePENDING {
		exData.logger.Warn().Msg("user is pending and therefore cannot confirm email")
		exData.err = catalogs.ErrUserNotActive

		return nil
	}

	if row.EmailConfirm {
		return nil
	}

	return confirmEmail
}

func confirmEmail(ctx context.Context, s *UserService, exData *emailConfirmExchangeData) emailConfirmState {
	if err := s.postgresService.Queries.ConfirmEmailUserByID(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.err = catalogs.ErrUserNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to confirm email in db")
		exData.err = catalogs.ErrDBResourceUpdate
	}

	return nil
}

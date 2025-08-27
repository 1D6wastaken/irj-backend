package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (u *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("id is required")
	}

	if err := processDeleteUser(r.Context(), u, id, &token); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type deleteUserExchangeData struct {
	logger *zerolog.Logger
	err    error
	id     string
	token  *jwt.SessionInfo
	user   *queries.GetUserByIDRow
}

type deleteUserState func(ctx context.Context, env *UserService, data *deleteUserExchangeData) deleteUserState

func processDeleteUser(ctx context.Context, s *UserService, id string, token *jwt.SessionInfo) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := deleteUserExchangeData{
		logger: logger,
		id:     id,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := deleteUserCheck; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during pending users retrieval")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

func deleteUserCheck(ctx context.Context, s *UserService, exData *deleteUserExchangeData) deleteUserState {
	if exData.id != exData.token.ID && exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot delete users")
		exData.err = catalogs.ErrUserNotAdmin

		return nil
	}

	user, err := s.postgresService.Queries.GetUserByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found and therefore can not be deleted")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get user by id")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.user = &user

	return deleteUser
}

func deleteUser(ctx context.Context, s *UserService, exData *deleteUserExchangeData) deleteUserState {
	if err := s.postgresService.Queries.DeleteUserByID(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found and therefore can not be deleted")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to delete user")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	return sendUserDeletionMail
}

func sendUserDeletionMail(_ context.Context, env *UserService, exData *deleteUserExchangeData) deleteUserState {
	env.stopper.Hold(1)

	//nolint:contextcheck
	go func(env *UserService, row queries.GetUserByIDRow, byAdmin bool) {
		defer env.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		_ = env.smtpService.SendDeletionMail(ctx, []smtp.EmailPerson{
			{
				Name:  row.Prenom + " " + row.Nom,
				Email: row.Email,
			},
		}, byAdmin)
	}(env, *exData.user, exData.id != exData.token.ID)

	return nil
}

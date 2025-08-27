package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/go-openapi/strfmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (u *UserService) GetUserInfo(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("id is required")
	}

	user, err := processGetUserInfo(r.Context(), u, id, &token)
	if err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	return _http.WriteJSONResponse(w, http.StatusOK, user)
}

type getUserInfoExchangeData struct {
	logger *zerolog.Logger
	result getUserInfoResult
	id     string
	token  *jwt.SessionInfo
}

type getUserInfoResult struct {
	user *api.GetUserInfo
	err  error
}

type getUserInfoState func(ctx context.Context, s *UserService, data *getUserInfoExchangeData) getUserInfoState

func processGetUserInfo(ctx context.Context, s *UserService, id string, token *jwt.SessionInfo) (*api.GetUserInfo, error) {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := getUserInfoExchangeData{
		logger: logger,
		id:     id,
		token:  token,
	}

	respChan := make(chan getUserInfoResult, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := getUserInfoCheck; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.result

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during pending users retrieval")

		return nil, catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp.user, resp.err
	}
}

func getUserInfoCheck(ctx context.Context, s *UserService, exData *getUserInfoExchangeData) getUserInfoState {
	if exData.id != exData.token.ID && exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot retrieve other users info")
		exData.result.err = catalogs.ErrUserNotAdmin

		return nil
	}

	user, err := s.postgresService.Queries.GetUserByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found")
			exData.result.err = catalogs.ErrUserNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get user by id")
		exData.result.err = catalogs.ErrUnexpectedError

		return nil
	}

	var uuid strfmt.UUID

	if err := uuid.UnmarshalText([]byte(exData.id)); err != nil {
		exData.logger.Error().Err(err).Msg("failed to unmarshal ID")
		exData.result.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.result.user = &api.GetUserInfo{
		ID:           &uuid,
		Firstname:    &user.Prenom,
		Lastname:     &user.Nom,
		Mail:         (*strfmt.Email)(&user.Email),
		Phone:        user.Telephone.String,
		Organization: user.Organisation.String,
		Domain:       utils.PtrTo(string(user.Domaine)),
		Motivation:   &user.Motivation.String,
		CreationDate: strfmt.Date(user.DateCreation.Time),
	}

	return nil
}

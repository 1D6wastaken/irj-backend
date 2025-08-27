package users

import (
	"context"
	"net/http"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	"irj/pkg/glog"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"
)

func (u *UserService) GetPendingUsers(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	users, err := processGetPendingUsers(r.Context(), u, &token)
	if err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	return _http.WriteJSONResponse(w, http.StatusOK, users)
}

type getPendingUsersExchangeData struct {
	logger *zerolog.Logger
	result getPendingUsersResult
	token  *jwt.SessionInfo
	rows   []queries.GetUsersByGradeRow
}

type getPendingUsersResult struct {
	users api.GetUsersInfo
	err   error
}

type getPendingUsersState func(ctx context.Context, s *UserService, data *getPendingUsersExchangeData) getPendingUsersState

func processGetPendingUsers(ctx context.Context, s *UserService, token *jwt.SessionInfo) (api.GetUsersInfo, error) {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := getPendingUsersExchangeData{
		logger: logger,
		result: getPendingUsersResult{},
		token:  token,
	}

	respChan := make(chan getPendingUsersResult, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := getPendingUsersIfAdmin; state != nil; {
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
		return resp.users, resp.err
	}
}

func getPendingUsersIfAdmin(ctx context.Context, s *UserService, exData *getPendingUsersExchangeData) getPendingUsersState {
	if exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot retrieve pending users")
		exData.result.err = catalogs.ErrUserNotAdmin

		return nil
	}

	users, err := s.postgresService.Queries.GetUsersByGrade(ctx, queries.UserGradePENDING)
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to get pending users")
		exData.result.err = catalogs.ErrDBResourceRetrieval

		return nil
	}

	exData.rows = users

	return getPendingUsersCreateResponse
}

//nolint:gocritic
func getPendingUsersCreateResponse(_ context.Context, _ *UserService, exData *getPendingUsersExchangeData) getPendingUsersState {
	r := make(api.GetUsersInfo, 0, len(exData.rows))

	for _, row := range exData.rows {
		model, err := convertUserInfoRowToModel(exData.logger, &row)
		if err != nil {
			continue
		}

		r = append(r, &model)
	}

	exData.result.users = r

	return nil
}

func convertUserInfoRowToModel(logger *glog.Logger, row *queries.GetUsersByGradeRow) (api.GetUserInfo, error) {
	var uuid strfmt.UUID

	if err := uuid.UnmarshalText([]byte(row.ID)); err != nil {
		logger.Error().Err(err).Str("userID", row.ID).Msg("failed to unmarshal ID")

		return api.GetUserInfo{}, err
	}

	var email strfmt.Email

	if err := email.UnmarshalText([]byte(row.Email)); err != nil {
		logger.Error().Err(err).Str("userID", row.ID).Msg("failed to unmarshal email")

		return api.GetUserInfo{}, err
	}

	return api.GetUserInfo{
		ID:           &uuid,
		Firstname:    &row.Prenom,
		Lastname:     &row.Nom,
		Mail:         &email,
		Phone:        row.Telephone.String,
		Organization: row.Organisation.String,
		Domain:       utils.PtrTo(string(row.Domaine)),
		Motivation:   &row.Motivation.String,
		CreationDate: strfmt.Date(row.DateCreation.Time),
	}, nil
}

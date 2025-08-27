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
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (u *UserService) ApproveRejectUser(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	req, err := _http.DecodeAndValidateJSONBody[*api.PatchUsersBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("id is required")
	}

	if err := processApproveRejectUser(r.Context(), u, id, req, &token); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type approveRejectExchangeData struct {
	logger *zerolog.Logger
	result error
	id     string
	params *api.PatchUsersBody
	token  *jwt.SessionInfo
	user   *queries.GetUserByIDRow
}

type approveRejectState func(ctx context.Context, s *UserService, data *approveRejectExchangeData) approveRejectState

func processApproveRejectUser(ctx context.Context, s *UserService, id string, req *api.PatchUsersBody, token *jwt.SessionInfo) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := approveRejectExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := approveOrRejectIfAdmin; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.result

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

func approveOrRejectIfAdmin(ctx context.Context, s *UserService, exData *approveRejectExchangeData) approveRejectState {
	if exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot approve or reject users")
		exData.result = catalogs.ErrUserNotAdmin

		return nil
	}

	user, err := s.postgresService.Queries.GetUserByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found and therefore can not be approved or rejected")
			exData.result = catalogs.ErrUserNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get user by id")
		exData.result = catalogs.ErrUnexpectedError

		return nil
	}

	exData.user = &user

	if *exData.params.Action == api.PatchUsersBodyActionActivate {
		return approveOrRejectActivateUser
	}

	return approveOrRejectDeleteUser
}

func approveOrRejectActivateUser(ctx context.Context, s *UserService, exData *approveRejectExchangeData) approveRejectState {
	if err := s.postgresService.Queries.ApproveUserByID(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found and therefore can not be approved")
			exData.result = catalogs.ErrUserNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to approve user")
		exData.result = catalogs.ErrUnexpectedError
	}

	exData.logger.Info().Str("userID", exData.id).Str("adminID", exData.token.ID).Msg("user approved")

	return sendUserActivationMail
}

func approveOrRejectDeleteUser(ctx context.Context, s *UserService, exData *approveRejectExchangeData) approveRejectState {
	if err := s.postgresService.Queries.DeleteUserByID(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to delete user")
		exData.result = catalogs.ErrUnexpectedError

		return nil
	}

	exData.logger.Info().Str("userID", exData.id).Str("adminID", exData.token.ID).Msg("user rejected")

	return sendUserRejectionMail
}

func sendUserActivationMail(_ context.Context, s *UserService, exData *approveRejectExchangeData) approveRejectState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(s *UserService, id string, row queries.GetUserByIDRow) {
		defer s.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		_ = s.smtpService.SendActivationMail(ctx, []smtp.EmailPerson{
			{
				Name:  row.Prenom + " " + row.Nom,
				Email: row.Email,
			},
		}, id)
	}(s, exData.id, *exData.user)

	return nil
}

func sendUserRejectionMail(_ context.Context, s *UserService, exData *approveRejectExchangeData) approveRejectState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(s *UserService, row queries.GetUserByIDRow) {
		defer s.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		_ = s.smtpService.SendRejectionMail(ctx, []smtp.EmailPerson{
			{
				Name:  row.Prenom + " " + row.Nom,
				Email: row.Email,
			},
		})
	}(s, *exData.user)

	return nil
}

package password

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"irj/internal/catalogs"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/api"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/rs/zerolog"
)

func (p *PasswordService) RequestNewPassword(w http.ResponseWriter, r *http.Request) *_http.APIError {
	req, err := _http.DecodeAndValidateJSONBody[*api.PostResetPwd](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	if err := processRequestNewPassword(r.Context(), p, req); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type newPasswordExchangeData struct {
	logger *zerolog.Logger
	err    error
	params *api.PostResetPwd
	uuid   string
	user   *queries.GetUserByEmailRow
}

type newPasswordState func(ctx context.Context, env *PasswordService, data *newPasswordExchangeData) newPasswordState

func processRequestNewPassword(ctx context.Context, env *PasswordService, req *api.PostResetPwd) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := newPasswordExchangeData{
		logger: logger,
		params: req,
	}
	respChan := make(chan error, 1)

	env.stopper.Hold(1)

	go func() {
		defer env.stopper.Release()

		for state := newPasswordCheckEmailExistence; state != nil; {
			state = state(ctx, env, &exData)
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

func newPasswordCheckEmailExistence(ctx context.Context, env *PasswordService, exData *newPasswordExchangeData) newPasswordState {
	row, err := env.postgresService.Queries.GetUserByEmail(ctx, exData.params.Email.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found so password can not be reset")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get user by email")
		exData.err = catalogs.ErrDBResourceRetrieval

		return nil
	}

	if row.Grade == queries.UserGradePENDING || !row.EmailConfirm {
		exData.logger.Warn().Msg("user is not active or has not confirmed email")

		return nil
	}

	exData.user = &row

	return createPasswordResetEntry
}

func createPasswordResetEntry(ctx context.Context, env *PasswordService, exData *newPasswordExchangeData) newPasswordState {
	uuid, err := utils.CreateUUID()
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to generate token")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	if err := env.postgresService.Queries.CreateResetToken(ctx, queries.CreateResetTokenParams{
		ID:    exData.user.ID,
		Token: uuid,
	}); err != nil {
		exData.logger.Error().Err(err).Msg("failed to create reset token")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.uuid = uuid

	return sendPasswordResetMail
}

func sendPasswordResetMail(_ context.Context, env *PasswordService, exData *newPasswordExchangeData) newPasswordState {
	env.stopper.Hold(1)

	//nolint:contextcheck
	go func(env *PasswordService, uuid string, params queries.GetUserByEmailRow) {
		defer env.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		_ = env.smtpService.SendPasswordResetMail(ctx, []smtp.EmailPerson{
			{
				Name:  params.Prenom + " " + params.Nom,
				Email: params.Email,
			},
		}, uuid)
	}(env, exData.uuid, *exData.user)

	return nil
}

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

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) Login(w http.ResponseWriter, r *http.Request) *_http.APIError {
	req, err := _http.DecodeAndValidateJSONBody[*api.PostLogin](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	token, err := processLogin(r.Context(), u, req)
	if err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	return _http.WriteJSONResponse(w, http.StatusOK, token)
}

type loginExchangeData struct {
	logger *zerolog.Logger
	result loginResult
	params *api.PostLogin
	user   queries.GetUserByEmailRow
}

type loginResult struct {
	token *api.UserLogin
	err   error
}

type loginState func(ctx context.Context, s *UserService, data *loginExchangeData) loginState

func processLogin(ctx context.Context, s *UserService, req *api.PostLogin) (*api.UserLogin, error) {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := loginExchangeData{
		logger: logger,
		result: loginResult{},
		params: req,
	}

	respChan := make(chan loginResult, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := loginGetUserByEmail; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.result

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during user login")

		return nil, catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp.token, resp.err
	}
}

func loginGetUserByEmail(ctx context.Context, s *UserService, exData *loginExchangeData) loginState {
	row, err := s.postgresService.Queries.GetUserByEmail(ctx, exData.params.Email.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.result.err = catalogs.ErrUserNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get user by email")
		exData.result.err = catalogs.ErrDBResourceRetrieval

		return nil
	}

	if row.Grade == queries.UserGradePENDING {
		exData.logger.Warn().Msg("user is pending")
		exData.result.err = catalogs.ErrUserNotActive

		return nil
	}

	if !row.EmailConfirm {
		exData.logger.Warn().Msg("user has not confirmed email")
		exData.result.err = catalogs.ErrMailNotConfirmed

		return nil
	}

	exData.user = row

	return loginCheckPassword
}

func loginCheckPassword(_ context.Context, _ *UserService, exData *loginExchangeData) loginState {
	passwordBytes := []byte(*exData.params.Password)
	hashedPasswordBytes := []byte(exData.user.MotDePasse)

	if err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes); err != nil {
		exData.result.err = catalogs.ErrWrongCredentials

		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// other errors should be logged
			exData.logger.Error().Err(err).Msg("unexpected error during user password verification")
		}

		return nil
	}

	return loginGenerateAccessToken
}

func loginGenerateAccessToken(_ context.Context, s *UserService, exData *loginExchangeData) loginState {
	duration := s.config.Session.Duration

	accessToken := jwt.NewToken(s.config.Session.Signer, duration, string(exData.user.Grade), exData.user.ID)

	signed, err := accessToken.Signed()
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to sign access token")

		exData.result.err = catalogs.ErrUnexpectedError
	}

	exData.result.token = &api.UserLogin{
		Token:     &signed,
		TokenType: utils.PtrTo("Bearer"),
		ExpiresIn: utils.PtrTo(utils.RoundFloat64(duration.Hours(), 0)),
		Firstname: &exData.user.Prenom,
	}

	return loginUpdateLastLogin
}

func loginUpdateLastLogin(ctx context.Context, s *UserService, exData *loginExchangeData) loginState {
	err := s.postgresService.Queries.UpdateLastLogin(ctx, exData.user.ID)
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to update last login")
	}

	return nil
}

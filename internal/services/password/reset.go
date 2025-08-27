package password

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"unicode"

	"irj/internal/catalogs"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (p *PasswordService) ResetPassword(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token := httprouter.ParamsFromContext(r.Context()).ByName("token")
	if token == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("token is required")
	}

	req, err := _http.DecodeAndValidateJSONBody[*api.PostNewPwd](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	if err := processResetPassword(r.Context(), p, token, req); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type resetPwdExchangeData struct {
	logger   *zerolog.Logger
	err      error
	token    string
	userID   string
	params   *api.PostNewPwd
	hashedPw string
}

type resetPwdState func(ctx context.Context, env *PasswordService, data *resetPwdExchangeData) resetPwdState

func processResetPassword(ctx context.Context, e *PasswordService, token string, req *api.PostNewPwd) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := resetPwdExchangeData{
		logger: logger,
		token:  token,
		params: req,
	}

	respChan := make(chan error, 1)

	e.stopper.Hold(1)

	go func() {
		defer e.stopper.Release()

		for state := checkResetPwdToken; state != nil; {
			state = state(ctx, e, &exData)
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

func checkResetPwdToken(ctx context.Context, env *PasswordService, exData *resetPwdExchangeData) resetPwdState {
	token, err := env.postgresService.Queries.GetResetPasswordByToken(ctx, exData.token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("token not found")
			exData.err = catalogs.ErrTokenNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get token")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.userID = token.UserID

	return resetPwdCheckPwd
}

//nolint:cyclop
func resetPwdCheckPwd(_ context.Context, _ *PasswordService, exData *resetPwdExchangeData) resetPwdState {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range *exData.params.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasUpper && hasLower && hasNumber && hasSpecial {
		return hashResetPassword
	}

	exData.logger.Warn().Err(catalogs.ErrInvalidPassword).Msg("password does not match rules")
	exData.err = catalogs.ErrInvalidPassword

	return nil
}

func hashResetPassword(_ context.Context, _ *PasswordService, exData *resetPwdExchangeData) resetPwdState {
	password := []byte(*exData.params.Password)

	hashedPwB, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			exData.err = catalogs.ErrPasswordTooLong

			return nil
		}

		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.hashedPw = string(hashedPwB)

	return resetPwd
}

func resetPwd(ctx context.Context, env *PasswordService, exData *resetPwdExchangeData) resetPwdState {
	err := env.postgresService.Queries.UpdatePasswordByID(ctx, queries.UpdatePasswordByIDParams{
		ID:       exData.userID,
		Password: exData.hashedPw,
	})
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to update password")
		exData.err = catalogs.ErrUnexpectedError
	}

	return nil
}

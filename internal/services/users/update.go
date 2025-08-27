package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"unicode"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) UpdateUserInfo(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("id is required")
	}

	req, err := _http.DecodeAndValidateJSONBody[*api.PutUsersBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	if err := processUpdateUser(r.Context(), u, id, &token, req); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type updateUserExchangeData struct {
	logger      *zerolog.Logger
	err         error
	id          string
	params      *api.PutUsersBody
	token       *jwt.SessionInfo
	user        *queries.GetUserByIDRow
	emailChange bool
	hashedPw    string
}

type updateUserState func(ctx context.Context, s *UserService, data *updateUserExchangeData) updateUserState

func processUpdateUser(ctx context.Context, s *UserService, id string, token *jwt.SessionInfo, params *api.PutUsersBody) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := updateUserExchangeData{
		logger: logger,
		id:     id,
		params: params,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := updateUserCheck; state != nil; {
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

func updateUserCheck(ctx context.Context, s *UserService, exData *updateUserExchangeData) updateUserState {
	if exData.id != exData.token.ID && exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot update other users")
		exData.err = catalogs.ErrUserNotAdmin

		return nil
	}

	if exData.params.Password != "" && exData.params.OldPassword == "" {
		exData.logger.Warn().Msg("password is changed but old password is not provided")
		exData.err = catalogs.ErrInvalidParam

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

	return updateUserCheckSensitiveChange
}

func updateUserCheckSensitiveChange(_ context.Context, _ *UserService, exData *updateUserExchangeData) updateUserState {
	if exData.user.Email != exData.params.Mail.String() {
		exData.emailChange = true
	}

	if exData.params.Password != "" {
		return checkUserUpdateOldPassword
	}

	exData.hashedPw = exData.user.MotDePasse

	return updateUser
}

func checkUserUpdateOldPassword(_ context.Context, _ *UserService, exData *updateUserExchangeData) updateUserState {
	passwordBytes := []byte(exData.params.OldPassword)
	hashedPasswordBytes := []byte(exData.user.MotDePasse)

	if err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes); err != nil {
		exData.err = catalogs.ErrInvalidOldPassword

		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// other errors should be logged
			exData.logger.Error().Err(err).Msg("unexpected error during user password verification")
		}

		return nil
	}

	return checkUserUpdatePassword
}

//nolint:cyclop
func checkUserUpdatePassword(_ context.Context, _ *UserService, exData *updateUserExchangeData) updateUserState {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range exData.params.Password {
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
		return hashUserUpdatePassword
	}

	exData.logger.Warn().Err(catalogs.ErrInvalidPassword).Msg("password does not match rules")
	exData.err = catalogs.ErrInvalidPassword

	return nil
}

func hashUserUpdatePassword(_ context.Context, _ *UserService, exData *updateUserExchangeData) updateUserState {
	password := []byte(exData.params.Password)

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

	return updateUser
}

func updateUser(ctx context.Context, s *UserService, exData *updateUserExchangeData) updateUserState {
	if err := s.postgresService.Queries.UpdateUserByID(ctx, queries.UpdateUserByIDParams{
		ID:           exData.id,
		Firstname:    *exData.params.Firstname,
		Name:         *exData.params.Lastname,
		Email:        exData.params.Mail.String(),
		EmailConfirm: !exData.emailChange,
		Password:     exData.hashedPw,
		Phone:        pgtype.Text{String: exData.params.Phone, Valid: true},
		Organization: pgtype.Text{String: exData.params.Organization, Valid: true},
		Domain:       queries.DomaineExpertise(*exData.params.Domain),
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found and therefore can not be updated")
			exData.err = catalogs.ErrUserNotFound

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to update user")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	if exData.emailChange {
		return sendUserUpdateEmailChange
	}

	return nil
}

func sendUserUpdateEmailChange(_ context.Context, s *UserService, exData *updateUserExchangeData) updateUserState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(s *UserService, id string, params api.PutUsersBody) {
		defer s.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		_ = s.smtpService.SendEmailConfirmationMail(ctx, []smtp.EmailPerson{
			{
				Name:  *params.Firstname + " " + *params.Lastname,
				Email: params.Mail.String(),
			},
		}, id)
	}(s, exData.id, *exData.params)

	return nil
}

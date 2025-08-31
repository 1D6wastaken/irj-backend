package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"unicode"

	"irj/internal/catalogs"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/api"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) RegisterUser(w http.ResponseWriter, r *http.Request) *_http.APIError {
	req, err := _http.DecodeAndValidateJSONBody[*api.PostUsersBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	token, err := processRegistration(r.Context(), u, req)
	if err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	return _http.WriteJSONResponse(w, http.StatusCreated, api.UserID{Token: (*strfmt.UUID)(&token)})
}

type registrationExchangeData struct {
	logger   *zerolog.Logger
	result   registrationResult
	params   *api.PostUsersBody
	hashedPw string
}

type registrationResult struct {
	token string
	err   error
}

type registrationState func(ctx context.Context, s *UserService, data *registrationExchangeData) registrationState

func processRegistration(ctx context.Context, s *UserService, req *api.PostUsersBody) (string, error) {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := registrationExchangeData{
		logger: logger,
		result: registrationResult{},
		params: req,
	}

	respChan := make(chan registrationResult, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := checkUserRegistrationEmailExistence; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.result

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during user registration")

		return "", catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp.token, resp.err
	}
}

func checkUserRegistrationEmailExistence(ctx context.Context, s *UserService, exData *registrationExchangeData) registrationState {
	_, err := s.postgresService.Queries.GetUserByEmail(ctx, exData.params.Mail.String())
	if errors.Is(err, sql.ErrNoRows) {
		return checkUserRegistrationPassword
	}

	if err != nil {
		exData.result.err = catalogs.ErrDBResourceRetrieval

		return nil
	}

	exData.result.err = catalogs.ErrUserAlreadyRegistered

	return nil
}

//nolint:cyclop
func checkUserRegistrationPassword(_ context.Context, _ *UserService, exData *registrationExchangeData) registrationState {
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
		return hashUserRegistrationPassword
	}

	exData.logger.Warn().Err(catalogs.ErrInvalidPassword).Msg("password does not match rules")
	exData.result.err = catalogs.ErrInvalidPassword

	return nil
}

func hashUserRegistrationPassword(_ context.Context, _ *UserService, exData *registrationExchangeData) registrationState {
	password := []byte(*exData.params.Password)

	hashedPwB, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			exData.result.err = catalogs.ErrPasswordTooLong

			return nil
		}

		exData.result.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.hashedPw = string(hashedPwB)

	return createUserRegistration
}

func createUserRegistration(ctx context.Context, s *UserService, exData *registrationExchangeData) registrationState {
	uuid, err := utils.CreateUUID()
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to generate token")
		exData.result.err = catalogs.ErrUnexpectedError

		return nil
	}

	err = s.postgresService.Queries.CreateUser(ctx, queries.CreateUserParams{
		ID:        uuid,
		Firstname: *exData.params.Firstname,
		Name:      *exData.params.Lastname,
		Email:     exData.params.Mail.String(),
		Password:  exData.hashedPw,
		Phone: pgtype.Text{
			String: exData.params.Phone,
			Valid:  true,
		},
		Organization: pgtype.Text{
			String: exData.params.Organization,
			Valid:  true,
		},
		Domain: queries.DomaineExpertise(*exData.params.Domain),
		Motivation: pgtype.Text{
			String: *exData.params.Motivation,
			Valid:  true,
		},
	})
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to create user")
		exData.result.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.result.token = uuid

	return storeUserRegistrationEvent
}

func storeUserRegistrationEvent(_ context.Context, s *UserService, exData *registrationExchangeData) registrationState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, id string) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.ContributorRegistrationEvent(context.Background(), id)
		if err != nil {
			logger.Error().Err(err).Msg("failed to store registration event")
		}
	}(exData.logger, exData.result.token)

	return sendUserRegistrationEmail
}

func sendUserRegistrationEmail(_ context.Context, s *UserService, exData *registrationExchangeData) registrationState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, s *UserService) {
		defer s.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		row, err := s.postgresService.Queries.GetUsersByGrade(ctx, queries.UserGradeADMIN)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get admin user")

			return
		}

		to := make([]smtp.EmailPerson, 0, len(row))

		//nolint:gocritic
		for _, user := range row {
			to = append(to, smtp.EmailPerson{
				Name:  user.Prenom + " " + user.Nom,
				Email: user.Email,
			})
		}

		_ = s.smtpService.SendNewUserEmail(ctx, to)
	}(exData.logger, s)

	return nil
}

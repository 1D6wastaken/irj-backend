package business

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

func (b *BusinessService) CreatePersonnePhysique(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	req, err := _http.DecodeAndValidateJSONBody[*api.PersonnePhysiqueCreationBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	if err := processCreatePersonnePhysique(r.Context(), b, &token, req); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

type createPersonnePhysiqueExchangeData struct {
	logger *zerolog.Logger
	err    error
	params *api.PersonnePhysiqueCreationBody
	token  *jwt.SessionInfo
	id     int32
}

//nolint:lll
type createPersonnePhysiqueState func(ctx context.Context, s *BusinessService, data *createPersonnePhysiqueExchangeData) createPersonnePhysiqueState

//nolint:lll
func processCreatePersonnePhysique(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.PersonnePhysiqueCreationBody) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := createPersonnePhysiqueExchangeData{
		logger: logger,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := createPersonnePhysique; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during personne physique creation")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

//nolint:lll
func createPersonnePhysique(ctx context.Context, s *BusinessService, exData *createPersonnePhysiqueExchangeData) createPersonnePhysiqueState {
	id, err := s.postgresService.Queries.CreatePersPhysique(ctx, queries.CreatePersPhysiqueParams{
		PrenomNomPersPhy:      pgtype.Text{String: *exData.params.Title, Valid: true},
		Commentaires:          pgtype.Text{String: exData.params.Comment, Valid: exData.params.Comment != ""},
		DateNaissance:         pgtype.Text{String: exData.params.Birthday, Valid: exData.params.Birthday != ""},
		DateDeces:             pgtype.Text{String: exData.params.Death, Valid: exData.params.Death != ""},
		Attestation:           pgtype.Text{String: exData.params.Attestation, Valid: exData.params.Attestation != ""},
		ElementsBiographiques: pgtype.Text{String: exData.params.BiographicalElements, Valid: exData.params.BiographicalElements != ""},
		ElementsPelerinage:    pgtype.Text{String: exData.params.PilgrimageElements, Valid: exData.params.PilgrimageElements != ""},
		NatureEvenement:       pgtype.Text{String: exData.params.Nature, Valid: exData.params.Nature != ""},
		CommutationVoeu:       pgtype.Text{String: exData.params.Commutation, Valid: exData.params.Commutation != ""},
		Bibliographie:         pgtype.Text{String: exData.params.Bibliography, Valid: exData.params.Bibliography != ""},
		Sources:               pgtype.Text{String: exData.params.Source, Valid: exData.params.Source != ""},
		Contributeurs:         pgtype.Text{String: strings.Join(exData.params.Contributors, ","), Valid: len(exData.params.Contributors) > 0},
		IDCommune: pgtype.Int4{
			Int32: exData.params.City,
			Valid: exData.params.City != 0,
		},
		IDPays: pgtype.Int4{
			Int32: exData.params.Country,
			Valid: exData.params.Country != 0,
		},
	})
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to insert personne physique")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.logger.Info().Int32("id", id).Msg("personne physique created")

	exData.id = id

	return storePersPhyDocumentSubmissionEvent
}

//nolint:lll
func storePersPhyDocumentSubmissionEvent(_ context.Context, s *BusinessService, exData *createPersonnePhysiqueExchangeData) createPersonnePhysiqueState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, userID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentSubmissionEvent(context.Background(), queries.DocumentSubmissionEventParams{
			UserID: userID,
			DocumentID: pgtype.Int4{
				Int32: documentID,
				Valid: true,
			},
			Comment: pgtype.Text{
				String: "personnes_physiques",
				Valid:  true,
			},
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to store document submission event")
		}
	}(exData.logger, exData.token.ID, exData.id)

	return linkPersonnePhysique
}

func linkPersonnePhysique(ctx context.Context, s *BusinessService, exData *createPersonnePhysiqueExchangeData) createPersonnePhysiqueState {
	err := s.postgresService.Queries.AttachSieclesToPersPhy(ctx, queries.AttachSieclesToPersPhyParams{
		SiecleID: exData.params.Centuries,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach centuries to personne physique")
	}

	err = s.postgresService.Queries.AttachMediasToPersPhy(ctx, queries.AttachMediasToPersPhyParams{
		MediaIds: exData.params.Medias,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach medias to personne physique")
	}

	err = s.postgresService.Queries.AttachThemesToPersPhy(ctx, queries.AttachThemesToPersPhyParams{
		ThemeIds: exData.params.Themes,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach themes to personne physique")
	}

	err = s.postgresService.Queries.AttachHistoricalPeriodsToPersPhy(ctx, queries.AttachHistoricalPeriodsToPersPhyParams{
		PeriodeIds: exData.params.HistoricalPeriods,
		ID:         exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach historical periods to personne physique")
	}

	err = s.postgresService.Queries.AttachProfessionsToPersPhy(ctx, queries.AttachProfessionsToPersPhyParams{
		ProfessionIds: exData.params.Professions,
		ID:            exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach professions to personne physique")
	}

	err = s.postgresService.Queries.AttachModeDeTransportsToPersPhy(ctx, queries.AttachModeDeTransportsToPersPhyParams{
		TravelIds: exData.params.Travels,
		ID:        exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach travel modes to personne physique")
	}

	err = s.postgresService.Queries.LinkPersPhyToMonuLieu(ctx, queries.LinkPersPhyToMonuLieuParams{
		MonuLieuIds: exData.params.LinkedMonumentsLieux,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link monuments lieux document to personne physique")
	}

	err = s.postgresService.Queries.LinkPersPhyToMobImg(ctx, queries.LinkPersPhyToMobImgParams{
		MobImgIds: exData.params.LinkedMobiliersImages,
		ID:        exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link mobiliers images document to personne physique")
	}

	err = s.postgresService.Queries.LinkPersPhyToPersMo(ctx, queries.LinkPersPhyToPersMoParams{
		PersoMoIds: exData.params.LinkedPersMorales,
		ID:         exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link personnes morales document to personne physique")
	}

	return addAuteurToPersonnePhysique
}

//nolint:lll
func addAuteurToPersonnePhysique(ctx context.Context, s *BusinessService, exData *createPersonnePhysiqueExchangeData) createPersonnePhysiqueState {
	user, err := s.postgresService.Queries.GetUserByID(ctx, exData.token.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("user not found and therefore can not be added as an author")
		}

		exData.logger.Error().Err(err).Msg("failed to get user")
	}

	var id int32

	auteurs, err := s.postgresService.Queries.GetAuteurByName(ctx, pgtype.Text{
		String: user.Prenom + " " + user.Nom,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			id, err = s.postgresService.Queries.CreateAuteur(ctx, pgtype.Text{
				String: user.Prenom + " " + user.Nom,
				Valid:  true,
			})
			if err != nil {
				exData.logger.Error().Err(err).Msg("failed to create author")
			}
		} else {
			exData.logger.Error().Err(err).Msg("failed to get auteur")
		}
	} else {
		id = auteurs.IDAuteurFiche
	}

	err = s.postgresService.Queries.AttachAuthorToPersPhy(ctx, queries.AttachAuthorToPersPhyParams{
		AuteurID: id,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach author to personne physique")
	}

	return sendNewPersonnePhysiqueDocumentMail
}

//nolint:lll
func sendNewPersonnePhysiqueDocumentMail(_ context.Context, s *BusinessService, exData *createPersonnePhysiqueExchangeData) createPersonnePhysiqueState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(s *BusinessService, logger *zerolog.Logger) {
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

		_ = s.smtpService.SendNewDocumentMail(ctx, to)
	}(s, exData.logger)

	return nil
}

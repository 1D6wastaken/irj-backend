package business

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	"irj/internal/postgres"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (b *BusinessService) UpdatePersonnePhysique(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	req, err := _http.DecodeJSONBody[*api.PersonnePhysiqueCreationBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	if !req.Draft {
		if err := req.Validate(nil); err != nil {
			return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
		}
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("id is required")
	}

	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid").Err(err)
	}

	if err := processUpdatePersonnePhysique(r.Context(), b, &token, req, int32(intID)); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type updatePersonnePhysiqueExchangeData struct {
	logger               *zerolog.Logger
	err                  error
	parentID             int32
	id                   int32
	params               *api.PersonnePhysiqueCreationBody
	token                *jwt.SessionInfo
	draftToDelete        int32
	originalDateCreation pgtype.Date
}

//nolint:lll
type updatePersonnePhysiqueState func(ctx context.Context, s *BusinessService, data *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState

//nolint:lll
func processUpdatePersonnePhysique(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.PersonnePhysiqueCreationBody, id int32) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := updatePersonnePhysiqueExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := getPersonnePhysiqueToUpdate; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during personne physique update")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

func getPersonnePhysiqueToUpdate(ctx context.Context, s *BusinessService,
	exData *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState {
	p, err := s.postgresService.Queries.GetPersonnePhysiqueByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Int32("id", exData.id).Msg("personne physique not found")

			exData.err = catalogs.ErrDBResourceNotFound

			return nil
		}

		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to get personne physique by id")
		exData.err = catalogs.ErrDBResourceRetrieval

		return nil
	}

	exData.originalDateCreation = p.DateCreation

	if p.PublicationStatus == postgres.DraftPublicationStatus {
		exData.logger.Info().Int32("id", exData.id).Msg("updating draft personne physique")

		exData.draftToDelete = exData.id

		if p.ParentID.Valid {
			exData.id = p.ParentID.Int32
		} else {
			exData.id = 0
		}
	}

	return updatePersonnePhysique
}

//nolint:lll
func updatePersonnePhysique(ctx context.Context, s *BusinessService, exData *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState {
	publicationStatus := postgres.PendingPublicationStatus
	if exData.params.Draft {
		publicationStatus = postgres.DraftPublicationStatus
	}

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
		IDDepartement: pgtype.Int4{
			Int32: exData.params.Department,
			Valid: exData.params.Department != 0,
		},
		IDRegion: pgtype.Int4{
			Int32: exData.params.Region,
			Valid: exData.params.Region != 0,
		},
		IDPays: pgtype.Int4{
			Int32: exData.params.Country,
			Valid: exData.params.Country != 0,
		},
		Historiographie:   pgtype.Text{String: exData.params.Historiography, Valid: exData.params.Historiography != ""},
		Evenements:        pgtype.Text{String: exData.params.Evenements, Valid: exData.params.Evenements != ""},
		Preparatifs:       pgtype.Text{String: exData.params.Preparatifs, Valid: exData.params.Preparatifs != ""},
		CheminSuivi:       pgtype.Text{String: exData.params.CheminSuivi, Valid: exData.params.CheminSuivi != ""},
		Arrivee:           pgtype.Text{String: exData.params.Arrivee, Valid: exData.params.Arrivee != ""},
		Retour:            pgtype.Text{String: exData.params.Retour, Valid: exData.params.Retour != ""},
		NonExecution:      pgtype.Text{String: exData.params.NonExecution, Valid: exData.params.NonExecution != ""},
		Age:               pgtype.Text{String: exData.params.Age, Valid: exData.params.Age != ""},
		CompositionGroupe: pgtype.Text{String: exData.params.CompositionGroupe, Valid: exData.params.CompositionGroupe != ""},
		DateCreation:      exData.originalDateCreation,
		PublicationStatus: publicationStatus,
		ParentID:          pgtype.Int4{Int32: exData.id, Valid: exData.id != 0},
		UserID: pgtype.Text{
			String: exData.token.ID,
			Valid:  true,
		},
	})
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to insert personne physique")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.logger.Info().Int32("id", id).Msg("personne physique created")

	exData.parentID = exData.id

	exData.id = id

	return linkUpdatedPersonnePhysique
}

//nolint:lll,cyclop
func linkUpdatedPersonnePhysique(ctx context.Context, s *BusinessService, exData *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState {
	err := s.postgresService.Queries.AttachSieclesToPersPhy(ctx, queries.AttachSieclesToPersPhyParams{
		SiecleID: exData.params.Centuries,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach centuries to personne physique")
	}

	origMediaIds := exData.params.Medias
	mediaIds := origMediaIds

	if exData.parentID != 0 {
		parentMediaIds, _ := s.postgresService.Queries.GetMediaIdsByPersPhy(ctx, exData.parentID)
		mediaIds = s.duplicateMediasIfShared(ctx, exData.logger, origMediaIds, parentMediaIds)
	}

	s.applyMediaTitleChanges(ctx, exData.logger, origMediaIds, mediaIds, exData.params.MediaTitles)

	err = s.postgresService.Queries.AttachMediasToPersPhy(ctx, queries.AttachMediasToPersPhyParams{
		MediaIds: mediaIds,
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

	err = s.postgresService.Queries.LinkPersPhyToPersPhy(ctx, queries.LinkPersPhyToPersPhyParams{
		PersPhyIds: exData.params.LinkedPersPhysiques,
		ID:         exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link personnes physiques document to updated personne physique")
	}

	return addAuteurToUpdatedPersonnePhysique
}

//nolint:lll
func addAuteurToUpdatedPersonnePhysique(ctx context.Context, s *BusinessService, exData *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState {
	auteurID, err := s.postgresService.Queries.GetAuteurIDByPersPhyID(ctx, exData.parentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("no author found for the original personne physique, adding a new one")

			return linkNewAuteurToPersPhy
		}

		exData.logger.Error().Err(err).Int32("id", exData.parentID).Msg("failed to get author by personne physique id")

		if exData.params.Draft {
			if exData.draftToDelete != 0 {
				return deleteOldDraftPersonnePhysiqueParent
			}

			return nil
		}

		return storePersPhyDocumentUpdateEvent
	}

	err = s.postgresService.Queries.AttachAuthorToPersPhy(ctx, queries.AttachAuthorToPersPhyParams{
		AuteurID: auteurID,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach author to personne physique")
	}

	if exData.params.Draft {
		if exData.draftToDelete != 0 {
			return deleteOldDraftPersonnePhysiqueParent
		}

		return nil
	}

	return storePersPhyDocumentUpdateEvent
}

//nolint:lll
func linkNewAuteurToPersPhy(ctx context.Context, s *BusinessService, exData *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState {
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
			id, err = s.postgresService.Queries.CreateAuteur(ctx, queries.CreateAuteurParams{
				Name: pgtype.Text{
					String: user.Prenom + " " + user.Nom,
					Valid:  true,
				},
				UserID: pgtype.Text{
					String: exData.token.ID,
					Valid:  true,
				},
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

	if exData.params.Draft {
		if exData.draftToDelete != 0 {
			return deleteOldDraftPersonnePhysiqueParent
		}

		return nil
	}

	return storePersPhyDocumentUpdateEvent
}

//nolint:lll
func storePersPhyDocumentUpdateEvent(_ context.Context, s *BusinessService, exData *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, userID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentUpdateEvent(context.Background(), queries.DocumentUpdateEventParams{
			UserID: pgtype.Text{
				String: userID,
				Valid:  true,
			},
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
			logger.Error().Err(err).Msg("failed to store document update event")
		}
	}(exData.logger, exData.token.ID, exData.id)

	if exData.draftToDelete != 0 {
		return deleteOldDraftPersonnePhysiqueParent
	}

	return nil
}

//nolint:lll
func deleteOldDraftPersonnePhysiqueParent(_ context.Context, s *BusinessService, exData *updatePersonnePhysiqueExchangeData) updatePersonnePhysiqueState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, id int32) {
		defer s.stopper.Release()

		if err := deletePersonnePhysique(context.Background(), logger, s, id); err != nil {
			logger.Error().Err(err).Int32("id", id).Msg("failed to delete old draft personne physique parent")
		}
	}(exData.logger, exData.draftToDelete)

	return nil
}

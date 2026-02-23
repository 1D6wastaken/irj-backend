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

func (b *BusinessService) UpdateMobilierImage(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	req, err := _http.DecodeJSONBody[*api.MobilierImageCreationBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	if !req.Draft {
		err = req.Validate(nil)
		if err != nil {
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

	if err := processUpdateMobilierImage(r.Context(), b, &token, req, int32(intID)); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type updateMobilierImageExchangeData struct {
	logger        *zerolog.Logger
	err           error
	parentID      int32
	id            int32
	params        *api.MobilierImageCreationBody
	token         *jwt.SessionInfo
	draftToDelete int32
}

type updateMobilierImageState func(ctx context.Context, s *BusinessService, data *updateMobilierImageExchangeData) updateMobilierImageState

//nolint:lll
func processUpdateMobilierImage(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.MobilierImageCreationBody, id int32) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := updateMobilierImageExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := getMobilierImageToUpdate; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during mobilier image update")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

func getMobilierImageToUpdate(ctx context.Context, s *BusinessService, exData *updateMobilierImageExchangeData) updateMobilierImageState {
	m, err := s.postgresService.Queries.GetMobilierImageByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Int32("id", exData.id).Msg("mobilier image not found")

			exData.err = catalogs.ErrDBResourceNotFound

			return nil
		}

		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to get mobilier image by id")
		exData.err = catalogs.ErrDBResourceRetrieval

		return nil
	}

	if m.PublicationStatus == postgres.DraftPublicationStatus {
		exData.logger.Info().Int32("id", exData.id).Msg("updating draft mobilier image")

		exData.draftToDelete = exData.id

		if m.ParentID.Valid {
			exData.id = m.ParentID.Int32
		} else {
			exData.id = 0
		}
	}

	return updateMobilierImage
}

func updateMobilierImage(ctx context.Context, s *BusinessService, exData *updateMobilierImageExchangeData) updateMobilierImageState {
	publicationStatus := postgres.PendingPublicationStatus
	if exData.params.Draft {
		publicationStatus = postgres.DraftPublicationStatus
	}

	id, err := s.postgresService.Queries.CreateMobilierImage(ctx, queries.CreateMobilierImageParams{
		TitreMobImg:   *exData.params.Title,
		Description:   pgtype.Text{String: exData.params.Description, Valid: true},
		Historique:    pgtype.Text{String: exData.params.History, Valid: exData.params.History != ""},
		Inscriptions:  pgtype.Text{String: exData.params.Inscription, Valid: exData.params.Inscription != ""},
		Origin:        pgtype.Text{String: exData.params.OriginPlace, Valid: exData.params.OriginPlace != ""},
		Place:         pgtype.Text{String: exData.params.PresentPlace, Valid: exData.params.PresentPlace != ""},
		Bibliographie: pgtype.Text{String: exData.params.Bibliography, Valid: exData.params.Bibliography != ""},
		Protection: pgtype.Bool{
			Bool:  exData.params.IsProtected,
			Valid: true,
		},
		ProtectionComment: pgtype.Text{String: exData.params.ProtectionComment, Valid: exData.params.ProtectionComment != ""},
		Source:            pgtype.Text{String: exData.params.Source, Valid: exData.params.Source != ""},
		Contributors:      pgtype.Text{String: strings.Join(exData.params.Contributors, ","), Valid: len(exData.params.Contributors) > 0},
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
		PublicationStatus: publicationStatus,
		ParentID:          pgtype.Int4{Int32: exData.id, Valid: exData.id != 0},
		UserID: pgtype.Text{
			String: exData.token.ID,
			Valid:  exData.params.Draft,
		},
		TemoinCommentaires: pgtype.Text{
			String: exData.params.TemoinComment,
			Valid:  exData.params.TemoinComment != "",
		},
	})
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to insert mobilier image")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.logger.Info().Int32("id", id).Msg("mobilier image created")

	exData.parentID = exData.id

	exData.id = id

	return linkUpdatedMobilierImage
}

//nolint:cyclop
func linkUpdatedMobilierImage(ctx context.Context, s *BusinessService, exData *updateMobilierImageExchangeData) updateMobilierImageState {
	err := s.postgresService.Queries.AttachSieclesToMobImg(ctx, queries.AttachSieclesToMobImgParams{
		SiecleID: exData.params.Centuries,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach centuries to mobilier image")
	}

	err = s.postgresService.Queries.AttachMediasToMobImg(ctx, queries.AttachMediasToMobImgParams{
		MediaIds: exData.params.Medias,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach medias to mobilier image")
	}

	err = s.postgresService.Queries.AttachThemesToMobImg(ctx, queries.AttachThemesToMobImgParams{
		ThemeIds: exData.params.Themes,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach themes to mobilier image")
	}

	err = s.postgresService.Queries.AttachNaturesToMobImg(ctx, queries.AttachNaturesToMobImgParams{
		NatureIds: exData.params.Natures,
		ID:        exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach natures to mobilier image")
	}

	err = s.postgresService.Queries.AttachEtatsToMobImg(ctx, queries.AttachEtatsToMobImgParams{
		EtatIds: exData.params.ConservationStates,
		ID:      exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach conservation states to mobilier image")
	}

	err = s.postgresService.Queries.AttachMateriauxToMobImg(ctx, queries.AttachMateriauxToMobImgParams{
		MateriauIds: exData.params.Materials,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach materials to mobilier image")
	}

	err = s.postgresService.Queries.AttachTechniquesToMobImg(ctx, queries.AttachTechniquesToMobImgParams{
		TechniquesIds: exData.params.Techniques,
		ID:            exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach techniques to mobilier image")
	}

	err = s.postgresService.Queries.LinkMobImgToMonuLieu(ctx, queries.LinkMobImgToMonuLieuParams{
		MonuLieuIds: exData.params.LinkedMonumentsLieux,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link monuments lieux document to mobilier image")
	}

	err = s.postgresService.Queries.LinkMobImgToPersMo(ctx, queries.LinkMobImgToPersMoParams{
		PersoMoIds: exData.params.LinkedPersMorales,
		ID:         exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link personnes morales document to mobilier image")
	}

	err = s.postgresService.Queries.LinkMobImgToPersPhy(ctx, queries.LinkMobImgToPersPhyParams{
		PersoPhyIds: exData.params.LinkedPersPhysiques,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link personnes physiques document to mobilier image")
	}

	return addAuteursToUpdatedMobilierImage
}

//nolint:lll
func addAuteursToUpdatedMobilierImage(ctx context.Context, s *BusinessService, exData *updateMobilierImageExchangeData) updateMobilierImageState {
	auteurID, err := s.postgresService.Queries.GetAuteurIDByMobImg(ctx, exData.parentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("no auteur id found for the original mobilier image, adding a new one")

			return linkNewAuteurToMobilierImage
		}

		exData.logger.Error().Err(err).Int32("id", exData.parentID).Msg("failed to get author by mobilier image")

		if exData.params.Draft {
			if exData.draftToDelete != 0 {
				return deleteOldDraftMobilierImageParent
			}

			return nil
		}

		return storeMobImgDocumentUpdateEvent
	}

	err = s.postgresService.Queries.AttachAuthorToMobImg(ctx, queries.AttachAuthorToMobImgParams{
		AuteurID: auteurID,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach author to updated mobilier image")
	}

	if exData.params.Draft {
		if exData.draftToDelete != 0 {
			return deleteOldDraftMobilierImageParent
		}

		return nil
	}

	return storeMobImgDocumentUpdateEvent
}

//nolint:lll
func linkNewAuteurToMobilierImage(ctx context.Context, s *BusinessService, exData *updateMobilierImageExchangeData) updateMobilierImageState {
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

	err = s.postgresService.Queries.AttachAuthorToMobImg(ctx, queries.AttachAuthorToMobImgParams{
		AuteurID: id,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach author to updated mobilier image")
	}

	err = s.postgresService.Queries.AttachAuthorToMobImg(ctx, queries.AttachAuthorToMobImgParams{
		AuteurID: id,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach author to updated mobilier image")
	}

	if exData.params.Draft {
		if exData.draftToDelete != 0 {
			return deleteOldDraftMobilierImageParent
		}

		return nil
	}

	return storeMobImgDocumentUpdateEvent
}

//nolint:lll
func storeMobImgDocumentUpdateEvent(_ context.Context, s *BusinessService, exData *updateMobilierImageExchangeData) updateMobilierImageState {
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
				String: "mobiliers_images",
				Valid:  true,
			},
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to store document update event")
		}
	}(exData.logger, exData.token.ID, exData.id)

	if exData.draftToDelete != 0 {
		return deleteOldDraftMobilierImageParent
	}

	return nil
}

//nolint:lll
func deleteOldDraftMobilierImageParent(_ context.Context, s *BusinessService, exData *updateMobilierImageExchangeData) updateMobilierImageState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, id int32) {
		defer s.stopper.Release()

		if err := deleteMobilierImage(context.Background(), logger, s, id); err != nil {
			logger.Error().Err(err).Int32("id", id).Msg("failed to delete old draft mobilier image parent")
		}
	}(exData.logger, exData.draftToDelete)

	return nil
}

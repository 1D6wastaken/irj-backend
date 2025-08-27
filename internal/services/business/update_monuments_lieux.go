package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (b *BusinessService) UpdateMonumentLieu(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	req, err := _http.DecodeAndValidateJSONBody[*api.MonumentsLieuxCreationBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("id is required")
	}

	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid").Err(err)
	}

	if err := processUpdateMonumentLieu(r.Context(), b, &token, req, int32(intID)); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type updateMonumentLieuExchangeData struct {
	logger *zerolog.Logger
	err    error
	id     int32
	params *api.MonumentsLieuxCreationBody
	token  *jwt.SessionInfo
}

type updateMonumentLieuState func(ctx context.Context, s *BusinessService, data *updateMonumentLieuExchangeData) updateMonumentLieuState

//nolint:lll
func processUpdateMonumentLieu(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.MonumentsLieuxCreationBody, id int32) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := updateMonumentLieuExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := updateMonumentLieu; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during monument lieu update")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

func updateMonumentLieu(ctx context.Context, s *BusinessService, exData *updateMonumentLieuExchangeData) updateMonumentLieuState {
	id, err := s.postgresService.Queries.CreateMonumentLieu(ctx, queries.CreateMonumentLieuParams{
		TitreMonuLieu: *exData.params.Title,
		Description:   pgtype.Text{String: *exData.params.Description, Valid: true},
		Histoire:      pgtype.Text{String: exData.params.History, Valid: exData.params.History != ""},
		Geolocalisation: pgtype.Text{
			String: fmt.Sprintf("%s, %s", exData.params.Latitude, exData.params.Longitude),
			Valid:  exData.params.Latitude != "" && exData.params.Longitude != "",
		},
		Bibliographie: pgtype.Text{String: exData.params.Bibliography, Valid: exData.params.Bibliography != ""},
		Protection: pgtype.Bool{
			Bool:  exData.params.IsProtected,
			Valid: true,
		},
		ProtectionCommentaires: pgtype.Text{String: exData.params.ProtectionComment, Valid: exData.params.ProtectionComment != ""},
		Source:                 pgtype.Text{String: exData.params.Source, Valid: exData.params.Source != ""},
		Contributeurs:          pgtype.Text{String: strings.Join(exData.params.Contributors, ","), Valid: len(exData.params.Contributors) > 0},
		IDCommune: pgtype.Int4{
			Int32: exData.params.City,
			Valid: exData.params.City != 0,
		},
		IDPays: pgtype.Int4{
			Int32: exData.params.Country,
			Valid: exData.params.Country != 0,
		},
		ParentID: pgtype.Int4{Int32: exData.id, Valid: true},
	})
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to insert monument lieu")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.logger.Info().Int32("id", id).Msg("updated monument lieu created")

	exData.id = id

	return linkUpdatedMonumentLieu
}

func linkUpdatedMonumentLieu(ctx context.Context, s *BusinessService, exData *updateMonumentLieuExchangeData) updateMonumentLieuState {
	err := s.postgresService.Queries.AttachSieclesToMonuLieu(ctx, queries.AttachSieclesToMonuLieuParams{
		SiecleID: exData.params.Centuries,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach centuries to updated monument lieu")
	}

	err = s.postgresService.Queries.AttachMediasToMonuLieu(ctx, queries.AttachMediasToMonuLieuParams{
		MediaIds: exData.params.Medias,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach medias to updated monument lieu")
	}

	err = s.postgresService.Queries.AttachThemesToMonuLieu(ctx, queries.AttachThemesToMonuLieuParams{
		ThemeIds: exData.params.Themes,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach themes to updated monument lieu")
	}

	err = s.postgresService.Queries.AttachNaturesToMonuLieu(ctx, queries.AttachNaturesToMonuLieuParams{
		NatureIds: exData.params.Natures,
		ID:        exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach natures to updated monument lieu")
	}

	err = s.postgresService.Queries.AttachEtatsToMonuLieu(ctx, queries.AttachEtatsToMonuLieuParams{
		EtatIds: exData.params.ConservationStates,
		ID:      exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach conservation states to updated monument lieu")
	}

	err = s.postgresService.Queries.AttachMateriauxToMonuLieu(ctx, queries.AttachMateriauxToMonuLieuParams{
		MateriauIds: exData.params.Materials,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach materials to updated monument lieu")
	}

	err = s.postgresService.Queries.LinkMonuLieuToMobImg(ctx, queries.LinkMonuLieuToMobImgParams{
		MobImgIds: exData.params.LinkedMobiliersImages,
		ID:        exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link mobiliers images document to updated monument lieu")
	}

	err = s.postgresService.Queries.LinkMonuLieuToPersMo(ctx, queries.LinkMonuLieuToPersMoParams{
		PersoMoIds: exData.params.LinkedPersMorales,
		ID:         exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link personnes morales document to updated monument lieu")
	}

	err = s.postgresService.Queries.LinkMonuLieuToPersPhy(ctx, queries.LinkMonuLieuToPersPhyParams{
		PersoPhyIds: exData.params.LinkedPersPhysiques,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link personnes physiques document to updated monument lieu")
	}

	return addAuteurToUpdatedMonumentLieu
}

//nolint:lll
func addAuteurToUpdatedMonumentLieu(ctx context.Context, s *BusinessService, exData *updateMonumentLieuExchangeData) updateMonumentLieuState {
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

	err = s.postgresService.Queries.AttachAuthorToMonuLieu(ctx, queries.AttachAuthorToMonuLieuParams{
		AuteurID: id,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach author to updated monument lieu")
	}

	return nil
}

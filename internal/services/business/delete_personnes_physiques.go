package business

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (b *BusinessService) DeleteDraftPersonnePhysique(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	_, ok := subCtx.Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	id := httprouter.ParamsFromContext(subCtx).ByName("id")
	if id == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("id is required")
	}

	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid").Err(err)
	}

	m, err := b.postgresService.Queries.GetPersonnePhysiqueByID(subCtx, int32(intID))
	if err != nil {
		return _http.ErrNotFound.Msg("unable to get personne physique").Err(err)
	}

	if m.PublicationStatus != queries.PublicationStatusDRAFT {
		return _http.ErrNotFound.Msg("resource not found").Err(err)
	}

	if err = deletePersonnePhysique(subCtx, zerolog.Ctx(subCtx), b, int32(intID)); err != nil {
		return _http.ErrInternalError.Msg("unable to delete personne physique").Err(err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

//nolint:cyclop
func deletePersonnePhysique(ctx context.Context, logger *zerolog.Logger, s *BusinessService, id int32) error {
	err := s.postgresService.Queries.DetachSieclesFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach siecles from personne physique")
	}

	err = s.postgresService.Queries.DetachMediasFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach medias from personne physique")
	}

	err = s.postgresService.Queries.DetachThemesFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach themes from personne physique")
	}

	err = s.postgresService.Queries.DetachHistoricalPeriodsFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach historical periods from personne physique")
	}

	err = s.postgresService.Queries.DetachProfessionsFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach professions from personne physique")
	}

	err = s.postgresService.Queries.DetachModeDeTransportsFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach travel modes from personne physique")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink monuments lieux document from personne physique")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink mobiliers images document from personne physique")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromPersMo(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink personnes morales document from personne physique")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink personnes physiques document from personne physique")
	}

	err = s.postgresService.Queries.DetachAuthorFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach author from personne physique")
	}

	if err := s.postgresService.Queries.DeletePersonnePhysique(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("document not found and therefore can not be rejected")

			return nil
		}

		logger.Error().Err(err).Msg("failed to reject document")

		return catalogs.ErrUnexpectedError
	}

	return nil
}

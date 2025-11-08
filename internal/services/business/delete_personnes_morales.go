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

func (b *BusinessService) DeleteDraftPersonneMorale(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	m, err := b.postgresService.Queries.GetPersonneMoraleByID(subCtx, int32(intID))
	if err != nil {
		return _http.ErrNotFound.Msg("unable to get personne morale").Err(err)
	}

	if m.PublicationStatus != queries.PublicationStatusDRAFT {
		return _http.ErrNotFound.Msg("resource not found").Err(err)
	}

	if err = deletePersonneMorale(subCtx, zerolog.Ctx(subCtx), b, int32(intID)); err != nil {
		return _http.ErrInternalError.Msg("unable to delete personne morale").Err(err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func deletePersonneMorale(ctx context.Context, logger *zerolog.Logger, s *BusinessService, id int32) error {
	err := s.postgresService.Queries.DetachSieclesFromPersMo(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach siecles from personne morale")
	}

	err = s.postgresService.Queries.DetachMediasFromPersMo(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach medias from personne morale")
	}

	err = s.postgresService.Queries.DetachThemesFromPersMo(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach themes from personne morale")
	}

	err = s.postgresService.Queries.DetachNaturesFromPersMo(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach natures from personne morale")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink monuments lieux document from personne morale")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink mobiliers images document from personne morale")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink personnes physiques document from personne morale")
	}

	err = s.postgresService.Queries.DetachAuthorFromPersMo(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach author from personne morale")
	}

	if err := s.postgresService.Queries.DeletePersonneMorale(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		logger.Error().Err(err).Msg("failed to reject document")

		return catalogs.ErrUnexpectedError
	}

	return nil
}

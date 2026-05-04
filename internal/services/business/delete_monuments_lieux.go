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

func (b *BusinessService) DeleteDraftMonumentLieu(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	m, err := b.postgresService.Queries.GetMonumentLieuByID(subCtx, int32(intID))
	if err != nil {
		return _http.ErrNotFound.Msg("unable to get monument lieu").Err(err)
	}

	if m.PublicationStatus != queries.PublicationStatusDRAFT {
		return _http.ErrNotFound.Msg("resource not found").Err(err)
	}

	if err = deleteMonumentLieu(subCtx, zerolog.Ctx(subCtx), b, int32(intID)); err != nil {
		return _http.ErrInternalError.Msg("unable to delete monument lieu").Err(err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

//nolint:cyclop
func deleteMonumentLieu(ctx context.Context, logger *zerolog.Logger, s *BusinessService, id int32) error {
	err := s.postgresService.Queries.DetachSieclesFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach siecles from monument lieu")
	}

	err = s.postgresService.Queries.DetachMediasFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach medias from monument lieu")
	}

	err = s.postgresService.Queries.DetachThemesFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach themes from monument lieu")
	}

	err = s.postgresService.Queries.DetachNaturesFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach natures from monument lieu")
	}

	err = s.postgresService.Queries.DetachEtatsFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach conservation states from monument lieu")
	}

	err = s.postgresService.Queries.DetachMateriauxFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach materials from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink mobiliers images document from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromPersMo(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink personnes morales document from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromPersPhy(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink personnes physiques document from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink monuments lieux document from monument lieu")
	}

	err = s.postgresService.Queries.DetachAuthorFromMonuLieu(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach author from monument lieu")
	}

	if err := s.postgresService.Queries.DeleteMonumentLieu(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		logger.Error().Err(err).Msg("failed to approve document")

		return catalogs.ErrUnexpectedError
	}

	return nil
}

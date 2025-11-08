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

func (b *BusinessService) DeleteDraftMobilierImage(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	m, err := b.postgresService.Queries.GetMobilierImageByID(subCtx, int32(intID))
	if err != nil {
		return _http.ErrNotFound.Msg("unable to get mobilier image").Err(err)
	}

	if m.PublicationStatus != queries.PublicationStatusDRAFT {
		return _http.ErrNotFound.Msg("resource not found").Err(err)
	}

	if err = deleteMobilierImage(subCtx, zerolog.Ctx(subCtx), b, int32(intID)); err != nil {
		return _http.ErrInternalError.Msg("unable to delete mobilier image").Err(err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func deleteMobilierImage(ctx context.Context, logger *zerolog.Logger, s *BusinessService, id int32) error {
	err := s.postgresService.Queries.DetachSieclesFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach siecles from mobilier image")
	}

	err = s.postgresService.Queries.DetachMediasFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach medias from mobilier image")
	}

	err = s.postgresService.Queries.DetachThemesFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach themes from mobilier image")
	}

	err = s.postgresService.Queries.DetachNaturesFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach natures from mobilier image")
	}

	err = s.postgresService.Queries.DetachEtatsFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach conservation states from mobilier image")
	}

	err = s.postgresService.Queries.DetachMateriauxFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach materials from mobilier image")
	}

	err = s.postgresService.Queries.DetachTechniquesFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach techniques from mobilier image")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink mobiliers images document from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink personnes morales document from mobilier image")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to unlink personnes physiques document from mobilier image")
	}

	err = s.postgresService.Queries.DetachAuthorFromMobImg(ctx, id)
	if err != nil {
		logger.Error().Err(err).Int32("id", id).Msg("failed to detach author from mobilier image")
	}

	if err = s.postgresService.Queries.DeleteMobilierImage(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		logger.Error().Err(err).Msg("failed to reject document")

		return catalogs.ErrUnexpectedError
	}

	return nil
}

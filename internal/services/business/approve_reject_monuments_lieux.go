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
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (b *BusinessService) ApproveRejectMonumentLieu(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	req, err := _http.DecodeAndValidateJSONBody[*api.PatchUsersBody](r)
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

	if err := processApproveRejectMonumentLieu(r.Context(), b, &token, req, int32(intID)); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type approveRejectMonumentLieuExchangeData struct {
	logger    *zerolog.Logger
	err       error
	id        int32
	params    *api.PatchUsersBody
	token     *jwt.SessionInfo
	documents *queries.GetMonumentLieuByIDRow
}

//nolint:lll
type approveRejectMonumentLieuState func(ctx context.Context, s *BusinessService, data *approveRejectMonumentLieuExchangeData) approveRejectMonumentLieuState

//nolint:lll
func processApproveRejectMonumentLieu(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.PatchUsersBody, id int32) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := approveRejectMonumentLieuExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := approveRejectMonumentLieuIfAdmin; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during personne physique approval")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

//nolint:lll
func approveRejectMonumentLieuIfAdmin(ctx context.Context, b *BusinessService, exData *approveRejectMonumentLieuExchangeData) approveRejectMonumentLieuState {
	if exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot approve or reject documents")
		exData.err = catalogs.ErrUserNotAdmin

		return nil
	}

	doc, err := b.postgresService.Queries.GetMonumentLieuByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Info().Msg("personne physique not found")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get personne physique by id")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.documents = &doc

	if *exData.params.Action == api.PatchUsersBodyActionActivate {
		return approveMonumentLieu
	}

	return rejectMonumentLieu
}

//nolint:lll
func approveMonumentLieu(ctx context.Context, s *BusinessService, exData *approveRejectMonumentLieuExchangeData) approveRejectMonumentLieuState {
	if err := s.postgresService.Queries.ValidatePendingMonumentLieu(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to approve document")
		exData.err = catalogs.ErrUnexpectedError
	}

	if exData.documents.ParentID.Valid {
		exData.id = exData.documents.ParentID.Int32

		return rejectMonumentLieu
	}

	return nil
}

//nolint:cyclop,lll
func rejectMonumentLieu(ctx context.Context, s *BusinessService, exData *approveRejectMonumentLieuExchangeData) approveRejectMonumentLieuState {
	err := s.postgresService.Queries.DetachSieclesFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach siecles from monument lieu")
	}

	err = s.postgresService.Queries.DetachMediasFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach medias from monument lieu")
	}

	err = s.postgresService.Queries.DetachThemesFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach themes from monument lieu")
	}

	err = s.postgresService.Queries.DetachNaturesFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach natures from monument lieu")
	}

	err = s.postgresService.Queries.DetachEtatsFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach conservation states from monument lieu")
	}

	err = s.postgresService.Queries.DetachMateriauxFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach materials from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink mobiliers images document from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromPersMo(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink personnes morales document from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink personnes physiques document from monument lieu")
	}

	err = s.postgresService.Queries.DetachAuthorFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach author from monument lieu")
	}

	if err := s.postgresService.Queries.DeletePendingMonumentLieu(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to approve document")
		exData.err = catalogs.ErrUnexpectedError
	}

	return nil
}

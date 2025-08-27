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

func (b *BusinessService) ApproveRejectPersonneMorale(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	if err := processApproveRejectPersonneMorale(r.Context(), b, &token, req, int32(intID)); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type approveRejectPersonneMoraleExchangeData struct {
	logger    *zerolog.Logger
	err       error
	id        int32
	params    *api.PatchUsersBody
	token     *jwt.SessionInfo
	documents *queries.GetPersonneMoraleByIDRow
}

//nolint:lll
type approveRejectPersonneMoraleState func(ctx context.Context, s *BusinessService, data *approveRejectPersonneMoraleExchangeData) approveRejectPersonneMoraleState

//nolint:lll
func processApproveRejectPersonneMorale(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.PatchUsersBody, id int32) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := approveRejectPersonneMoraleExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := approveRejectPersonneMoraleIfAdmin; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during personne morale approval")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

//nolint:lll
func approveRejectPersonneMoraleIfAdmin(ctx context.Context, b *BusinessService, exData *approveRejectPersonneMoraleExchangeData) approveRejectPersonneMoraleState {
	if exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot approve or reject documents")
		exData.err = catalogs.ErrUserNotAdmin

		return nil
	}

	doc, err := b.postgresService.Queries.GetPersonneMoraleByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Info().Msg("personne morale image not found")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get personne morale by id")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.documents = &doc

	if *exData.params.Action == api.PatchUsersBodyActionActivate {
		return approvePersonneMorale
	}

	return rejectPersonneMorale
}

//nolint:lll
func approvePersonneMorale(ctx context.Context, s *BusinessService, exData *approveRejectPersonneMoraleExchangeData) approveRejectPersonneMoraleState {
	if err := s.postgresService.Queries.ValidatePendingPersonneMorales(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to approve document")
		exData.err = catalogs.ErrUnexpectedError
	}

	if exData.documents.ParentID.Valid {
		exData.id = exData.documents.ParentID.Int32

		return rejectPersonneMorale
	}

	return nil
}

//nolint:cyclop,lll
func rejectPersonneMorale(ctx context.Context, s *BusinessService, exData *approveRejectPersonneMoraleExchangeData) approveRejectPersonneMoraleState {
	err := s.postgresService.Queries.DetachSieclesFromPersMo(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach siecles from personne morale")
	}

	err = s.postgresService.Queries.DetachMediasFromPersMo(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach medias from personne morale")
	}

	err = s.postgresService.Queries.DetachThemesFromPersMo(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach themes from personne morale")
	}

	err = s.postgresService.Queries.DetachNaturesFromPersMo(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach natures from personne morale")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink monuments lieux document from personne morale")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink mobiliers images document from personne morale")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink personnes physiques document from personne morale")
	}

	err = s.postgresService.Queries.DetachAuthorFromPersMo(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach author from personne morale")
	}

	if err := s.postgresService.Queries.DeletePendingPersonneMorale(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to reject document")
		exData.err = catalogs.ErrUnexpectedError
	}

	return nil
}

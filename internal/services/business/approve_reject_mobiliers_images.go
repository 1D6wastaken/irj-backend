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

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (b *BusinessService) ApproveRejectMobilierImage(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	if err := processApproveRejectMobilierImage(r.Context(), b, &token, req, int32(intID)); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type approveRejectMobilierImageExchangeData struct {
	logger    *zerolog.Logger
	err       error
	id        int32
	params    *api.PatchUsersBody
	token     *jwt.SessionInfo
	documents *queries.GetMobilierImageByIDRow
}

//nolint:lll
type approveRejectMobilierImageState func(ctx context.Context, s *BusinessService, data *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState

//nolint:lll
func processApproveRejectMobilierImage(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.PatchUsersBody, id int32) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := approveRejectMobilierImageExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := approveRejectMobilierImageIfAdmin; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during mobilier image approval")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

//nolint:lll
func approveRejectMobilierImageIfAdmin(ctx context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	if exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot approve or reject documents")
		exData.err = catalogs.ErrUserNotAdmin

		return nil
	}

	doc, err := s.postgresService.Queries.GetMobilierImageByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Info().Msg("mobilier image not found")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get mobilier image by id")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.documents = &doc

	if *exData.params.Action == api.PatchUsersBodyActionActivate {
		return approveMobilierImage
	}

	return rejectMobilierImage
}

//nolint:lll
func approveMobilierImage(ctx context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	if err := s.postgresService.Queries.ValidatePendingMobilierImage(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to approve document")
		exData.err = catalogs.ErrUnexpectedError
	}

	if exData.documents.ParentID.Valid {
		exData.id = exData.documents.ParentID.Int32

		return rejectMobilierImage
	}

	return storeMobImgDocumentValidationEvent
}

//nolint:lll
func storeMobImgDocumentValidationEvent(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, adminID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentValidationEvent(context.Background(), queries.DocumentValidationEventParams{
			AdminID: pgtype.Text{
				String: adminID,
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
			logger.Error().Err(err).Msg("failed to store document validation event")
		}
	}(exData.logger, exData.token.ID, exData.id)

	return nil
}

//nolint:cyclop,lll
func rejectMobilierImage(ctx context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	err := s.postgresService.Queries.DetachSieclesFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach siecles from mobilier image")
	}

	err = s.postgresService.Queries.DetachMediasFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach medias from mobilier image")
	}

	err = s.postgresService.Queries.DetachThemesFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach themes from mobilier image")
	}

	err = s.postgresService.Queries.DetachNaturesFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach natures from mobilier image")
	}

	err = s.postgresService.Queries.DetachEtatsFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach conservation states from mobilier image")
	}

	err = s.postgresService.Queries.DetachMateriauxFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach materials from mobilier image")
	}

	err = s.postgresService.Queries.DetachTechniquesFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach techniques from mobilier image")
	}

	err = s.postgresService.Queries.UnlinkMonuLieuFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink mobiliers images document from monument lieu")
	}

	err = s.postgresService.Queries.UnlinkPersMoFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink personnes morales document from mobilier image")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink personnes physiques document from mobilier image")
	}

	err = s.postgresService.Queries.DetachAuthorFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach author from mobilier image")
	}

	if err = s.postgresService.Queries.DeletePendingMobilierImage(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to reject document")
		exData.err = catalogs.ErrUnexpectedError
	}

	return storeMobImgDocumentRejectionEvent
}

//nolint:lll
func storeMobImgDocumentRejectionEvent(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, adminID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentRejectionEvent(context.Background(), queries.DocumentRejectionEventParams{
			AdminID: pgtype.Text{
				String: adminID,
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
			logger.Error().Err(err).Msg("failed to store document rejection event")
		}
	}(exData.logger, exData.token.ID, exData.id)

	return nil
}

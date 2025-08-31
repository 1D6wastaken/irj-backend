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

func (b *BusinessService) ApproveRejectPersonnePhysique(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	if err := processApproveRejectPersonnePhysique(r.Context(), b, &token, req, int32(intID)); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type approveRejectPersonnePhysiqueExchangeData struct {
	logger    *zerolog.Logger
	err       error
	id        int32
	params    *api.PatchUsersBody
	token     *jwt.SessionInfo
	documents *queries.GetPersonnePhysiqueByIDRow
}

//nolint:lll
type approveRejectPersonnePhysiqueState func(ctx context.Context, s *BusinessService, data *approveRejectPersonnePhysiqueExchangeData) approveRejectPersonnePhysiqueState

//nolint:lll
func processApproveRejectPersonnePhysique(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.PatchUsersBody, id int32) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := approveRejectPersonnePhysiqueExchangeData{
		logger: logger,
		id:     id,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := approveRejectPersonnePhysiqueIfAdmin; state != nil; {
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
func approveRejectPersonnePhysiqueIfAdmin(ctx context.Context, b *BusinessService, exData *approveRejectPersonnePhysiqueExchangeData) approveRejectPersonnePhysiqueState {
	if exData.token.Grade != string(queries.UserGradeADMIN) {
		exData.logger.Warn().Msg("user is not an admin and therefore cannot approve or reject documents")
		exData.err = catalogs.ErrUserNotAdmin

		return nil
	}

	doc, err := b.postgresService.Queries.GetPersonnePhysiqueByID(ctx, exData.id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Info().Msg("personne physique image not found")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to get personne physique by id")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.documents = &doc

	if *exData.params.Action == api.PatchUsersBodyActionActivate {
		return approvePersonnePhysique
	}

	return rejectPersonnePhysique
}

//nolint:lll
func approvePersonnePhysique(ctx context.Context, s *BusinessService, exData *approveRejectPersonnePhysiqueExchangeData) approveRejectPersonnePhysiqueState {
	if err := s.postgresService.Queries.ValidatePendingPersonnePhysique(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to approve document")
		exData.err = catalogs.ErrUnexpectedError
	}

	if exData.documents.ParentID.Valid {
		exData.id = exData.documents.ParentID.Int32

		return rejectPersonnePhysique
	}

	return storePersPhyDocumentValidationEvent
}

//nolint:lll
func storePersPhyDocumentValidationEvent(_ context.Context, s *BusinessService, exData *approveRejectPersonnePhysiqueExchangeData) approveRejectPersonnePhysiqueState {
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
				String: "personnes_physiques",
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
func rejectPersonnePhysique(ctx context.Context, s *BusinessService, exData *approveRejectPersonnePhysiqueExchangeData) approveRejectPersonnePhysiqueState {
	err := s.postgresService.Queries.DetachSieclesFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach siecles from personne physique")
	}

	err = s.postgresService.Queries.DetachMediasFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach medias from personne physique")
	}

	err = s.postgresService.Queries.DetachThemesFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach themes from personne physique")
	}

	err = s.postgresService.Queries.DetachHistoricalPeriodsFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach historical periods from personne physique")
	}

	err = s.postgresService.Queries.DetachProfessionsFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach professions from personne physique")
	}

	err = s.postgresService.Queries.DetachModeDeTransportsFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach travel modes from personne physique")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromMonuLieu(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink monuments lieux document from personne physique")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromMobImg(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink mobiliers images document from personne physique")
	}

	err = s.postgresService.Queries.UnlinkPersPhyFromPersMo(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to unlink personnes morales document from personne physique")
	}

	err = s.postgresService.Queries.DetachAuthorFromPersPhy(ctx, exData.id)
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to detach author from personne physique")
	}

	if err := s.postgresService.Queries.DeletePendingPersonnePhysique(ctx, exData.id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exData.logger.Warn().Msg("document not found and therefore can not be approved")

			return nil
		}

		exData.logger.Error().Err(err).Msg("failed to approve document")
		exData.err = catalogs.ErrUnexpectedError
	}

	return storePersPhyDocumentRejectionEvent
}

//nolint:lll
func storePersPhyDocumentRejectionEvent(_ context.Context, s *BusinessService, exData *approveRejectPersonnePhysiqueExchangeData) approveRejectPersonnePhysiqueState {
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
				String: "personnes_physiques",
				Valid:  true,
			},
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to store document rejection event")
		}
	}(exData.logger, exData.token.ID, exData.id)

	return nil
}

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
	"irj/internal/smtp"
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
		if err := deleteMobilierImage(ctx, exData.logger, s, exData.documents.ParentID.Int32); err != nil {
			exData.err = err
		}

		return storeMobImgDocumentUpdateValidationEvent
	}

	return storeMobImgDocumentSubmissionValidationEvent
}

//nolint:lll
func storeMobImgDocumentSubmissionValidationEvent(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, adminID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentSubmissionValidationEvent(context.Background(), queries.DocumentSubmissionValidationEventParams{
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

	return sendApprovalEmailMobImg
}

//nolint:lll
func storeMobImgDocumentUpdateValidationEvent(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, adminID string, documentID, parentID int32) {
		defer s.stopper.Release()

		if err := s.postgresService.Queries.DocumentUpdateValidationEvent(context.Background(), queries.DocumentUpdateValidationEventParams{
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
		}); err != nil {
			logger.Error().Err(err).Msg("failed to store document validation event")
		}

		if err := s.postgresService.Queries.UpdateDocumentIDAfterUpdate(context.Background(), queries.UpdateDocumentIDAfterUpdateParams{
			NewDocID: pgtype.Int4{
				Int32: documentID,
				Valid: true,
			},
			ParentDocID: pgtype.Int4{
				Int32: parentID,
				Valid: true,
			},
		}); err != nil {
			logger.Error().Err(err).Msg("failed to store document validation event")
		}
	}(exData.logger, exData.token.ID, exData.id, exData.documents.ParentID.Int32)

	return sendApprovalEmailMobImg
}

//nolint:lll
func rejectMobilierImage(ctx context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	if err := deleteMobilierImage(ctx, exData.logger, s, exData.id); err != nil {
		exData.err = err

		return nil
	}

	if exData.documents.ParentID.Valid {
		return storeMobImgDocumentUpdateRejectionEvent
	}

	return storeMobImgDocumentSubmissionRejectionEvent
}

//nolint:lll
func storeMobImgDocumentSubmissionRejectionEvent(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, adminID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentSubmissionRejectionEvent(context.Background(), queries.DocumentSubmissionRejectionEventParams{
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

	return sendRejectionEmailMobImg
}

//nolint:lll
func storeMobImgDocumentUpdateRejectionEvent(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, adminID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentUpdateRejectionEvent(context.Background(), queries.DocumentUpdateRejectionEventParams{
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

	return sendRejectionEmailMobImg
}

//nolint:lll
func sendApprovalEmailMobImg(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	if !exData.documents.UserID.Valid {
		return nil
	}

	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, userID string, id int32, title string, isUpdate bool) {
		defer s.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		user, err := s.postgresService.Queries.GetUserByID(ctx, userID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get user for approval email")

			return
		}

		to := []smtp.EmailPerson{{Name: user.Prenom + " " + user.Nom, Email: user.Email}}

		if err := s.smtpService.SendDocumentApprovedMail(ctx, to, smtp.SourceMobiliersImages, id, isUpdate, title); err != nil {
			logger.Error().Err(err).Msg("failed to send approval email")
		}
	}(exData.logger, exData.documents.UserID.String, exData.id, exData.documents.Title, exData.documents.ParentID.Valid)

	return nil
}

//nolint:lll
func sendRejectionEmailMobImg(_ context.Context, s *BusinessService, exData *approveRejectMobilierImageExchangeData) approveRejectMobilierImageState {
	if !exData.documents.UserID.Valid {
		return nil
	}

	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, userID string, id int32, title string, isUpdate bool) {
		defer s.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		user, err := s.postgresService.Queries.GetUserByID(ctx, userID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get user for rejection email")

			return
		}

		to := []smtp.EmailPerson{{Name: user.Prenom + " " + user.Nom, Email: user.Email}}

		if err := s.smtpService.SendDocumentRejectedMail(ctx, to, isUpdate, title, id); err != nil {
			logger.Error().Err(err).Msg("failed to send rejection email")
		}
	}(exData.logger, exData.documents.UserID.String, exData.id, exData.documents.Title, exData.documents.ParentID.Valid)

	return nil
}

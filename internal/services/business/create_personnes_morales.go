package business

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

func (b *BusinessService) CreatePersonneMorale(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	req, err := _http.DecodeAndValidateJSONBody[*api.PersonneMoraleCreationBody](r)
	if err != nil {
		return _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	if err := processCreatePersonneMorale(r.Context(), b, &token, req); err != nil {
		status, body := catalogs.GetError(err)

		return _http.WriteJSONResponse(w, status, body)
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

type createPersonneMoraleExchangeData struct {
	logger *zerolog.Logger
	err    error
	params *api.PersonneMoraleCreationBody
	token  *jwt.SessionInfo
	id     int32
}

//nolint:lll
type createPersonneMoraleState func(ctx context.Context, s *BusinessService, data *createPersonneMoraleExchangeData) createPersonneMoraleState

//nolint:lll
func processCreatePersonneMorale(ctx context.Context, s *BusinessService, token *jwt.SessionInfo, req *api.PersonneMoraleCreationBody) error {
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeOut)

	defer cancel()

	exData := createPersonneMoraleExchangeData{
		logger: logger,
		params: req,
		token:  token,
	}

	respChan := make(chan error, 1)

	s.stopper.Hold(1)

	go func() {
		defer s.stopper.Release()

		for state := createPersonneMorale; state != nil; {
			state = state(ctx, s, &exData)
		}

		respChan <- exData.err

		close(respChan)
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("deadline was reached during mobilier image creation")

		return catalogs.ErrServerTimeout
	case resp := <-respChan:
		return resp
	}
}

func createPersonneMorale(ctx context.Context, s *BusinessService, exData *createPersonneMoraleExchangeData) createPersonneMoraleState {
	id, err := s.postgresService.Queries.CreatePersMorale(ctx, queries.CreatePersMoraleParams{
		Title:      pgtype.Text{String: *exData.params.Title, Valid: true},
		Comment:    pgtype.Text{String: exData.params.Comment, Valid: exData.params.Comment != ""},
		Historique: pgtype.Text{String: exData.params.History, Valid: exData.params.History != ""},
		ActeFondation: pgtype.Bool{
			Bool:  exData.params.FoundationDeed,
			Valid: true,
		},
		SimpleMention: pgtype.Bool{
			Bool:  exData.params.SimpleMention,
			Valid: true,
		},
		TexteStatuts: pgtype.Bool{
			Bool:  exData.params.StatusText,
			Valid: true,
		},
		Fonctionnement:      pgtype.Text{String: exData.params.Functioning, Valid: exData.params.Functioning != ""},
		ParticipationVieSoc: pgtype.Text{String: exData.params.SocialInvolvement, Valid: exData.params.SocialInvolvement != ""},
		Objets:              pgtype.Text{String: exData.params.LinkedObjects, Valid: exData.params.LinkedObjects != ""},
		Bibliographie:       pgtype.Text{String: exData.params.Bibliography, Valid: exData.params.Bibliography != ""},
		Sources:             pgtype.Text{String: exData.params.Source, Valid: exData.params.Source != ""},
		Contributeurs:       pgtype.Text{String: strings.Join(exData.params.Contributors, ","), Valid: len(exData.params.Contributors) > 0},
		IDCommune: pgtype.Int4{
			Int32: exData.params.City,
			Valid: exData.params.City != 0,
		},
		IDPays: pgtype.Int4{
			Int32: exData.params.Country,
			Valid: exData.params.Country != 0,
		},
	})
	if err != nil {
		exData.logger.Error().Err(err).Msg("failed to insert personne morale")
		exData.err = catalogs.ErrUnexpectedError

		return nil
	}

	exData.logger.Info().Int32("id", id).Msg("personne morale created")

	exData.id = id

	return storePersMoDocumentSubmissionEvent
}

//nolint:lll
func storePersMoDocumentSubmissionEvent(_ context.Context, s *BusinessService, exData *createPersonneMoraleExchangeData) createPersonneMoraleState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(logger *zerolog.Logger, userID string, documentID int32) {
		defer s.stopper.Release()

		err := s.postgresService.Queries.DocumentSubmissionEvent(context.Background(), queries.DocumentSubmissionEventParams{
			UserID: userID,
			DocumentID: pgtype.Int4{
				Int32: documentID,
				Valid: true,
			},
			Comment: pgtype.Text{
				String: "personnes_morales",
				Valid:  true,
			},
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to store document submission event")
		}
	}(exData.logger, exData.token.ID, exData.id)

	return linkPersonneMorale
}

func linkPersonneMorale(ctx context.Context, s *BusinessService, exData *createPersonneMoraleExchangeData) createPersonneMoraleState {
	err := s.postgresService.Queries.AttachSieclesToPersMo(ctx, queries.AttachSieclesToPersMoParams{
		SiecleID: exData.params.Centuries,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach centuries to personne morale")
	}

	err = s.postgresService.Queries.AttachMediasToPersMo(ctx, queries.AttachMediasToPersMoParams{
		MediaIds: exData.params.Medias,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach medias to personne morale")
	}

	err = s.postgresService.Queries.AttachThemesToPersMo(ctx, queries.AttachThemesToPersMoParams{
		ThemeIds: exData.params.Themes,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach themes to personne morale")
	}

	err = s.postgresService.Queries.AttachNaturesToPersMo(ctx, queries.AttachNaturesToPersMoParams{
		NatureIds: exData.params.Natures,
		ID:        exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach natures to personne morale")
	}

	err = s.postgresService.Queries.LinkPersMoToMonuLieu(ctx, queries.LinkPersMoToMonuLieuParams{
		MonuLieuIds: exData.params.LinkedMonumentsLieux,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link monuments lieux document to personne morale")
	}

	err = s.postgresService.Queries.LinkPersMoToMobImg(ctx, queries.LinkPersMoToMobImgParams{
		MobImgIds: exData.params.LinkedMobiliersImages,
		ID:        exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link mobiliers images document to personne morale")
	}

	err = s.postgresService.Queries.LinkPersMoToPersPhy(ctx, queries.LinkPersMoToPersPhyParams{
		PersoPhyIds: exData.params.LinkedPersPhysiques,
		ID:          exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to link personnes physiques document to personne morale")
	}

	return addAuteurToPersonneMorale
}

//nolint:lll
func addAuteurToPersonneMorale(ctx context.Context, s *BusinessService, exData *createPersonneMoraleExchangeData) createPersonneMoraleState {
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

	err = s.postgresService.Queries.AttachAuthorToPersMo(ctx, queries.AttachAuthorToPersMoParams{
		AuteurID: id,
		ID:       exData.id,
	})
	if err != nil {
		exData.logger.Error().Err(err).Int32("id", exData.id).Msg("failed to attach author to mobilier image")
	}

	return sendNewPersonneMoraleDocumentMail
}

//nolint:lll
func sendNewPersonneMoraleDocumentMail(_ context.Context, s *BusinessService, exData *createPersonneMoraleExchangeData) createPersonneMoraleState {
	s.stopper.Hold(1)

	//nolint:contextcheck
	go func(s *BusinessService, logger *zerolog.Logger) {
		defer s.stopper.Release()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
		defer cancel()

		row, err := s.postgresService.Queries.GetUsersByGrade(ctx, queries.UserGradeADMIN)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get admin user")

			return
		}

		to := make([]smtp.EmailPerson, 0, len(row))

		//nolint:gocritic
		for _, user := range row {
			to = append(to, smtp.EmailPerson{
				Name:  user.Prenom + " " + user.Nom,
				Email: user.Email,
			})
		}

		_ = s.smtpService.SendNewDocumentMail(ctx, to)
	}(s, exData.logger)

	return nil
}

package business

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	"irj/pkg/collections"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func (b *BusinessService) GetPendingPersonnesPhysiques(w http.ResponseWriter, r *http.Request) *_http.APIError {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	logger := zerolog.Ctx(ctx)

	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	if token.Grade != string(queries.UserGradeADMIN) {
		logger.Warn().Msg("user is not an admin and therefore cannot see pending documents")

		status, body := catalogs.GetError(catalogs.ErrUserNotAdmin)

		return _http.WriteJSONResponse(w, status, body)
	}

	persPhy, err := b.postgresService.Queries.GetPendingPersonnesPhysiques(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get pending personnes physiques")

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	items := make([]*api.PendingDocuments, 0, len(persPhy))

	for i := range persPhy {
		row := persPhy[i]

		var (
			medias      []*api.Media
			commune     string
			departement string
			region      string
			pays        string
		)

		nocoMedias, err := b.parseMedias(row.Medias)
		if err == nil {
			medias = collections.Map(nocoMedias, func(m NocoMedia) *api.Media {
				return &api.Media{
					ID:    &m.ID,
					Title: &m.Title,
				}
			})

			// sort medias by Title
			slices.SortFunc(medias, func(a, b *api.Media) int {
				return strings.Compare(*a.Title, *b.Title)
			})
		} else {
			logger.Error().Err(err).Msg("failed to parse medias")
		}

		if row.Commune != nil {
			commune = row.Commune.(string)
		}

		if row.Departement != nil {
			departement = row.Departement.(string)
		}

		if row.Region != nil {
			region = row.Region.(string)
		}

		if row.Pays != nil {
			pays = row.Pays.(string)
		}

		items = append(items, &api.PendingDocuments{
			ID:           &row.ID,
			Title:        &row.Firstname.String,
			CreationDate: utils.PtrTo(strfmt.Date(row.DateCreation.Time)),
			Authors:      collections.InterfaceToStringSlice(row.Redacteurs),
			City:         commune,
			Department:   departement,
			Region:       region,
			Country:      pays,
			Professions:  collections.InterfaceToStringSlice(row.Professions),
			Medias:       medias,
			Centuries:    collections.InterfaceToStringSlice(row.Siecles),
			ParentID:     row.ParentID.Int32,
		})
	}

	return _http.WriteJSONResponse(w, http.StatusOK, items)
}

func (b *BusinessService) GetDraftPersonnesPhysiques(w http.ResponseWriter, r *http.Request) *_http.APIError {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	logger := zerolog.Ctx(ctx)

	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	persPhy, err := b.postgresService.Queries.GetDraftPersonnesPhysiques(ctx, pgtype.Text{
		String: token.ID,
		Valid:  true,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get draft personnes physiques")

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	items := make([]*api.PendingDocuments, 0, len(persPhy))

	for i := range persPhy {
		row := persPhy[i]

		var (
			medias      []*api.Media
			commune     string
			departement string
			region      string
			pays        string
		)

		nocoMedias, err := b.parseMedias(row.Medias)
		if err == nil {
			medias = collections.Map(nocoMedias, func(m NocoMedia) *api.Media {
				return &api.Media{
					ID:    &m.ID,
					Title: &m.Title,
				}
			})
		} else {
			logger.Error().Err(err).Msg("failed to parse medias")
		}

		if row.Commune != nil {
			commune = row.Commune.(string)
		}

		if row.Departement != nil {
			departement = row.Departement.(string)
		}

		if row.Region != nil {
			region = row.Region.(string)
		}

		if row.Pays != nil {
			pays = row.Pays.(string)
		}

		items = append(items, &api.PendingDocuments{
			ID:           &row.ID,
			Title:        &row.Firstname.String,
			CreationDate: utils.PtrTo(strfmt.Date(row.DateCreation.Time)),
			Authors:      collections.InterfaceToStringSlice(row.Redacteurs),
			City:         commune,
			Department:   departement,
			Region:       region,
			Country:      pays,
			Professions:  collections.InterfaceToStringSlice(row.Professions),
			Medias:       medias,
			Centuries:    collections.InterfaceToStringSlice(row.Siecles),
			ParentID:     row.ParentID.Int32,
		})
	}

	return _http.WriteJSONResponse(w, http.StatusOK, items)
}

func (b *BusinessService) GetPersonnePhysique(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	params := httprouter.ParamsFromContext(subCtx)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid").Err(err)
	}

	ppy, err := b.postgresService.Queries.GetPersonnePhysiqueByID(subCtx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _http.ErrNotFound.Msg("personne physique not found").Err(err)
		}

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	var medias []*api.Media

	nocoMedias, err := b.parseMedias(ppy.Medias)
	if err == nil {
		medias = collections.Map(nocoMedias, func(m NocoMedia) *api.Media {
			return &api.Media{
				ID:    &m.ID,
				Title: &m.Title,
			}
		})
	}

	city, department, region, country, err := parseLocation(ppy.City, ppy.Department, ppy.Region, ppy.Country)
	if err != nil {
		return _http.ErrInternalError.Msg("error while parsing location").Err(err)
	}

	return _http.WriteJSONResponse(w, http.StatusOK, api.PersonnePhysique{
		ID:                    ppy.ID,
		Firstname:             ppy.Firstname.String,
		Birthdate:             ppy.DateNaissance.String,
		Death:                 ppy.DateDeces.String,
		Attestation:           ppy.Attestation.String,
		HistoricalPeriod:      interfaceSliceToBasicFilterSlice(ppy.HistoricalPeriod.([]any)),
		Bibliography:          ppy.Bibliographie.String,
		BiographicalElements:  ppy.ElementsBiographiques.String,
		PilgrimageElement:     ppy.ElementsPelerinage.String,
		Commutation:           ppy.CommutationVoeu.String,
		Sources:               ppy.Sources.String,
		CreationDate:          strfmt.Date(ppy.DateCreation.Time),
		UpdateDate:            strfmt.Date(ppy.DateMaj.Time),
		Published:             ppy.Publie.Bool,
		Contributors:          ppy.Contributeurs.String,
		Comment:               ppy.Commentaires.String,
		Authors:               interfaceSliceToBasicFilterSlice(ppy.Authors.([]any)),
		City:                  city,
		Department:            department,
		Region:                region,
		Country:               country,
		Travels:               interfaceSliceToBasicFilterSlice(ppy.Travels.([]any)),
		Professions:           interfaceSliceToBasicFilterSlice(ppy.Professions.([]any)),
		EventNature:           ppy.NatureEvenement.String,
		Medias:                medias,
		Centuries:             interfaceSliceToBasicFilterSlice(ppy.Centuries.([]any)),
		LinkedMonumentsPlaces: collections.InterfaceToInt32Slice(ppy.MonumentsLieuxLiees),
		LinkedLegalEntities:   collections.InterfaceToInt32Slice(ppy.PersonnesMoralesLiees),
		LinkedFurnitureImages: collections.InterfaceToInt32Slice(ppy.MobiliersImagesLiees),
		Themes:                interfaceSliceToBasicFilterSlice(ppy.Themes.([]any)),
		ParentID:              ppy.ParentID.Int32,
		Historiography:        ppy.Historiographie.String,
		Evenements:            ppy.Evenements.String,
		Preparatifs:           ppy.Preparatifs.String,
		CheminSuivi:           ppy.CheminSuivi.String,
		Arrivee:               ppy.Arrivee.String,
		Retour:                ppy.Retour.String,
		NonExecution:          ppy.NonExecution.String,
		Age:                   ppy.Age.String,
		CompositionGroupe:     ppy.CompositionGroupe.String,
	})
}

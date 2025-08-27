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
	"irj/pkg/collections"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/go-openapi/strfmt"
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
		logger.Error().Err(err).Msg("failed to get pending personnes morales")

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

	pmo, err := b.postgresService.Queries.GetPersonnePhysiqueByID(subCtx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _http.ErrNotFound.Msg("personne physique not found").Err(err)
		}

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	var medias []*api.Media

	nocoMedias, err := b.parseMedias(pmo.Medias)
	if err == nil {
		medias = collections.Map(nocoMedias, func(m NocoMedia) *api.Media {
			return &api.Media{
				ID:    &m.ID,
				Title: &m.Title,
			}
		})
	}

	return _http.WriteJSONResponse(w, http.StatusOK, api.PersonnePhysique{
		ID:                    pmo.ID,
		Firstname:             pmo.Firstname.String,
		Birthdate:             pmo.DateNaissance.String,
		Death:                 pmo.DateDeces.String,
		Attestation:           pmo.Attestation.String,
		HistoricalPeriod:      collections.InterfaceToStringSlice(pmo.HistoricalPeriod),
		Bibliography:          pmo.Bibliographie.String,
		BiographicalElements:  pmo.ElementsBiographiques.String,
		PilgrimageElement:     pmo.ElementsPelerinage.String,
		Commutation:           pmo.CommutationVoeu.String,
		Sources:               pmo.Sources.String,
		CreationDate:          strfmt.Date(pmo.DateCreation.Time),
		UpdateDate:            strfmt.Date(pmo.DateMaj.Time),
		Published:             pmo.Publie.Bool,
		Contributors:          pmo.Contributeurs.String,
		Comment:               pmo.Commentaires.String,
		Authors:               collections.InterfaceToStringSlice(pmo.Redacteurs),
		City:                  pmo.Commune.String,
		Department:            pmo.Departement.String,
		Region:                pmo.Region.String,
		Country:               pmo.Pays.String,
		Travels:               collections.InterfaceToStringSlice(pmo.Travels),
		Professions:           collections.InterfaceToStringSlice(pmo.Professions),
		EventNature:           pmo.NatureEvenement.String,
		Medias:                medias,
		Centuries:             collections.InterfaceToStringSlice(pmo.Siecles),
		LinkedMonumentsPlaces: collections.InterfaceToInt32Slice(pmo.MonumentsLieuxLiees),
		LinkedLegalEntities:   collections.InterfaceToInt32Slice(pmo.PersonnesMoralesLiees),
		LinkedFurnitureImages: collections.InterfaceToInt32Slice(pmo.MobiliersImagesLiees),
		Themes:                collections.InterfaceToStringSlice(pmo.Themes),
	})
}

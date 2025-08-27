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

func (b *BusinessService) GetPendingPersonnesMorales(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	persMo, err := b.postgresService.Queries.GetPendingPersonnesMorales(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get pending personnes morales")

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	items := make([]*api.PendingDocuments, 0, len(persMo))

	for i := range persMo {
		row := persMo[i]

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
			Title:        &row.Title.String,
			CreationDate: utils.PtrTo(strfmt.Date(row.DateCreation.Time)),
			Authors:      collections.InterfaceToStringSlice(row.Redacteurs),
			City:         commune,
			Department:   departement,
			Region:       region,
			Country:      pays,
			Natures:      collections.InterfaceToStringSlice(row.Natures),
			Medias:       medias,
			Centuries:    collections.InterfaceToStringSlice(row.Siecles),
		})
	}

	return _http.WriteJSONResponse(w, http.StatusOK, items)
}

func (b *BusinessService) GetPersonneMorale(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	params := httprouter.ParamsFromContext(subCtx)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid").Err(err)
	}

	pmo, err := b.postgresService.Queries.GetPersonneMoraleByID(subCtx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _http.ErrNotFound.Msg("personne morale not found").Err(err)
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

	return _http.WriteJSONResponse(w, http.StatusOK, api.PersonneMorale{
		ID:                    pmo.ID,
		Title:                 pmo.Title.String,
		FoundationDeed:        pmo.FoundationDeed.Bool,
		History:               pmo.Historique.String,
		Bibliography:          pmo.Bibliographie.String,
		SimpleMention:         pmo.SimpleMention.Bool,
		Process:               pmo.Fonctionnement.String,
		SocialInvolvement:     pmo.ParticipationVieSoc.String,
		Objects:               pmo.Objets.String,
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
		Natures:               collections.InterfaceToStringSlice(pmo.Natures),
		Medias:                medias,
		Centuries:             collections.InterfaceToStringSlice(pmo.Siecles),
		LinkedMonumentsPlaces: collections.InterfaceToInt32Slice(pmo.MonumentsLieuxLiees),
		LinkedIndividuals:     collections.InterfaceToInt32Slice(pmo.PersonnesPhysiquesLiees),
		LinkedFurnitureImages: collections.InterfaceToInt32Slice(pmo.MobiliersImagesLiees),
		Themes:                collections.InterfaceToStringSlice(pmo.Themes),
	})
}

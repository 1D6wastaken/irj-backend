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

func (b *BusinessService) GetPendingMonumentsLieux(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	monuments, err := b.postgresService.Queries.GetPendingMonumentsLieux(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get pending monuments lieux")

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	items := make([]*api.PendingDocuments, 0, len(monuments))

	for i := range monuments {
		row := monuments[i]

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
			Title:        &row.Title,
			CreationDate: utils.PtrTo(strfmt.Date(row.DateCreation.Time)),
			Authors:      collections.InterfaceToStringSlice(row.Redacteurs),
			City:         commune,
			Department:   departement,
			Region:       region,
			Country:      pays,
			Natures:      collections.InterfaceToStringSlice(row.Natures),
			Medias:       medias,
			Centuries:    collections.InterfaceToStringSlice(row.Siecles),
			ParentID:     row.ParentID.Int32,
		})
	}

	return _http.WriteJSONResponse(w, http.StatusOK, items)
}

func (b *BusinessService) GetDraftMonumentsLieux(w http.ResponseWriter, r *http.Request) *_http.APIError {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	logger := zerolog.Ctx(ctx)

	token, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	monuments, err := b.postgresService.Queries.GetDraftMonumentsLieux(ctx, pgtype.Text{
		String: token.ID,
		Valid:  true,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get draft monuments lieux")

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	items := make([]*api.PendingDocuments, 0, len(monuments))

	for i := range monuments {
		row := monuments[i]

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
			Title:        &row.Title,
			CreationDate: utils.PtrTo(strfmt.Date(row.DateCreation.Time)),
			Authors:      collections.InterfaceToStringSlice(row.Redacteurs),
			City:         commune,
			Department:   departement,
			Region:       region,
			Country:      pays,
			Natures:      collections.InterfaceToStringSlice(row.Natures),
			Medias:       medias,
			Centuries:    collections.InterfaceToStringSlice(row.Siecles),
			ParentID:     row.ParentID.Int32,
		})
	}

	return _http.WriteJSONResponse(w, http.StatusOK, items)
}

func (b *BusinessService) GetMonumentLieu(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	params := httprouter.ParamsFromContext(subCtx)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid").Err(err)
	}

	monument, err := b.postgresService.Queries.GetMonumentLieuByID(subCtx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _http.ErrNotFound.Msg("monument lieu not found").Err(err)
		}

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	var medias []*api.Media

	nocoMedias, err := b.parseMedias(monument.Medias)
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
	}

	city, department, region, country, err := parseLocation(monument.City, monument.Department, monument.Region, monument.Country)
	if err != nil {
		return _http.ErrInternalError.Msg("error while parsing location").Err(err)
	}

	return _http.WriteJSONResponse(w, http.StatusOK, api.MonumentLieu{
		ID:                    monument.ID,
		Title:                 monument.Title,
		Description:           monument.Description.String,
		History:               monument.Histoire.String,
		Geolocation:           monument.Geolocalisation.String,
		Bibliography:          monument.Bibliographie.String,
		Sources:               monument.Source.String,
		CreationDate:          strfmt.Date(monument.DateCreation.Time),
		UpdateDate:            strfmt.Date(monument.DateMaj.Time),
		Published:             monument.Publie.Bool,
		Contributors:          monument.Contributeurs.String,
		Protected:             monument.Protection.Bool,
		ProtectionComment:     monument.ProtectionCommentaires.String,
		Authors:               interfaceSliceToBasicFilterSlice(monument.Authors.([]any)),
		City:                  city,
		Department:            department,
		Region:                region,
		Country:               country,
		Conservation:          interfaceSliceToBasicFilterSlice(monument.Conservation.([]any)),
		Materials:             interfaceSliceToBasicFilterSlice(monument.Materials.([]any)),
		Natures:               interfaceSliceToBasicFilterSlice(monument.Natures.([]any)),
		Medias:                medias,
		Centuries:             interfaceSliceToBasicFilterSlice(monument.Centuries.([]any)),
		LinkedFurnitureImages: collections.InterfaceToInt32Slice(monument.MobiliersImagesLiees),
		LinkedIndividuals:     collections.InterfaceToInt32Slice(monument.PersonnesPhysiquesLiees),
		LinkedLegalEntities:   collections.InterfaceToInt32Slice(monument.PersonnesMoralesLiees),
		Themes:                interfaceSliceToBasicFilterSlice(monument.Themes.([]any)),
		ParentID:              monument.ParentID.Int32,
		Dimensions:            monument.Dimensions.String,
		Altitude:              monument.Altitude.String,
		Emplacement:           monument.Emplacement.String,
		DateConstruction:      monument.DateConstruction.String,
		PremiereMention:       monument.PremiereMention.String,
		ProprietaireActuel:    monument.ProprietaireActuel.String,
		Architecte:            monument.Architecte.String,
		Commanditaire:         monument.Commanditaire.String,
	})
}

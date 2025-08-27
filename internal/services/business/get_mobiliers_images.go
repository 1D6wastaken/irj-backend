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

func (b *BusinessService) GetPendingMobiliersImages(w http.ResponseWriter, r *http.Request) *_http.APIError {
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

	mobiliers, err := b.postgresService.Queries.GetPendingMobiliersImages(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get pending mobiliers images")

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	items := make([]*api.PendingDocuments, 0, len(mobiliers))

	for i := range mobiliers {
		row := mobiliers[i]

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
			CreationDate: utils.PtrTo(strfmt.Date(row.DateCrAtion.Time)),
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

func (b *BusinessService) GetMobilierImage(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	params := httprouter.ParamsFromContext(subCtx)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid").Err(err)
	}

	mobilier, err := b.postgresService.Queries.GetMobilierImageByID(subCtx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _http.ErrNotFound.Msg("mobilier image not found").Err(err)
		}

		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	var medias []*api.Media

	nocoMedias, err := b.parseMedias(mobilier.Medias)
	if err == nil {
		medias = collections.Map(nocoMedias, func(m NocoMedia) *api.Media {
			return &api.Media{
				ID:    &m.ID,
				Title: &m.Title,
			}
		})
	}

	return _http.WriteJSONResponse(w, http.StatusOK, api.MobilierImage{
		ID:                    mobilier.ID,
		Title:                 mobilier.Title,
		Description:           mobilier.Description.String,
		History:               mobilier.Historique.String,
		Bibliography:          mobilier.Bibliographie.String,
		Inscriptions:          mobilier.Inscriptions.String,
		CreationDate:          strfmt.Date(mobilier.DateCrAtion.Time),
		UpdateDate:            strfmt.Date(mobilier.DateMaj.Time),
		Published:             mobilier.Publie.Bool,
		Contributors:          mobilier.Contributeurs.String,
		Sources:               mobilier.Source.String,
		Protected:             mobilier.Protection.Bool,
		ProtectionComment:     mobilier.ProtectionCommentaires.String,
		ConservationPlace:     mobilier.LieuConservation.String,
		OriginPlace:           mobilier.LieuOrigine.String,
		Authors:               collections.InterfaceToStringSlice(mobilier.Redacteurs),
		City:                  mobilier.Commune.String,
		Department:            mobilier.Departement.String,
		Region:                mobilier.Region.String,
		Country:               mobilier.Pays.String,
		Conservation:          collections.InterfaceToStringSlice(mobilier.EtatsConservation),
		Materials:             collections.InterfaceToStringSlice(mobilier.Materiaux),
		Natures:               collections.InterfaceToStringSlice(mobilier.Natures),
		Medias:                medias,
		Centuries:             collections.InterfaceToStringSlice(mobilier.Siecles),
		Techniques:            collections.InterfaceToStringSlice(mobilier.Techniques),
		LinkedMonumentsPlaces: collections.InterfaceToInt32Slice(mobilier.MonumentsLieuxLiees),
		LinkedIndividuals:     collections.InterfaceToInt32Slice(mobilier.PersonnesPhysiquesLiees),
		LinkedLegalEntities:   collections.InterfaceToInt32Slice(mobilier.PersonnesMoralesLiees),
		Themes:                collections.InterfaceToStringSlice(mobilier.Themes),
	})
}

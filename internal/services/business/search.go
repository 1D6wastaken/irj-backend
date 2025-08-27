package business

import (
	"context"
	"net/http"
	"strconv"

	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	"irj/pkg/collections"
	_http "irj/pkg/http"

	"github.com/rs/zerolog"
)

func (b *BusinessService) Search(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	var (
		result *api.SearchResult
		apiErr *_http.APIError
	)

	if r.URL.Query().Get("q") != "" {
		result, apiErr = searchWithText(subCtx, b, r)
		if apiErr != nil {
			return apiErr
		}

		return _http.WriteJSONResponse(w, http.StatusOK, result)
	}

	result, apiErr = searchWithoutText(subCtx, b, r)
	if apiErr != nil {
		return apiErr
	}

	return _http.WriteJSONResponse(w, http.StatusOK, result)
}

func parseQueryWithText(r *http.Request) (*queries.SearchGlobalParams, *_http.APIError) {
	q := r.URL.Query().Get("q")

	req, err := _http.DecodeAndValidateJSONBody[*api.SearchFilters](r)
	if err != nil {
		return nil, _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	limit, offset := extractLimitAndOffset(r)

	params := queries.SearchGlobalParams{
		Q:                      q,
		IncludeMonumentsLieux:  false,
		Siecles:                req.Centuries,
		Pays:                   req.Countries,
		Region:                 req.Regions,
		Departement:            req.Departments,
		Commune:                req.Cities,
		NaturesMonu:            nil,
		EtatsMonu:              nil,
		MateriauxMonu:          nil,
		IncludeMobiliersImages: false,
		NaturesMob:             nil,
		EtatsMob:               nil,
		MateriauxMob:           nil,
		TechniquesMob:          nil,
		IncludePersMorales:     false,
		NaturesPersMo:          nil,
		IncludePersPhysiques:   false,
		Professions:            nil,
		ModesDeplacements:      nil,
		OffsetParam:            offset,
		LimitParam:             limit,
	}

	allSources := true

	if req.MonumentsLieux != nil {
		allSources = false
		params.IncludeMonumentsLieux = true
		params.NaturesMonu = req.MonumentsLieux.Natures
		params.EtatsMonu = req.MonumentsLieux.States
		params.MateriauxMonu = req.MonumentsLieux.Materials
	}

	if req.MobiliersImages != nil {
		allSources = false
		params.IncludeMobiliersImages = true
		params.NaturesMob = req.MobiliersImages.Natures
		params.EtatsMob = req.MobiliersImages.States
		params.MateriauxMob = req.MobiliersImages.Materials
		params.TechniquesMob = req.MobiliersImages.Techniques
	}

	if req.PersMorales != nil {
		allSources = false
		params.IncludePersMorales = true
		params.NaturesPersMo = req.PersMorales.Natures
	}

	if req.PersPhysiques != nil {
		allSources = false
		params.IncludePersPhysiques = true
		params.Professions = req.PersPhysiques.Professions
		params.ModesDeplacements = req.PersPhysiques.Travels
	}

	if allSources {
		params.IncludeMonumentsLieux = true
		params.IncludeMobiliersImages = true
		params.IncludePersMorales = true
		params.IncludePersPhysiques = true
	}

	return &params, nil
}

func parseQueryWithoutText(r *http.Request) (*queries.SearchGlobalNoTextParams, *_http.APIError) {
	req, err := _http.DecodeAndValidateJSONBody[*api.SearchFilters](r)
	if err != nil {
		return nil, _http.ErrBadRequest.Msg("unable to decode request body").Err(err)
	}

	limit, offset := extractLimitAndOffset(r)

	params := queries.SearchGlobalNoTextParams{
		IncludeMonumentsLieux:  false,
		Siecles:                req.Centuries,
		Pays:                   req.Countries,
		Region:                 req.Regions,
		Departement:            req.Departments,
		Commune:                req.Cities,
		NaturesMonu:            nil,
		EtatsMonu:              nil,
		MateriauxMonu:          nil,
		IncludeMobiliersImages: false,
		NaturesMob:             nil,
		EtatsMob:               nil,
		MateriauxMob:           nil,
		TechniquesMob:          nil,
		IncludePersMorales:     false,
		NaturesPersMo:          nil,
		IncludePersPhysiques:   false,
		Professions:            nil,
		ModesDeplacements:      nil,
		OffsetParam:            offset,
		LimitParam:             limit,
	}

	allSources := true

	if req.MonumentsLieux != nil {
		allSources = false
		params.IncludeMonumentsLieux = true
		params.NaturesMonu = req.MonumentsLieux.Natures
		params.EtatsMonu = req.MonumentsLieux.States
		params.MateriauxMonu = req.MonumentsLieux.Materials
	}

	if req.MobiliersImages != nil {
		allSources = false
		params.IncludeMobiliersImages = true
		params.NaturesMob = req.MobiliersImages.Natures
		params.EtatsMob = req.MobiliersImages.States
		params.MateriauxMob = req.MobiliersImages.Materials
		params.TechniquesMob = req.MobiliersImages.Techniques
	}

	if req.PersMorales != nil {
		allSources = false
		params.IncludePersMorales = true
		params.NaturesPersMo = req.PersMorales.Natures
	}

	if req.PersPhysiques != nil {
		allSources = false
		params.IncludePersPhysiques = true
		params.Professions = req.PersPhysiques.Professions
		params.ModesDeplacements = req.PersPhysiques.Travels
	}

	if allSources {
		params.IncludeMonumentsLieux = true
		params.IncludeMobiliersImages = true
		params.IncludePersMorales = true
		params.IncludePersPhysiques = true
	}

	return &params, nil
}

func extractLimitAndOffset(r *http.Request) (int32, int32) {
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	page, _ := strconv.ParseInt(pageStr, 10, 32)

	if limit <= 0 {
		limit = 25
	}

	if page <= 0 {
		page = 1
	}

	offset := (int32(page) - 1) * int32(limit)

	return int32(limit), offset
}

func searchWithText(ctx context.Context, s *BusinessService, r *http.Request) (*api.SearchResult, *_http.APIError) {
	logger := zerolog.Ctx(ctx)

	params, apiErr := parseQueryWithText(r)
	if apiErr != nil {
		return nil, apiErr
	}

	search, err := s.postgresService.Queries.SearchGlobal(ctx, *params)
	if err != nil {
		return nil, _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	total := int64(0)
	if len(search) > 0 {
		total = search[0].TotalCount
	}

	items := make([]*api.ListItem, 0, len(search))

	for i := range search {
		row := search[i]

		var medias []*api.Media

		nocoMedias, err := s.parseMedias(row.Medias)
		if err == nil {
			medias = collections.Map(nocoMedias, func(m NocoMedia) *api.Media {
				return &api.Media{
					ID:    &m.ID,
					Title: &m.Title,
				}
			})
		} else {
			logger.Error().Err(err).Msg("error while parsing medias")
		}

		items = append(items, &api.ListItem{
			Source:      &row.Source,
			ID:          &row.ID,
			Title:       &row.Title,
			Medias:      medias,
			Natures:     collections.InterfaceToStringSlice(row.Natures),
			Centuries:   collections.InterfaceToStringSlice(row.Siecles),
			Professions: collections.InterfaceToStringSlice(row.Professions),
		})
	}

	return &api.SearchResult{
		Total: &total,
		Items: items,
	}, nil
}

func searchWithoutText(ctx context.Context, s *BusinessService, r *http.Request) (*api.SearchResult, *_http.APIError) {
	logger := zerolog.Ctx(ctx)

	params, apiErr := parseQueryWithoutText(r)
	if apiErr != nil {
		return nil, apiErr
	}

	search, err := s.postgresService.Queries.SearchGlobalNoText(ctx, *params)
	if err != nil {
		return nil, _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	total := int64(0)
	if len(search) > 0 {
		total = search[0].TotalCount
	}

	items := make([]*api.ListItem, 0, len(search))

	for i := range search {
		row := search[i]

		var medias []*api.Media

		nocoMedias, err := s.parseMedias(row.Medias)
		if err == nil {
			medias = collections.Map(nocoMedias, func(m NocoMedia) *api.Media {
				return &api.Media{
					ID:    &m.ID,
					Title: &m.Title,
				}
			})
		} else {
			logger.Error().Err(err).Msg("error while parsing medias")
		}

		items = append(items, &api.ListItem{
			Source:      &row.Source,
			ID:          &row.ID,
			Title:       &row.Title,
			Medias:      medias,
			Natures:     collections.InterfaceToStringSlice(row.Natures),
			Centuries:   collections.InterfaceToStringSlice(row.Siecles),
			Professions: collections.InterfaceToStringSlice(row.Professions),
		})
	}

	return &api.SearchResult{
		Total: &total,
		Items: items,
	}, nil
}

package business

import (
	"context"
	"net/http"
	"strconv"

	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	"irj/pkg/collections"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

func (b *BusinessService) GetPays(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	p, err := b.postgresService.Queries.GetPays(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	pJSON := collections.FilterMap(p, func(row queries.GetPaysRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, pJSON)
}

func (b *BusinessService) GetRegions(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	reg, err := b.postgresService.Queries.GetRegions(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	regJSON := collections.FilterMap(reg, func(row queries.GetRegionsRow) (*api.RegionFilter, bool) {
		var pf api.BasicFilter

		if row.PaysID.Valid {
			pf = api.BasicFilter{
				ID:   utils.PtrTo(int64(row.PaysID.Int32)),
				Name: utils.PtrTo(row.PaysName.String),
			}
		}

		return &api.RegionFilter{
			BasicFilter: api.BasicFilter{
				ID:   utils.PtrTo(int64(row.ID)),
				Name: utils.PtrTo(row.Name.String),
			},
			Pays: &pf,
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, regJSON)
}

func (b *BusinessService) GetDepartements(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	d, err := b.postgresService.Queries.GetDepartements(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	dJSON := collections.FilterMap(d, func(row queries.GetDepartementsRow) (*api.DepartmentFilter, bool) {
		var (
			pf api.BasicFilter
			rf api.RegionFilter
		)

		if row.PaysID.Valid {
			pf = api.BasicFilter{
				ID:   utils.PtrTo(int64(row.PaysID.Int32)),
				Name: utils.PtrTo(row.PaysName.String),
			}
		}

		if row.RegionID.Valid {
			rf = api.RegionFilter{
				BasicFilter: api.BasicFilter{
					ID:   utils.PtrTo(int64(row.RegionID.Int32)),
					Name: utils.PtrTo(row.RegionName.String),
				},
				Pays: &pf,
			}
		}

		return &api.DepartmentFilter{
			BasicFilter: api.BasicFilter{
				ID:   utils.PtrTo(int64(row.ID)),
				Name: utils.PtrTo(row.Name.String),
			},
			Region: &rf,
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, dJSON)
}

func (b *BusinessService) GetCommunes(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	q := r.URL.Query().Get("q")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	if limit == 0 {
		limit = 10
	}

	c, err := b.postgresService.Queries.SearchCommunesPaginated(subCtx, queries.SearchCommunesPaginatedParams{
		NomCommune: pgtype.Text{String: q + "%%", Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	cJSON := collections.FilterMap(c, func(row queries.SearchCommunesPaginatedRow) (*api.CommuneFilter, bool) {
		var (
			pf api.BasicFilter
			rf api.RegionFilter
			df api.DepartmentFilter
		)

		if row.PaysID.Valid {
			pf = api.BasicFilter{
				ID:   utils.PtrTo(int64(row.PaysID.Int32)),
				Name: utils.PtrTo(row.PaysName.String),
			}
		}

		if row.RegionID.Valid {
			rf = api.RegionFilter{
				BasicFilter: api.BasicFilter{
					ID:   utils.PtrTo(int64(row.RegionID.Int32)),
					Name: utils.PtrTo(row.RegionName.String),
				},
				Pays: &pf,
			}
		}

		if row.DepartementID.Valid {
			df = api.DepartmentFilter{
				BasicFilter: api.BasicFilter{
					ID:   utils.PtrTo(int64(row.DepartementID.Int32)),
					Name: utils.PtrTo(row.DepartementName.String),
				},
				Region: &rf,
			}
		}

		return &api.CommuneFilter{
			BasicFilter: api.BasicFilter{
				ID:   utils.PtrTo(int64(row.ID)),
				Name: utils.PtrTo(row.CommuneName.String),
			},
			Department: &df,
		}, row.CommuneName.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, cJSON)
}

func (b *BusinessService) GetSiecles(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	s, err := b.postgresService.Queries.GetSiecles(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	sJSON := collections.FilterMap(s, func(row queries.GetSieclesRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, sJSON)
}

func (b *BusinessService) GetEtatsConservation(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	ec, err := b.postgresService.Queries.GetEtatsConservation(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	ecJSON := collections.FilterMap(ec, func(row queries.GetEtatsConservationRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, ecJSON)
}

func (b *BusinessService) GetNaturesMonument(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	n, err := b.postgresService.Queries.GetNaturesMonu(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	nJSON := collections.FilterMap(n, func(row queries.GetNaturesMonuRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, nJSON)
}

func (b *BusinessService) GetMateriaux(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	m, err := b.postgresService.Queries.GetMateriaux(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	mJSON := collections.FilterMap(m, func(row queries.GetMateriauxRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, mJSON)
}

func (b *BusinessService) GetNaturesMobilier(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	n, err := b.postgresService.Queries.GetNaturesMob(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	nJSON := collections.FilterMap(n, func(row queries.GetNaturesMobRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, nJSON)
}

func (b *BusinessService) GetTechniquesMobilier(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	t, err := b.postgresService.Queries.GetTechniquesMob(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	tJSON := collections.FilterMap(t, func(row queries.GetTechniquesMobRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, tJSON)
}

func (b *BusinessService) GetNaturesPersonnesMorales(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	p, err := b.postgresService.Queries.GetNaturesPersonnesMorales(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	pJSON := collections.FilterMap(p, func(row queries.GetNaturesPersonnesMoralesRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, pJSON)
}

func (b *BusinessService) GetProfessions(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	p, err := b.postgresService.Queries.GetProfessions(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	pJSON := collections.FilterMap(p, func(row queries.GetProfessionsRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, pJSON)
}

func (b *BusinessService) GetDeplacements(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	d, err := b.postgresService.Queries.GetDeplacements(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	dJSON := collections.FilterMap(d, func(row queries.GetDeplacementsRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, dJSON)
}

func (b *BusinessService) GetHistoricalPeriods(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	d, err := b.postgresService.Queries.GetHistoricalPeriods(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	dJSON := collections.FilterMap(d, func(row queries.GetHistoricalPeriodsRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, dJSON)
}

func (b *BusinessService) GetThemes(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	d, err := b.postgresService.Queries.GetThemes(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("error while fetching data").Err(err)
	}

	dJSON := collections.FilterMap(d, func(row queries.GetThemesRow) (*api.BasicFilter, bool) {
		return &api.BasicFilter{
			ID:   utils.PtrTo(int64(row.ID)),
			Name: utils.PtrTo(row.Name.String),
		}, row.Name.Valid
	})

	return _http.WriteJSONResponse(w, http.StatusOK, dJSON)
}

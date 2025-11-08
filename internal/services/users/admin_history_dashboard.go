package users

import (
	"context"
	"net/http"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	"irj/pkg/collections"
	_http "irj/pkg/http"
	"irj/pkg/utils"

	"github.com/go-openapi/strfmt"
)

func (u *UserService) GetAdminHistoryDashboard(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	token, ok := subCtx.Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	if token.Grade != string(queries.UserGradeADMIN) {
		return _http.ErrForbidden.Msg("admin access required")
	}

	limit, offset := extractLimitAndOffset(r)

	history, err := u.postgresService.Queries.GetAllHistory(subCtx, queries.GetAllHistoryParams{
		OffsetParam: offset,
		LimitParam:  limit,
	})
	if err != nil {
		return _http.ErrInternalError.Msg("unable to get admin history").Err(err)
	}

	result := collections.Map(history, func(item queries.GetAllHistoryRow) *api.GetAdminContributionsItem {
		return &api.GetAdminContributionsItem{
			Date:           utils.PtrTo(strfmt.DateTime(item.Date.Time)),
			Event:          utils.PtrTo(string(item.Type)),
			Category:       utils.PtrTo(item.Comment.String),
			DocumentID:     item.DocumentID.Int32,
			UserID:         strfmt.UUID(item.UserID.String),
			UserFirstname:  item.UserPrenom.String,
			UserLastname:   item.UserNom.String,
			UserMail:       strfmt.Email(item.UserEmail.String),
			AdminID:        strfmt.UUID(item.AdminID.String),
			AdminFirstname: item.AdminPrenom.String,
			AdminLastname:  item.AdminNom.String,
			AdminMail:      strfmt.Email(item.AdminEmail.String),
		}
	})

	var total int64

	if len(history) > 0 {
		total = history[0].TotalCount
	}

	return _http.WriteJSONResponse(w, http.StatusOK, api.GetAdminContributions{
		Total:  &total,
		Events: result,
	})
}

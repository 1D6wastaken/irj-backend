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
	"github.com/jackc/pgx/v5/pgtype"
)

func (u *UserService) GetUserHistory(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	token, ok := subCtx.Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	limit, offset := extractLimitAndOffset(r)

	history, err := u.postgresService.Queries.GetHistoryByUserID(subCtx, queries.GetHistoryByUserIDParams{
		OffsetParam: offset,
		LimitParam:  limit,
		UserID: pgtype.Text{
			String: token.ID,
			Valid:  true,
		},
	})
	if err != nil {
		return _http.ErrInternalError.Msg("unable to get user history").Err(err)
	}

	result := collections.Map(history, func(item queries.GetHistoryByUserIDRow) *api.GetUserHistoryItem {
		return &api.GetUserHistoryItem{
			Date:       utils.PtrTo(strfmt.DateTime(item.Date.Time)),
			Event:      utils.PtrTo(string(item.Type)),
			Category:   utils.PtrTo(item.Comment.String),
			DocumentID: item.DocumentID.Int32,
		}
	})

	var total int64

	if len(history) > 0 {
		total = history[0].TotalCount
	}

	return _http.WriteJSONResponse(w, http.StatusOK, api.GetUserHistory{
		Total:  &total,
		Events: result,
	})
}

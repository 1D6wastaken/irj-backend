package password

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
)

func (p *PasswordService) CheckNewPasswordToken(w http.ResponseWriter, r *http.Request) *_http.APIError {
	token := httprouter.ParamsFromContext(r.Context()).ByName("token")
	if token == "" {
		return _http.ErrBadRequest.Msg("Missing path parameter").WithDetails("token is required")
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)

	defer cancel()

	_, err := p.postgresService.Queries.GetResetPasswordByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _http.ErrNotFound.Msg("Token not found")
		}

		return _http.ErrInternalError.Msg("Error while fetching data").Err(err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

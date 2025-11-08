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

func (u *UserService) GetAdminUsersDashboard(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	token, ok := subCtx.Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	if token.Grade != string(queries.UserGradeADMIN) {
		return _http.ErrForbidden.Msg("admin access required")
	}

	users, err := u.postgresService.Queries.GetUsers(subCtx)
	if err != nil {
		return _http.ErrInternalError.Msg("unable to get users").Err(err)
	}

	result := collections.Map(users, func(u queries.GetUsersRow) api.GetAdminUserInfo {
		return api.GetAdminUserInfo{
			ID:                   utils.PtrTo(strfmt.UUID(u.ID)),
			Firstname:            utils.PtrTo(u.Prenom),
			Lastname:             utils.PtrTo(u.Nom),
			Mail:                 utils.PtrTo(strfmt.Email(u.Email)),
			MailConfirm:          u.EmailConfirm,
			Phone:                u.Telephone.String,
			Organization:         u.Organisation.String,
			Domain:               string(u.Domaine),
			LastLogin:            strfmt.DateTime(u.LastLogin.Time),
			ValidatedByFirstname: u.AdminPrenom.String,
			ValidatedByLastname:  u.AdminNom.String,
			ValidatedByMail:      strfmt.Email(u.AdminEmail.String),
			Grade:                utils.PtrTo(string(u.Grade)),
			CreationDate:         utils.PtrTo(strfmt.DateTime(u.DateCreation.Time)),
		}
	})

	return _http.WriteJSONResponse(w, http.StatusOK, result)
}

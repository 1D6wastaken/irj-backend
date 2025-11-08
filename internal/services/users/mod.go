package users

import (
	"net/http"
	"strconv"
	"time"

	"irj/internal/config"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/framework"
	"irj/pkg/utils"
)

const defaultTimeOut = time.Second * 30

type UserService struct {
	stopper         *utils.AppStopper
	config          *config.Config
	smtpService     *smtp.SMTPService
	postgresService *framework.DB[*queries.Queries]
}

//nolint:lll
func NewUserService(stopper *utils.AppStopper, cfg *config.Config, smtpService *smtp.SMTPService, postgresService *framework.DB[*queries.Queries]) *UserService {
	return &UserService{
		stopper:         stopper,
		config:          cfg,
		smtpService:     smtpService,
		postgresService: postgresService,
	}
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

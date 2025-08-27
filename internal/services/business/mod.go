package business

import (
	"time"

	"irj/internal/config"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/framework"
	"irj/pkg/utils"
)

const defaultTimeOut = time.Second * 30

type BusinessService struct {
	stopper         *utils.AppStopper
	config          *config.Config
	smtpService     *smtp.SMTPService
	postgresService *framework.DB[*queries.Queries]
}

//nolint:lll
func NewBusinessService(stopper *utils.AppStopper, cfg *config.Config, smtpService *smtp.SMTPService, postgresService *framework.DB[*queries.Queries]) *BusinessService {
	return &BusinessService{
		stopper:         stopper,
		config:          cfg,
		smtpService:     smtpService,
		postgresService: postgresService,
	}
}

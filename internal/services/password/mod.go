package password

import (
	"time"

	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/framework"
	"irj/pkg/utils"
)

const defaultTimeOut = time.Second * 30

type PasswordService struct {
	stopper         *utils.AppStopper
	smtpService     *smtp.SMTPService
	postgresService *framework.DB[*queries.Queries]
}

//nolint:lll
func NewPasswordService(stopper *utils.AppStopper, smtpService *smtp.SMTPService, postgresService *framework.DB[*queries.Queries]) *PasswordService {
	return &PasswordService{
		stopper:         stopper,
		smtpService:     smtpService,
		postgresService: postgresService,
	}
}

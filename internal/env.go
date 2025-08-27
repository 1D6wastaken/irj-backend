package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"irj/internal/config"
	"irj/internal/middlewares"
	queries "irj/internal/postgres/_generated"
	"irj/internal/services/business"
	"irj/internal/services/password"
	"irj/internal/services/users"
	"irj/internal/smtp"
	"irj/pkg/dto"
	"irj/pkg/framework"
	"irj/pkg/glog"
	_http "irj/pkg/http"
	"irj/pkg/postgres"
	"irj/pkg/utils"

	"github.com/go-co-op/gocron"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

const (
	defaultServiceName = "irj-backend"
)

type Env struct {
	Stopper         *utils.AppStopper
	Config          *config.Config
	Logger          *glog.Logger
	DBService       *postgres.Service
	SMTPService     *smtp.SMTPService
	PostgresService *framework.DB[*queries.Queries]
	UserService     *users.UserService
	PasswordService *password.PasswordService
	BusinessService *business.BusinessService
	Scheduler       *gocron.Scheduler
}

func InitEnv(ctx context.Context, stopper *utils.AppStopper) (framework.Server, error) {
	cfg, err := config.LoadConfig("conf/config.yaml")
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(cfg.Logs.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := glog.InitLogger(cfg.Logs.Level, defaultServiceName, file)

	dbService, err := postgres.Connect(ctx, cfg.Database.DSN, cfg.Database.MaxConns)
	if err != nil {
		return nil, err
	}

	postgresService, err := framework.NewDB[*queries.Queries](ctx, cfg.Database.DSN, func(db *pgxpool.Pool) *queries.Queries {
		return queries.New(db)
	}, nil)
	if err != nil {
		return nil, err
	}

	smtpService := smtp.NewSMTPService(&logger, cfg.SMTP.Host, cfg.SMTP.APIKey, cfg.SMTP.From.Name, cfg.SMTP.From.Email)

	userService := users.NewUserService(stopper, cfg, smtpService, postgresService)

	passwordService := password.NewPasswordService(stopper, smtpService, postgresService)

	businessService := business.NewBusinessService(stopper, cfg, smtpService, postgresService)

	env := &Env{
		Stopper:         stopper,
		Config:          cfg,
		Logger:          &logger,
		DBService:       dbService,
		SMTPService:     smtpService,
		PostgresService: postgresService,
		UserService:     userService,
		PasswordService: passwordService,
		BusinessService: businessService,
		Scheduler:       gocron.NewScheduler(time.UTC),
	}

	if err := env.Cron(ctx); err != nil {
		return nil, err
	}

	return env, nil
}

func (e *Env) GetHTTPServerAddr() string {
	return fmt.Sprintf(":%d", e.Config.Port)
}

func (e *Env) GetServerName() string {
	return defaultServiceName
}

func (e *Env) Probes() []func(context.Context) *dto.ServiceItem {
	return []func(ctx context.Context) *dto.ServiceItem{
		e.DBService.Status,
	}
}

//nolint:lll,funlen
func (e *Env) Routes(router *httprouter.Router) {
	injectLoggerMiddleware := _http.InjectDefaultLogger(e.Logger)
	logMiddleware := _http.LogInputOutput()

	public := _http.Middlewares(injectLoggerMiddleware, logMiddleware)

	// Public APIs
	router.HandlerFunc(http.MethodGet, "/api/v1/pays", public(_http.HandlerMiddleware(e.BusinessService.GetPays)))
	router.HandlerFunc(http.MethodGet, "/api/v1/regions", public(_http.HandlerMiddleware(e.BusinessService.GetRegions)))
	router.HandlerFunc(http.MethodGet, "/api/v1/departments", public(_http.HandlerMiddleware(e.BusinessService.GetDepartements)))
	router.HandlerFunc(http.MethodGet, "/api/v1/cities", public(_http.HandlerMiddleware(e.BusinessService.GetCommunes)))
	router.HandlerFunc(http.MethodGet, "/api/v1/centuries", public(_http.HandlerMiddleware(e.BusinessService.GetSiecles)))
	router.HandlerFunc(http.MethodGet, "/api/v1/conservation_states", public(_http.HandlerMiddleware(e.BusinessService.GetEtatsConservation)))
	router.HandlerFunc(http.MethodGet, "/api/v1/building_natures", public(_http.HandlerMiddleware(e.BusinessService.GetNaturesMonument)))
	router.HandlerFunc(http.MethodGet, "/api/v1/materials", public(_http.HandlerMiddleware(e.BusinessService.GetMateriaux)))
	router.HandlerFunc(http.MethodGet, "/api/v1/furnitures_natures", public(_http.HandlerMiddleware(e.BusinessService.GetNaturesMobilier)))
	router.HandlerFunc(http.MethodGet, "/api/v1/furnitures_techniques", public(_http.HandlerMiddleware(e.BusinessService.GetTechniquesMobilier)))
	router.HandlerFunc(http.MethodGet, "/api/v1/legal_entity_natures", public(_http.HandlerMiddleware(e.BusinessService.GetNaturesPersonnesMorales)))
	router.HandlerFunc(http.MethodGet, "/api/v1/professions", public(_http.HandlerMiddleware(e.BusinessService.GetProfessions)))
	router.HandlerFunc(http.MethodGet, "/api/v1/travels", public(_http.HandlerMiddleware(e.BusinessService.GetDeplacements)))
	router.HandlerFunc(http.MethodGet, "/api/v1/historical_periods", public(_http.HandlerMiddleware(e.BusinessService.GetHistoricalPeriods)))
	router.HandlerFunc(http.MethodGet, "/api/v1/themes", public(_http.HandlerMiddleware(e.BusinessService.GetThemes)))

	router.HandlerFunc(http.MethodGet, "/api/v1/monuments_lieux/:id", public(_http.HandlerMiddleware(e.BusinessService.GetMonumentLieu)))
	router.HandlerFunc(http.MethodGet, "/api/v1/mobiliers_images/:id", public(_http.HandlerMiddleware(e.BusinessService.GetMobilierImage)))
	router.HandlerFunc(http.MethodGet, "/api/v1/personnes_morales/:id", public(_http.HandlerMiddleware(e.BusinessService.GetPersonneMorale)))
	router.HandlerFunc(http.MethodGet, "/api/v1/personnes_physiques/:id", public(_http.HandlerMiddleware(e.BusinessService.GetPersonnePhysique)))

	router.HandlerFunc(http.MethodPost, "/api/v1/search", public(_http.HandlerMiddleware(e.BusinessService.Search)))

	router.HandlerFunc(http.MethodGet, "/media/:id", public(e.BusinessService.GetMediaByID))

	router.HandlerFunc(http.MethodPost, "/api/v1/login", public(_http.HandlerMiddleware(e.UserService.Login)))

	router.HandlerFunc(http.MethodPost, "/api/v1/users", public(_http.HandlerMiddleware(e.UserService.RegisterUser)))
	router.HandlerFunc(http.MethodPost, "/api/v1/password-reset", public(_http.HandlerMiddleware(e.PasswordService.RequestNewPassword)))
	router.HandlerFunc(http.MethodGet, "/api/v1/password-reset/:token", public(_http.HandlerMiddleware(e.PasswordService.CheckNewPasswordToken)))
	router.HandlerFunc(http.MethodPost, "/api/v1/password-reset/:token", public(_http.HandlerMiddleware(e.PasswordService.ResetPassword)))
	router.HandlerFunc(http.MethodGet, "/api/v1/email/:token/validate", public(_http.HandlerMiddleware(e.UserService.UserEmailConfirm)))

	// Protected APIs
	protected := _http.Middlewares(injectLoggerMiddleware, logMiddleware, middlewares.CheckJWT(e.Config))

	router.HandlerFunc(http.MethodGet, "/api/v1/users", protected(_http.HandlerMiddleware(e.UserService.GetPendingUsers)))
	router.HandlerFunc(http.MethodPatch, "/api/v1/users/:id", protected(_http.HandlerMiddleware(e.UserService.ApproveRejectUser)))
	router.HandlerFunc(http.MethodDelete, "/api/v1/users/:id", protected(_http.HandlerMiddleware(e.UserService.DeleteUser)))
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id", protected(_http.HandlerMiddleware(e.UserService.GetUserInfo)))
	router.HandlerFunc(http.MethodPut, "/api/v1/users/:id", protected(_http.HandlerMiddleware(e.UserService.UpdateUserInfo)))

	router.HandlerFunc(http.MethodPost, "/api/v1/medias", protected(_http.HandlerMiddleware(e.BusinessService.UploadImage)))
	router.HandlerFunc(http.MethodPost, "/api/v1/monuments_lieux", protected(_http.HandlerMiddleware(e.BusinessService.CreateMonumentLieu)))
	router.HandlerFunc(http.MethodPost, "/api/v1/mobiliers_images", protected(_http.HandlerMiddleware(e.BusinessService.CreateMobilierImage)))
	router.HandlerFunc(http.MethodPost, "/api/v1/personnes_morales", protected(_http.HandlerMiddleware(e.BusinessService.CreatePersonneMorale)))
	router.HandlerFunc(http.MethodPost, "/api/v1/personnes_physiques", protected(_http.HandlerMiddleware(e.BusinessService.CreatePersonnePhysique)))

	router.HandlerFunc(http.MethodGet, "/api/v1/monuments_lieux", protected(_http.HandlerMiddleware(e.BusinessService.GetPendingMonumentsLieux)))
	router.HandlerFunc(http.MethodGet, "/api/v1/mobiliers_images", protected(_http.HandlerMiddleware(e.BusinessService.GetPendingMobiliersImages)))
	router.HandlerFunc(http.MethodGet, "/api/v1/personnes_morales", protected(_http.HandlerMiddleware(e.BusinessService.GetPendingPersonnesMorales)))
	router.HandlerFunc(http.MethodGet, "/api/v1/personnes_physiques", protected(_http.HandlerMiddleware(e.BusinessService.GetPendingPersonnesPhysiques)))

	router.HandlerFunc(http.MethodPatch, "/api/v1/monuments_lieux/:id", protected(_http.HandlerMiddleware(e.BusinessService.ApproveRejectMonumentLieu)))
	router.HandlerFunc(http.MethodPatch, "/api/v1/mobiliers_images/:id", protected(_http.HandlerMiddleware(e.BusinessService.ApproveRejectMobilierImage)))
	router.HandlerFunc(http.MethodPatch, "/api/v1/personnes_morales/:id", protected(_http.HandlerMiddleware(e.BusinessService.ApproveRejectPersonneMorale)))
	router.HandlerFunc(http.MethodPatch, "/api/v1/personnes_physiques/:id", protected(_http.HandlerMiddleware(e.BusinessService.ApproveRejectPersonnePhysique)))

	router.HandlerFunc(http.MethodPut, "/api/v1/monuments_lieux/:id", protected(_http.HandlerMiddleware(e.BusinessService.UpdateMonumentLieu)))
	router.HandlerFunc(http.MethodPut, "/api/v1/mobiliers_images/:id", protected(_http.HandlerMiddleware(e.BusinessService.UpdateMobilierImage)))
	router.HandlerFunc(http.MethodPut, "/api/v1/personnes_morales/:id", protected(_http.HandlerMiddleware(e.BusinessService.UpdatePersonneMorale)))
	router.HandlerFunc(http.MethodPut, "/api/v1/personnes_physiques/:id", protected(_http.HandlerMiddleware(e.BusinessService.UpdatePersonnePhysique)))
}

func (e *Env) GetLogger() *glog.Logger {
	return e.Logger
}

func (e *Env) Close(_ context.Context) error {
	errs := make([]error, 0)

	e.DBService.Close()

	return errors.Join(errs...)
}

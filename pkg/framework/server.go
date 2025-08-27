package framework

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"irj/pkg/dto"
	"irj/pkg/glog"
	"irj/pkg/utils"

	_http "irj/pkg/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type Server interface {
	GetHTTPServerAddr() string
	GetServerName() string
	Probes() []func(context.Context) *dto.ServiceItem
	Routes(router *httprouter.Router)
	GetLogger() *glog.Logger
	Close(ctx context.Context) error
}

// "ReadinessProbe ensures that the container is healthy to serve incoming traffic".
func readiness(server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		status := &dto.Service{
			Name:    server.GetServerName(),
			Status:  dto.ServiceStatusOK,
			Details: "",
			Services: utils.Map(server.Probes(), func(probe func(context.Context) *dto.ServiceItem) *dto.ServiceItem {
				return probe(ctx)
			}),
		}

		httpStatus := http.StatusOK

		if utils.Some(status.Services, func(s *dto.ServiceItem) bool {
			return s.Status == dto.ServiceItemStatusKO
		}) {
			status.Status = dto.ServiceItemStatusKO
			httpStatus = http.StatusInternalServerError
		}

		_ = _http.WriteJSONResponse(w, httpStatus, status)
	}
}

// "LivenessProbe serves as a diagnostic check to confirm if the container is alive".
func liveness(server Server) http.HandlerFunc {
	status := &dto.Service{
		Name:     server.GetServerName(),
		Status:   dto.ServiceStatusOK,
		Details:  "",
		Services: nil,
	}

	return func(w http.ResponseWriter, _ *http.Request) {
		_ = _http.WriteJSONResponse(w, http.StatusOK, status)
	}
}

func Run(serverBuilder func(ctx context.Context, stopper *utils.AppStopper) (Server, error)) {
	var (
		ctx     = context.Background()
		stopper = utils.NewAppStopper(ctx)
		server  = utils.Must(serverBuilder(ctx, stopper))
		logger  = server.GetLogger()
	)

	router := httprouter.New()

	server.Routes(router)

	// Add probes and default handlers
	router.HandlerFunc(http.MethodGet, "/health", readiness(server))
	router.HandlerFunc(http.MethodGet, "/status", liveness(server))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_http.ErrNotFound.
			WithDetails(fmt.Sprintf("route %s %s doesn't exist", r.Method, r.URL.Path)).
			Write(w)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_http.ErrMethodNotAllowed.Write(w)
	})

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		zerolog.Ctx(r.Context()).Error().
			Str("stacktrace", string(debug.Stack())).
			Str("error", fmt.Sprintf("%v", err)).
			Msg("panic unexpectedly")

		_http.ErrInternalError.WithDetails(fmt.Sprintf("%v", err)).Write(w)
	}

	var (
		//nolint:exhaustruct
		apiServer = &http.Server{
			Addr:              server.GetHTTPServerAddr(),
			Handler:           _http.CORSMiddleware(router),
			ReadHeaderTimeout: time.Minute,
		}

		apiServerCtx = make(chan struct{})
	)

	go func() {
		defer close(apiServerCtx)

		logger.Info().Msgf("start HTTP API server on %s", apiServer.Addr)

		if err := apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Err(err).Msg("HTTP API server has stopped unexpectedly")
		}
	}()

	select {
	case <-stopper.Done():
		logger.Info().Msg("exit signal received: shutting down")

		if err := apiServer.Shutdown(context.WithoutCancel(ctx)); err != nil {
			logger.Err(err).Msg("HTTP API server failed to be shutdown")
		}
	case <-apiServerCtx:
		stopper.Cancel()
	}

	logger.Info().Msg("app is closing")
	defer logger.Info().Msg("app has been gracefully closed")

	if err := stopper.Wait(25 * time.Second); err != nil {
		logger.Err(err).Msg("failed to wait app stopper")
	}

	if err := server.Close(ctx); err != nil {
		logger.Err(err).Msg("app failed to close")
	}
}

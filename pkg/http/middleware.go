package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"

	"irj/pkg/glog"
)

type ResponseSpy struct {
	ResponseWriter http.ResponseWriter
	body           string
	status         int
}

func (w *ResponseSpy) Write(p []byte) (int, error) {
	w.body = string(p)

	return w.ResponseWriter.Write(p)
}

func (w *ResponseSpy) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseSpy) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseSpy) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Middlewares(middlewares ...Middleware) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}

		return handler
	}
}

// InjectDefaultLogger will inject a logger into incoming request context.
func InjectDefaultLogger(logger *glog.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(logger.WithContext(r.Context()))
			next.ServeHTTP(w, r)
		}
	}
}

func LogInputOutput(excludes ...string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if (r.Method == http.MethodGet && (r.URL.Path == "/readiness" || r.URL.Path == "/liveness")) ||
				slices.Contains(excludes, fmt.Sprintf("%s %s", r.Method, r.URL.Path)) {
				next(w, r)

				return
			}

			startTime := time.Now()
			logger := zerolog.Ctx(r.Context())

			logger.Info().Object(glog.KeyHTTP, glog.LogRequest(r)).Msg("incoming HTTP request")

			reSpy := &ResponseSpy{
				ResponseWriter: w,
				body:           "",
				status:         0,
			}

			next(reSpy, r)

			res := http.Response{
				Status:           "",
				StatusCode:       reSpy.status,
				Proto:            "",
				ProtoMajor:       0,
				ProtoMinor:       0,
				Header:           w.Header(),
				Body:             io.NopCloser(bytes.NewBufferString(reSpy.body)),
				ContentLength:    0,
				TransferEncoding: nil,
				Close:            false,
				Uncompressed:     false,
				Trailer:          nil,
				Request:          nil,
				TLS:              nil,
			}

			logger.Debug().Object(glog.KeyHTTP, glog.LogResponse(r, &res, startTime)).Msg("outgoing HTTP response")
		}
	}
}

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method found, HS256 expected")
	ErrFailedToFindClaims      = errors.New("failed to find claims inside jwt")
	ErrFailedToCheckClaims     = errors.New("invalid claims found inside jwt")
)

func ParseAndCheckJWT(secret []byte, authorization string, checkClaims func(claims jwt.MapClaims) bool) error {
	authorization = strings.TrimPrefix(authorization, "Bearer ")

	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return secret, nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrFailedToFindClaims
	}

	if !checkClaims(claims) {
		return ErrFailedToCheckClaims
	}

	return nil
}

func ExpectJWT(hmacSecret []byte, condition func(jwt.MapClaims) bool) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				ErrUnauthorized.WithDetails(`"Authorization" header is required`).Write(w)

				return
			}

			if err := ParseAndCheckJWT(hmacSecret, authorization, condition); err != nil {
				ErrUnauthorized.Err(err).Write(w)

				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func ExtractJWT(secret []byte, authorization string) (jwt.MapClaims, error) {
	authorization = strings.TrimPrefix(authorization, "Bearer ")

	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrFailedToFindClaims
	}

	return claims, nil
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // ou "*" pour tester
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Gérer la requête preflight OPTIONS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		// Passer au handler suivant
		next.ServeHTTP(w, r)
	})
}

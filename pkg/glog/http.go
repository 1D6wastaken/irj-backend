package glog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type (
	requestLogger struct {
		path     string
		method   string
		rawQuery string
		body     string
	}

	responseLogger struct {
		path       string
		method     string
		rawQuery   string
		statusCode int
		body       string
		startTime  time.Time
	}
)

func LogRequest(r *http.Request) zerolog.LogObjectMarshaler {
	body := ""

	if r.Body != nil {
		rawBody, err := io.ReadAll(r.Body)
		if err == nil && len(rawBody) > 0 {
			body = string(rawBody)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	}

	return requestLogger{
		path:     r.URL.Path,
		method:   r.Method,
		rawQuery: r.URL.RawQuery,
		body:     body,
	}
}

func LogResponse(r *http.Request, w *http.Response, startTime time.Time) zerolog.LogObjectMarshaler {
	body := ""

	if w.Body != nil {
		rawBody, err := io.ReadAll(w.Body)
		if err == nil && len(rawBody) > 0 {
			body = string(rawBody)
		}

		w.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	}

	return &responseLogger{
		path:       r.URL.Path,
		method:     r.Method,
		rawQuery:   r.URL.RawQuery,
		statusCode: w.StatusCode,
		body:       body,
		startTime:  startTime,
	}
}

func (r requestLogger) MarshalZerologObject(e *zerolog.Event) {
	e.Str(KeyURI, r.path).
		Str(KeyMethod, r.method)

	if r.rawQuery != "" {
		e.Str(KeyQuery, r.rawQuery)
	}

	if r.body != "" {
		e.Str(KeyRequestBody, r.body)
	}
}

func (r *responseLogger) MarshalZerologObject(e *zerolog.Event) {
	e.Str(KeyURI, r.path).
		Str(KeyMethod, r.method)

	if r.rawQuery != "" {
		e.Str(KeyQuery, r.rawQuery)
	}

	e.Int(KeyStatusCode, r.statusCode).
		Str(KeyDuration, fmt.Sprintf("%dms", time.Since(r.startTime).Milliseconds()))

	if r.body != "" {
		e.Str(KeyResponseBody, r.body)
	}
}

package handler

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
)

func logRequest(r *http.Request, status, size int, duration time.Duration) {
	ll := hlog.FromRequest(r)

	ref := r.Header.Get("Referer")

	ll.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("referer", ref).
		Int64("req_size", r.ContentLength).
		Int("resp_code", status).
		Int("resp_size", size).
		Int64("resp_duration_ms", duration.Milliseconds()).
		Msg("request processed")
}

// LogRequest middleware that logs the request
func (h *Handler) LogRequest(next http.Handler) http.Handler {
	next = hlog.AccessHandler(logRequest)(next)
	next = hlog.RequestIDHandler("req_id", "Request-Id")(next)

	return hlog.NewHandler(h.Logger)(next)
}

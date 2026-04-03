package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func SlogLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" || r.URL.Path == "/favicon.ico" || strings.HasPrefix(r.URL.Path, "/static/") {
				next.ServeHTTP(w, r)
				return
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				reqID := middleware.GetReqID(r.Context())
				l := logger.With(
					slog.String("class", "middleware_http"),
					slog.String("req_id", reqID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("remote_addr", r.RemoteAddr),
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.BytesWritten()),
					slog.Duration("duration", time.Since(t1)),
				)
				if ww.Status() >= 500 {
					l.Error("Request completed")
				} else {
					l.Info("Request completed")
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

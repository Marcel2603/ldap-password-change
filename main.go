package main

import (
	"ldap-password-change/cmd/config"
	changepassword "ldap-password-change/internal/handler/change-password"
	"ldap-password-change/internal/handler/health"
	"ldap-password-change/internal/handler/index"
	staticfiles "ldap-password-change/internal/handler/static-files"
	custommw "ldap-password-change/internal/middleware"
	"ldap-password-change/internal/service/ldap"
	"ldap-password-change/internal/validation"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:generate go tool github.com/a-h/templ/cmd/templ generate

func main() {
	configuration := config.Configuration

	var level slog.Level
	if err := level.UnmarshalText([]byte(configuration.Log.Level)); err != nil {
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting server")

	r, err := setupApp(configuration, ldap.CreateWrapper(), logger)
	if err != nil {
		logger.Error("Error starting server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Listening on :" + configuration.Server.Port)
	err = http.ListenAndServe(":"+configuration.Server.Port, r)
	if err != nil {
		logger.Error("Error starting server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func setupApp(configuration config.Config, wrapper ldap.Wrapper, logger *slog.Logger) (*chi.Mux, error) {
	r := setupServerRouter(configuration, logger)

	service := ldap.CreateService(configuration.Ldap, wrapper, logger)
	validator, errValidator := validation.CreateValidator(configuration.Validation)
	if errValidator != nil {
		return nil, errValidator
	}

	r.Get("/", index.Handler)
	r.Get("/health", health.LivenessHandler)
	r.Get("/health/live", health.LivenessHandler)
	r.Get("/health/ready", health.ReadinessHandler(service))

	r.Get("/favicon.ico", staticfiles.HandleFavicon)
	r.Get("/static/*", staticfiles.Handler)
	r.Get("/custom/*", staticfiles.Handler)

	r.Post("/change-password", changepassword.Handler(service, validator, logger))

	return r, nil
}

func setupServerRouter(configuration config.Config, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{configuration.Server.Host},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           600, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(custommw.SlogLogger(logger))
	r.Use(middleware.Recoverer)
	return r
}

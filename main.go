package main

import (
	"errors"
	"github.com/go-chi/cors"
	"ldap-password-change/cmd/config"
	changepassword "ldap-password-change/internal/handler/change-password"
	"ldap-password-change/internal/handler/index"
	staticfiles "ldap-password-change/internal/handler/static-files"
	"ldap-password-change/internal/service/ldap"
	"ldap-password-change/internal/validation"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:generate go tool github.com/a-h/templ/cmd/templ generate

func main() {
	slog.Info("Starting server")
	configuration := config.Configuration
	r := setupServerRouter(configuration)

	r.Get("/", index.Handler)

	r.Get("/favicon.ico", staticfiles.HandleFavicon)
	r.Get("/static/*", staticfiles.Handler)

	service, errService := ldap.CreateService(configuration.Ldap, ldap.CreateWrapper())
	validator, errValidator := validation.CreateValidator(configuration.Validation)
	if errService != nil || errValidator != nil {
		log.Fatalf("Error creating services: %s\n", errors.Join(errService, errValidator).Error())
	}
	r.Post("/change-password", changepassword.Handler(service, validator))

	slog.Info("Listening on :" + configuration.Server.Port)
	err := http.ListenAndServe(":"+configuration.Server.Port, r)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err.Error())
	}
}

func setupServerRouter(configuration config.Config) *chi.Mux {
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
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

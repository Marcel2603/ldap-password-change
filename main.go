package main

import (
	"fmt"
	"github.com/go-chi/cors"
	"ldap-password-change/cmd/config"
	changepassword "ldap-password-change/internal/handler/change-password"
	"ldap-password-change/internal/handler/index"
	staticfiles "ldap-password-change/internal/handler/static-files"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:generate go run github.com/a-h/templ/cmd/templ generate

func main() {
	fmt.Println("Starting server")
	configuration := config.Configuration
	r := setupServerRouter(configuration)

	r.Get("/", index.Handler)

	r.Get("/favicon.ico", staticfiles.HandleFavicon)
	r.Get("/static/*", staticfiles.Handler)

	r.Post("/change-password", changepassword.Handler)
	fmt.Println("Listening on :" + configuration.Server.Port)
	http.ListenAndServe(":"+configuration.Server.Port, r)
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

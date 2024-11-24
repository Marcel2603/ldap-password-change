package main

import (
	"fmt"
	"github.com/go-chi/cors"
	"ldap-password-change/cmd/config"
	changepassword "ldap-password-change/internal/handler/change-password"
	"ldap-password-change/internal/handler/index"
	staticfiles "ldap-password-change/internal/handler/static-files"
	"ldap-password-change/internal/handler/validate"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("Starting server")
	configuration := config.Get()
	r := setupServerRouter(configuration)

	r.Get("/", index.Handler)

	r.Get("/favicon.ico", staticfiles.HandleFavicon)
	r.Get("/static/*", staticfiles.Handler)

	r.Post("/change-password", changepassword.Handler)
	r.Post("/validate/*", validate.Handler)
	fmt.Println("Listening on :" + configuration.Port)
	http.ListenAndServe(":"+configuration.Port, r)
}

func setupServerRouter(configuration config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{configuration.Domain},
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

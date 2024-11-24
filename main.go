package main

import (
	"fmt"
	"github.com/go-chi/cors"
	changepassword "ldap-password-change/internal/handler/change-password"
	"ldap-password-change/internal/handler/index"
	staticfiles "ldap-password-change/internal/handler/static-files"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3333"},
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

	r.Get("/", index.Handler)

	r.Get("/static/*", staticfiles.Handler)

	r.Post("/change-password", changepassword.Handler)

	fmt.Println("Listening on :3333")
	http.ListenAndServe(":3333", r)
}
